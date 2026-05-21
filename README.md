# AC Library Go

AtCoder Library (AC Library) を Idiomatic な Go 言語 (Go 1.18+ の Generics 対応) に移植したライブラリです。競技プログラミングなどで要求される高パフォーマンスと安全性を両立しています。

## 含まれるモジュール

- **データ構造**: `dsu` (Union-Find), `fenwicktree`, `segtree`, `lazysegtree`
- **グラフ・フロー**: `scc`, `twosat`, `maxflow`, `mincostflow`
- **数学・文字列**: `math`, `modint`, `mint` (モジュラ演算ヘルパー), `convolution`, `string`

## 使い方

本ライブラリは、用途に合わせて以下の2つの方法でご利用いただけます。

---

### 方法1. GitHubからGoモジュールとしてインポートして使う

Go プロジェクトで通常の外部パッケージとしてインポートし、利用する方法です。

**1. ライブラリのインストール**
ターミナルで以下のコマンドを実行し、プロジェクトにモジュールを追加します。
（※リポジトリURLは実際のGitHubパスに置き換えてください）

```bash
go get github.com/aruaru0/ac-library-go
```

**2. コード内でのインポートと利用**
使いたいモジュールを `import` に追加して利用します。

```go
package main

import (
	"fmt"

	// 必要なモジュールをインポート
	"github.com/aruaru0/ac-library-go/dsu"
	"github.com/aruaru0/ac-library-go/maxflow"
)

func main() {
	// DSU (Union-Find) の例
	uf := dsu.NewDSU(10)
	uf.Merge(1, 2)
	fmt.Println("1 and 2 are same?:", uf.Same(1, 2))

	// MaxFlow の例
	mf := maxflow.NewMFGraph[int](4)
	mf.AddEdge(0, 1, 10)
	mf.AddEdge(1, 3, 5)
	mf.AddEdge(0, 2, 5)
	mf.AddEdge(2, 3, 10)
	
	fmt.Println("Max Flow:", mf.Flow(0, 3)) // 出力: 10
}
```

---

### 方法2. `package main` のコード内にコピー＆ペーストして使う (競技プログラミング向け)

競技プログラミングの提出システム (AtCoderなど) では外部モジュールを import できない場合があります。その際は、必要なモジュールのソースコードを直接ご自身のコードに貼り付けてください。

**手順**
1. 自分の解答用ファイル (例: `main.go`) を作成し、`package main` を宣言します。
2. 利用したい機能の `.go` ファイル (例: `dsu/dsu.go`) の中身を開きます。
3. `package dsu` 等のパッケージ宣言と `import` 文 **以外** の部分（型定義、関数、メソッド群など）をコピーします。
4. コピーしたコードをご自身の `main.go` 内のファイル末尾などにペーストします。
5. ペーストした構造体や関数は同じ `package main` に属するため、パッケージ名をつけずにそのまま呼び出すことができます。

**例 (DSUをコピペして使う場合):**

```go
package main

import (
	"fmt"
)

func main() {
	// コピペしたコードを使うため dsu.New ではなく NewDSU などの名前に変更、
	// もしくはそのまま関数名を利用します (ここでは元が New() なので、衝突を避けるため NewDSU 等にすると安全です)。
	uf := NewDSU(10)
	uf.Merge(1, 2)
	fmt.Println(uf.Same(1, 2))
}

// ---------------- 以下、コピー＆ペーストしたコード ----------------

type DSU struct {
	parentOrSize []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := range p {
		p[i] = -1
	}
	return &DSU{parentOrSize: p}
}

func (d *DSU) Merge(a, b int) int {
	// ... (実装の続き) ...
}

func (d *DSU) Same(a, b int) bool {
	// ... (実装の続き) ...
}

// ...
```

---

### 方法3. 提出用コード自動生成ツール (`acl-go-expander`) を使う (おすすめ)

手動でのコピー＆ペーストやパッケージプレフィックスの削除を自動化するコマンドラインツールを用意しています。本ツールは、通常のインポート文を使ったコードから、提出用の自己完結した単一ファイルを自動生成します。

**1. ツールのインストール**
```bash
go install github.com/aruaru0/ac-library-go/cmd/acl-go-expander@latest
```
*(※ `~/go/bin` にパスが通っている必要があります)*

**2. 開発コードの作成**
通常通り外部パッケージをインポートして解答コードを書きます（エディタのコード補完やコンパイルチェックの恩恵をすべて受けられます）。

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/dsu"
)

func main() {
	uf := dsu.NewDSU(5)
	uf.Merge(0, 1)
	fmt.Println(uf.Same(0, 1))
}
```

**3. 提出用コードの生成**
ターミナルで以下のコマンドを実行し、ライブラリコードがマージされた単一ファイルを生成します。

```bash
acl-go-expander -o submit.go main.go
```

生成された `submit.go` をそのまま AtCoder などのジャッジシステムに提出してください。

---

## modint と mint の違いについて

Go言語には演算子オーバーロードがないため、本家 ACL 同様の構造体に値をラップする `modint` パッケージを使用すると、数式の記述が `a.Add(b).Mul(c)` のように冗長になり、直感的な実装が難しくなる傾向があります。

これに対処するため、書きやすさを重視したモジュラ演算ヘルパーモジュール **`mint`** を用意しました。
`mint` は値自体をラップするのではなく、「法（mod）と計算用バッファを持つコンテクスト」を提供し、通常の `int` 同士で演算を行います（例: `m.Add(a, m.Mul(b, c))`）。

さらに、`mint` には実用性を高めるため以下の便利な拡張も内包されています：
* **組み合わせ計算**: $_n\text{C}_r$ や階乗・逆元テーブルの $O(N)$ 一括事前構築
* **行列累乗**: $A^p \pmod{mod}$ の高速計算（フィボナッチ数の高速計算等に有用）

## ドキュメント

各機能の詳細な仕様や計算量については、`doc/` フォルダ内のマークダウンドキュメントをご参照ください。
