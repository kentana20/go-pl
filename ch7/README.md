7章 - インタフェース
---

## 7.7 - http.Handler インタフェース

ListenAndServe関数はHandlerインタフェースを必要とする。

- セール商品をドルの価格に対応させているDBを持つECサイトの例

```go
func main() {
	db := database{"shoes": 50, "socks": 5}
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}
```

マップ型のdatabase型にServeHTTPメソッドを結びつけている → http.Handlerインタフェースを満足している。
このサンプルだと、すべてのHTTPリクエストに対して同じレスポンスを返すので、今度はパスごとに処理をわける例。
`req.URL.path` に基づいてswitch文で処理を分岐させている。

```go
func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/list":
		for item, price := range db {
			fmt.Fprintf(w, "%s: %s\n", item, price)
		}
	case "/price":
		item := req.URL.Query().Get("item")
		price, ok := db[item]
		if !ok {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "no such item: %q\n", item)
			return
		}
		fmt.Fprintf(w, "%s\n", price)
	default:
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such page: %s\n", req.URL)
	}
}

// curl http://localhost:8000/price\?item\=socks // $5.00
```

URLとハンドラの関連付けを簡単にするためにリクエストマルチプレクサである **ServeMux** をnet/httpパッケージが提供している。 **ServeMux** はhttp.Handlerのコレクションを単一のhttp.Handlerへ集める（同じインタフェースを満足する型はそれぞれ代入可能な特性を利用している）。↓は **ServeMux** 利用の例。

```go
func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.Handle("/price", http.HandlerFunc(db.price))
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
```

**ServeMux** を http.ListenAndServe のハンドラとして指定する形。
http.HandlerFuncはhttp.Handlerインタフェースを満足する関数型。関数値にインタフェースを満足させるアダプタ？（よくわからん）
↑の例だと、db.list関数をhttp.HandlerFunc関数型に変換してhttp.Handlerインタフェースを満足させてハンドラ登録している。
**ServeMux** はHandleFuncメソッドを持っているので、型変換せずにハンドラ登録することも可能。

```go
mux.HandleFunc(db.list)
mux.HandleFunc(db.price)
```

さらにnet/httpパッケージはDefaultServeMuxというグローバルなServeMuxインスタンスを持っているのでこれを使えば新しく **ServeMux** インスタンスを作る必要がなくなる。

```go
http.HandleFunc("/list", db.list)
http.HandleFunc("/price", db.price)
log.Fatal(http.ListenAndServe("localhost:8000", nil)) // DefaultServeMuxを使う場合はnilを渡す
```

WebサーバはそれぞれのハンドラをGroutineで起動するので、変数の扱い（書き換えなど）に注意する必要がある。詳しくは8章で。

## 7-8 - error インタフェース

エラーメッセージを返す単一メソッドをもったインタフェース。

```go
type error interface {
    Error() string
}
```

errorsパッケージもシンプル。

```go
package errors

func New(text string) error { return &errorString{text} }
type errorString struct { text string }
func (e *errorString) Error() string { return e.text }
```

errorStringのポインタ型がerrorインタフェースを満足しているので、エラーが複数起こった場合でもそれぞれのerrors.Newは別のインスタンスが割り当てられる（それぞれのエラーが競合しない）。
実際には errors.New よりフォーマット付きのラッパー関数である fmt.Errorf 関数をよく使う。

```go
package fmt

import "errors"

func Errorf(format string, args ...interface{}) error {
    return errors.New(Sprintf(format, args...))
}
```

errorインタフェースを満足する他の型については 7-11 で。

## 7-9 - 式評価器

算術式に対する評価器の構築。インタフェース `Expr` を使う。

- 式言語の構成要素
  - 浮動小数点数リテラル
  - 二項演算子 +, -, *, /
  - 単項演算 -x, +x
  - 関数呼び出し pow(x, y), sin(x), sqrt(x)
  - 変数 pi

すべての値はfloat64。

```go
//Var - 変数
type Var string

// 数値定数
type literal float64

// 単項演算式
type unary struct {
	op rune
	x  Expr
}

// 二項演算式
type binary struct {
	op   rune
	x, y Expr
}

// 関数呼び出し式
type call struct {
	fn   string
	args []Expr
}

//Env - 変数名を値へ対応づける環境
type Env map[Var]float64

//Expr - インタフェース
type Expr interface {
	Eval(env Env) float64
}

//Eval - 環境の検索
func (v Var) Eval(env Env) float64 {
	return env[v]
}

//Eval - リテラル値の返却
func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

//unaryに対するEvalメソッド
func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

//binaryに対するEvalメソッド
func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

// callに対するEvalメソッド
func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}
```

メソッドのいくつかは失敗する可能性があるのでEvalに対するテストを書く。

```go
func TestEval(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
		{"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
		{"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
	}
	var prevExpr string
	for _, test := range tests {
		// Print expr only when it changes.
		if test.expr != prevExpr {
			fmt.Printf("\n%s\n", test.expr)
			prevExpr = test.expr
		}
		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err) // parse error
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			t.Errorf("%s.Eval() in %v = %q, want %q\n",
				test.expr, test.env, got, test.want)
		}
	}
}
```

Parse関数をgo getで見つけろって書いてあるけど、よくわからん。
テストが通ったものとして次へ。

静的な検査をするためにExprインタフェースにCheckメソッドを追加する。

```go
type Expr interface {
    Eval(env Env) float64
    Check(vars map[Var]bool) error
}
```

以下、それぞれの型に対するCheckメソッド。

```go
func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (literal) Check(vars map[Var]bool) error {
	return nil
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}
```

Parse関数は構文エラーを報告し、Check関数は意味論的エラーを報告するらしい。
これでEvalパッケージは完成して、次はそれを利用して式を受け取ってその関数による面を描画するWebアプリケーションを書いてみる。

```go
const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func surface(w io.Writer, f func(x, y float64) float64) {
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, judgea := corner(i+1, j)
			bx, by, judgeb := corner(i, j)
			cx, cy, judgec := corner(i, j+1)
			dx, dy, judged := corner(i+1, j+1)

			if !judgea || !judgeb || !judgec || !judged {
				continue
			}

			fmt.Fprintf(w, "<polygon points='%g,%g,%g,%g,%g,%g,%g,%g' />\n", ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintln(w, "</svg>")
}

func corner(i, j int) (float64, float64, bool) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z, judge := f(x, y)

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	return sx, sy, judge
}

func parseAndCheck(s string) (eval.Expr, error) {
	if s == "" {
		return nil, fmt.Errorf("empty expression")
	}
	expr, err := eval.Parse(s)
	if err != nil {
		return nil, err
	}
	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		return nil, err
	}
	for v := range vars {
		if v != "x" && v != "y" && v != "r" {
			return nil, fmt.Errorf("undefined variable: %s", v)
		}
	}
	return expr, nil
}

func f(x, y float64) (float64, bool) {
	r := math.Hypot(x, y)

	if math.IsInf((math.Sin(r) / r), 0) {
		return 0, false
	}

	return math.Sin(r) / r, true
}

func plot(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	expr, err := parseAndCheck(r.Form.Get("expr"))
	if err != nil {
		http.Error(w, "bad expr: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	surface(w, func(x, y float64) float64 {
		r := math.Hypot(x, y) // distance from (0,0)
		return expr.Eval(eval.Env{"x": x, "y": y, "r": r})
	})
}

func main() {
	http.HandleFunc("/plot", plot)
	log.Fatal(http.ListenAndServe("localhost:8001", nil))
}
```

## 7-10 - 型アサーション

インタフェース値へ提供される演算。  `x.(T)` と書く。

- x: インタフェース型の式
- T: 断定型(Assert)

そのオペランドの動的な型が断定型と一致するかを検査できる。

### 例1: Tが具象型の場合はxの動的な型がTと同一かを検査する

検査成功時は x の動的な値を返す。失敗時はpanicになる。

```go
var w io.Writer
w = os.Stdout
f := w.(*os.File)
c := w.(*bytes.Buffer)
```

### 例2: Tがインタフェースの場合はxの動的な型がTを満足するかを検査する

検査成功時はインタフェース型Tを持つ。

```go
var w io.Writer
w = os.Stdout
rw := w.(io.ReadWriter) // rw はio.ReadWriter型

w = new(ByteCounter)
rw = w.(io.ReadWriter) // panic
```

オペランドがnilのインタフェース値の場合、型アサーションは失敗する。

```go
var w io.Writer = os.Stdout
f, ok := w.(*os.File)
b, ok := w.(*bytes.Buffer) // panicにはならず、okにはfalseが入る
```

## 7-11 - 型アサーションによるエラーの区別

osパッケージによるファイル操作を例にとった型アサーションのエラー区別。

- IsExist: ファイル作成時の存在チェック
- IsNotExist: ファイル読み込み時の存在チェック
- IsPermission: ファイル操作の権限チェック

素朴な実装。エラー文言から判定する。こんなこと絶対にやらない。

```go
func IsNotExist(err error) bool {
    return strings.Contains(err.Error(), "file does not exist")
}
```

実践的な例。構造化されたエラー値を表現する。ファイル操作時の失敗はPathError型として定義されている。

```go
type PathError struct {
    Op string // 操作
    Path string // ファイルパス
    Err error // エラー内容
}

func (e *PathError) Error() string {
    return e.Op + " " + e.Path + ": " + e.Err.Error()
}
```

エラーの種類を区別する場合は型アサーションが使える。

```go
import (
    "errors"
    "syscall"
)

var ErrNotExist = errors.New("file does not exist")

func IsNotExist(err error) bool {
    // 型アサーション
    if pe, ok := err.(*PathError); ok {
        err = pe.Err
    }
    return err == syscall.ENOENT || err == ErrNotExist
}
```

## 7-12 - インタフェース型アサーションによる振る舞いの問い合わせ

io.WriterのWriteメソッドは型変換時にメモリを割り当ててコピーしている。

```go
func writeHeader(w io.Writer, contentType string) error {
    if _, err := w.Write([]byte("Content-Type: ")); err != nil {
        return err
    }
    if _, err := w.Write([]byte(contentType)); err != nil {
        return err
    }
    // ...
}
```

これを一時コピーのメモリ割り当てをせずに処理することを考える。一時コピーしないWriteStringメソッドを持つインタフェースを定義してインタフェース型アサーションを使えば、WriteStringメソッドを持つ型の場合にこのメソッドを使うように処理を書くことができる。

```go
func writeString(w io.Writer, s string) (n int, err error) {
    type stringWriter interface {
        WriteString(string) (n int, err error)
    }
    if sw, ok := w.(stringWriter); ok {
        return sw.WriteString(s) // WriteStringメソッドが使える場合はこちら
    }
    return w.Write([]byte(s))
}
```

標準ライブラリにもある。
https://golang.org/pkg/io/#WriteString

## 7-13 - 型 switch

型switch文型アサーションのif-elseの連なりをシンプルに記述できる

```go
switch x.(type) {
    case nil:
    case int, uint:
    case bool:
    case string:
    default:
}
```

インタフェースは2つのスタイルで使われる。

- メソッドが強調される
  - インタフェースのメソッドは、そのインタフェースを満足する具象型の類似性を表現する
- インタフェースを満足する具象型が強調される
  - インタフェース値がさまざまな具象型の値を保持できる → 型の共用体
 
## 7-14 - トークンに基づくXMLデコード

encoding/xml パッケージはXMLをデコードするための低レベルのトークンに基づくAPIを提供している。トークンは

- StartElement
- EndElement
- CharData
- Comment

の４種類。

```go
// 判別共用体のインタフェース
type Token interface{}

type StartElement struct {
    Name Name
    Attr []Attr
}

type EndElement struct {
    Name Name
}

type CharData []byte 

type Comment []byte

type Decoder struct {
    /*...*/
}
```

判別共用体を満足する具象型の集まりは意図的に固定されている（？）

## 7-15 - ちょっとした助言


