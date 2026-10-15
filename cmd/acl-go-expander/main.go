package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// 対象ライブラリのインポートパスの接頭辞（末尾スラッシュ付き）
const libraryPrefix = "github.com/aruaru0/ac-library-go/"

type GoPackage struct {
	Dir        string   `json:"Dir"`
	ImportPath string   `json:"ImportPath"`
	GoFiles    []string `json:"GoFiles"`
	Imports    []string `json:"Imports"`
	Module     Module   `json:"Module"`
}

// Module は go list -json が返すモジュール情報のうち、出典URLの生成に必要な項目
type Module struct {
	Path    string  `json:"Path"`
	Version string  `json:"Version"`
	Dir     string  `json:"Dir"`
	Main    bool    `json:"Main"`
	Replace *Module `json:"Replace"`
}

// lexChunk は 1 つのライブラリソースファイルに由来する宣言のまとまり
type lexChunk struct {
	sourceURL string
	decls     []ast.Decl
}

func main() {
	var outPath string
	flag.StringVar(&outPath, "o", "", "出力ファイルパス (デフォルトは標準出力)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("使用方法: acl-go-expander [-o <出力ファイル>] <解析対象のファイル.go>")
	}
	inputPath := args[0]

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, inputPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("入力ファイルのパースに失敗しました: %v", err)
	}

	// 1. main.go 自体で定義されているトップレベルの名前を登録
	declaredNames := make(map[string]bool)
	for _, decl := range file.Decls {
		if name := getDeclName(decl); name != "" {
			declaredNames[name] = true
		}
	}

	// 処理するパッケージのキューと訪問済みマップ
	queue := []string{}
	targetImports := make(map[string]string) // pkgName -> importPath (例: "dsu" -> "github.com...")
	visitedPkgs := make(map[string]bool)

	// 2. main.go の import から対象ライブラリを抽出
	var toDelete []string
	for _, imp := range file.Imports {
		if imp == nil || imp.Path == nil {
			continue
		}
		path := strings.Trim(imp.Path.Value, `"`)
		if strings.HasPrefix(path, libraryPrefix) {
			pkgName := filepath.Base(path)
			if imp.Name != nil {
				pkgName = imp.Name.Name
			}
			targetImports[pkgName] = path
			queue = append(queue, path)
			toDelete = append(toDelete, path)
		}
	}

	// 安全に一括削除
	for _, path := range toDelete {
		astutil.DeleteImport(fset, file, path)
	}

	// 対象のインポートがない場合はそのまま出力して終了
	if len(queue) == 0 {
		writeResult(fset, file, nil, "", outPath)
		return
	}

	var chunks []*lexChunk
	stdlibImports := make(map[string]bool)
	var bannerRepoBase, bannerRev string

	// 3. 依存関係の再帰的な解決とソースコードの収集
	for len(queue) > 0 {
		currPath := queue[0]
		queue = queue[1:]

		if visitedPkgs[currPath] {
			continue
		}
		visitedPkgs[currPath] = true

		pkgInfo, err := getPackageInfo(currPath)
		if err != nil {
			log.Fatalf("パッケージ情報の取得に失敗しました (%s): %v", currPath, err)
		}

		// 出典URL生成のためのリポジトリ情報（最初のパッケージのものを代表として使う）
		if bannerRepoBase == "" {
			bannerRepoBase, bannerRev = repoBaseAndRev(pkgInfo.Module)
		}

		// 依存関係を走査
		for _, imp := range pkgInfo.Imports {
			if strings.HasPrefix(imp, libraryPrefix) {
				// ライブラリ内部での別モジュールへの依存を発見した場合、キューに追加
				depPkgName := filepath.Base(imp)
				targetImports[depPkgName] = imp
				queue = append(queue, imp)
			} else if isStdLib(imp) {
				// 標準ライブラリ依存を収集
				stdlibImports[imp] = true
			}
		}

		// パッケージ内のすべてのGoソースファイルをパース
		for _, goFile := range pkgInfo.GoFiles {
			fullPath := filepath.Join(pkgInfo.Dir, goFile)
			pkgFile, err := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
			if err != nil {
				log.Fatalf("ライブラリファイルのパースに失敗しました (%s): %v", fullPath, err)
			}

			// 非公開のトップレベル定義をパッケージ名付きにリネームして競合を防ぐ
			renamePrivateDecls(pkgFile, filepath.Base(pkgInfo.Dir))

			chunk := &lexChunk{sourceURL: sourceURL(pkgInfo, fullPath)}
			for _, decl := range pkgFile.Decls {
				// import 宣言はスキップ
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
					continue
				}

				// 重複定義（ヘルパー関数 min, max 等）の競合を防ぐ
				name := getDeclName(decl)
				if name != "" {
					if declaredNames[name] {
						continue
					}
					declaredNames[name] = true
				}
				chunk.decls = append(chunk.decls, decl)
			}
			if len(chunk.decls) > 0 {
				chunks = append(chunks, chunk)
			}
		}
	}

	// 4. 標準ライブラリの import を main.go にマージ
	for std := range stdlibImports {
		astutil.AddImport(fset, file, std)
	}

	// 5. ASTの書き換え: パッケージプレフィックス (例: dsu.NewDSU -> NewDSU) の削除
	removePrefix := func(c *astutil.Cursor) bool {
		n := c.Node()
		if expr, ok := n.(*ast.SelectorExpr); ok {
			if ident, ok := expr.X.(*ast.Ident); ok {
				if _, exists := targetImports[ident.Name]; exists {
					c.Replace(expr.Sel)
				}
			}
		}
		return true
	}

	// main.go のコードを書き換え
	file = astutil.Apply(file, removePrefix, nil).(*ast.File)

	// インポートしたライブラリ側のコードも同様に書き換え（内部参照の解決）
	for _, chunk := range chunks {
		for i, decl := range chunk.decls {
			chunk.decls[i] = astutil.Apply(decl, removePrefix, nil).(ast.Decl)
		}
	}

	// 6. 出力（埋め込みコードごとにソースURLコメントを付与）
	banner := ""
	if bannerRepoBase != "" {
		banner = buildBanner(bannerRepoBase, bannerRev)
	}
	writeResult(fset, file, chunks, banner, outPath)
}

// go list を使用してパッケージ情報を JSON で取得
func getPackageInfo(importPath string) (*GoPackage, error) {
	cmd := exec.Command("go", "list", "-json", importPath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("go list 実行エラー: %w", err)
	}

	var pkg GoPackage
	if err := json.Unmarshal(stdout.Bytes(), &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

// 標準ライブラリかどうかの判定
func isStdLib(importPath string) bool {
	pkg, err := build.Import(importPath, "", build.FindOnly)
	return err == nil && pkg.Goroot
}

// 出典URLのリポジトリbaseとリビジョンを求める
func repoBaseAndRev(mod Module) (string, string) {
	return "https://" + urlModulePath(mod), moduleRev(mod)
}

// 指定ファイルのGitHub permalink (ファイル単位) を生成する
func sourceURL(pkg *GoPackage, absFile string) string {
	base, rev := repoBaseAndRev(pkg.Module)
	rel, err := filepath.Rel(moduleDir(pkg.Module), absFile)
	if err != nil {
		rel = filepath.Base(absFile)
	}
	return base + "/blob/" + rev + "/" + filepath.ToSlash(rel)
}

// 出典URLに使うモジュールパスを決める（replace されている場合はフォーク先を優先）
func urlModulePath(mod Module) string {
	if mod.Replace != nil && looksLikeModulePath(mod.Replace.Path) {
		return mod.Replace.Path
	}
	if mod.Path != "" {
		return mod.Path
	}
	// モジュール情報が取得できない場合は既知のライブラリのリポジトリにフォールバック
	return strings.TrimSuffix(libraryPrefix, "/")
}

// モジュールのソースディレクトリ（replace されている場合は置換先）
func moduleDir(mod Module) string {
	if mod.Replace != nil && mod.Replace.Dir != "" {
		return mod.Replace.Dir
	}
	return mod.Dir
}

// 出典URLに埋め込むリビジョン（タグ／コミットハッシュ）を求める
func moduleRev(mod Module) string {
	version := mod.Version
	if mod.Replace != nil && mod.Replace.Version != "" {
		version = mod.Replace.Version
	}
	if version != "" {
		// 疑似バージョン (v0.0.0-20240101120000-abcdef123456) はコミットハッシュを使う
		if hash := pseudoVersionHash(version); hash != "" {
			return hash
		}
		return version
	}
	// ローカルの開発チェックアウト等では git の HEAD を試す
	if dir := moduleDir(mod); dir != "" {
		if hash := gitHEAD(dir); hash != "" {
			return hash
		}
	}
	return "main"
}

// 疑似バージョン形式の文字列からコミットハッシュ部分を取り出す
func pseudoVersionHash(version string) string {
	i := strings.LastIndex(version, "-")
	if i < 0 {
		return ""
	}
	hash := version[i+1:]
	if len(hash) == 12 && isHex(hash) {
		return hash
	}
	return ""
}

func isHex(s string) bool {
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}

// looksLikeModulePath はモジュールパスらしさ（先頭要素がドメイン）をざっくり判定する
func looksLikeModulePath(path string) bool {
	if path == "" || strings.HasPrefix(path, ".") || strings.HasPrefix(path, "/") {
		return false
	}
	first := path
	if i := strings.Index(path, "/"); i >= 0 {
		first = path[:i]
	}
	return strings.Contains(first, ".")
}

// git rev-parse HEAD でコミットハッシュを取得する（取得できない場合は空文字）
func gitHEAD(dir string) string {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "HEAD")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(stdout.String())
}

// 先頭に付与する生成バナーを組み立てる
func buildBanner(repoBase, rev string) string {
	return fmt.Sprintf("// This file was assembled by acl-go-expander before the contest.\n"+
		"// Tool: %s/blob/%s/cmd/acl-go-expander/main.go\n"+
		"// The library code embedded below (see \"source:\" comments) was published\n"+
		"// beforehand, in accordance with AtCoder's rules on pre-contest AI-generated code.",
		repoBase, rev)
}

// 宣言からトップレベルの識別子名を取得する
func getDeclName(decl ast.Decl) string {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		if d.Recv == nil { // 通常の関数のみ（メソッドは除く）
			return d.Name.Name
		}
	case *ast.GenDecl:
		if d.Tok == token.TYPE || d.Tok == token.CONST || d.Tok == token.VAR {
			if len(d.Specs) > 0 {
				switch s := d.Specs[0].(type) {
				case *ast.TypeSpec:
					return s.Name.Name
				case *ast.ValueSpec:
					if len(s.Names) > 0 {
						return s.Names[0].Name
					}
				}
			}
		}
	}
	return ""
}

// 結果を出力 (フォーマット済みのGoコード)
func writeResult(fset *token.FileSet, file *ast.File, chunks []*lexChunk, banner string, outPath string) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		log.Fatalf("ASTのフォーマットに失敗しました: %v", err)
	}

	var sb strings.Builder
	if banner != "" {
		sb.WriteString(banner)
		sb.WriteString("\n\n")
	}
	sb.WriteString(buf.String())

	// 埋め込みコードを、その出典URLコメントとともに結合する
	for _, chunk := range chunks {
		if len(chunk.decls) == 0 {
			continue
		}
		sb.WriteString("\n// ---- source: ")
		sb.WriteString(chunk.sourceURL)
		sb.WriteString(" ----\n\n")
		for _, decl := range chunk.decls {
			var dbuf bytes.Buffer
			if err := format.Node(&dbuf, fset, decl); err != nil {
				log.Fatalf("埋め込みコードのフォーマットに失敗しました: %v", err)
			}
			sb.WriteString(dbuf.String())
			sb.WriteString("\n\n")
		}
	}

	// 全体をまとめて整形し直す（コメント挿入後の体裁を整える）
	src, err := format.Source([]byte(sb.String()))
	if err != nil {
		log.Fatalf("出力コードの整形に失敗しました: %v", err)
	}

	if outPath == "" {
		fmt.Print(string(src))
	} else {
		if err := os.WriteFile(outPath, src, 0644); err != nil {
			log.Fatalf("出力ファイルの書き込みに失敗しました: %v", err)
		}
	}
}

// ライブラリの非公開トップレベル定義をリネームして競合を防ぐ
func renamePrivateDecls(file *ast.File, pkgName string) {
	renameMap := make(map[string]string)
	for _, decl := range file.Decls {
		name := getDeclName(decl)
		if name != "" && isPrivateName(name) {
			renameMap[name] = pkgName + "_" + name
		}
	}

	if len(renameMap) == 0 {
		return
	}

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		n := c.Node()
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}

		newName, exists := renameMap[ident.Name]
		if !exists {
			return true
		}

		parent := c.Parent()
		if parent == nil {
			return true
		}

		switch p := parent.(type) {
		case *ast.SelectorExpr:
			if p.Sel == ident {
				return true
			}
		case *ast.KeyValueExpr:
			if p.Key == ident {
				return true
			}
		case *ast.FuncDecl:
			if p.Recv != nil && p.Name == ident {
				return true
			}
		case *ast.Field:
			return true
		}

		c.Replace(&ast.Ident{Name: newName, NamePos: ident.NamePos})
		return true
	}, nil)
}

func isPrivateName(name string) bool {
	if len(name) == 0 {
		return false
	}
	r := name[0]
	return r >= 'a' && r <= 'z'
}
