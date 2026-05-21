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

type GoPackage struct {
	Dir     string   `json:"Dir"`
	GoFiles []string `json:"GoFiles"`
	Imports []string `json:"Imports"`
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

	const libraryPrefix = "github.com/aruaru0/ac-library-go/"

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
		writeResult(fset, file, outPath)
		return
	}

	importedDecls := []ast.Decl{}
	stdlibImports := make(map[string]bool)

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
				importedDecls = append(importedDecls, decl)
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
	for i, decl := range importedDecls {
		importedDecls[i] = astutil.Apply(decl, removePrefix, nil).(ast.Decl)
	}

	// 6. 収集したライブラリコードを main.go の末尾に結合
	file.Decls = append(file.Decls, importedDecls...)

	// 7. 出力
	writeResult(fset, file, outPath)
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
func writeResult(fset *token.FileSet, file *ast.File, outPath string) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		log.Fatalf("ASTのフォーマットに失敗しました: %v", err)
	}

	if outPath == "" {
		fmt.Print(buf.String())
	} else {
		if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
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
