7章 - インタフェース
---

インタフェース型はほかの型のふるまいに関する一般化あるいは抽象化を表す。多くのGoプログラムは自分で定義したインタフェースと同じくらい多くの標準インタフェースを使う。

# 7.1 - 契約としてのインタフェース

インタフェースは抽象型。内部構造を公開せず、メソッドのいくつかを公開しているのみ。インタフェース型の値がある場合、その値が何であるかはわからず、何ができるかを知っているだけ。

- fmt.Printf
- fmt.Sprintf

での例を見てみる。

```go
package fmt

func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)

func Printf(format string, args ...interface{}) (int, error) {
    return Fprintf(os.Stdout, format, args...)
}

func Sprintf(format string, args ...interface{}) string {
    var buf bytes.Buffer
    Fprintf(&buf, format, args...)
    return buf.String()
}
```

Fprintfの第一引数はio.Writerというインタフェース型。渡されたファイルへ書き出される。 `Printf` は引数に os.Stdout を、`Sprintf`は引数にメモリバッファへのポインタを渡している。

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

io.Writer インタフェースによってFprintfと呼び出し元の契約を定義している。Writeメソッドを持った具象型の値を渡せばFprintfは機能する。同じインタフェースを満足する別の型で置き換え可能 → 代替可能性がある。

```go
type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // intをByteCounterへ変換
	return len(p), nil
}
```

ByteCounter型にWriteメソッドを定義している → Fprintfへ渡せる。

# 7.2 - インタフェース型

インタフェース型は具象型がそのインタフェースのインスタンスとしてみなされるために持たなければならないメソッドの集まりを定義する。
既存の型の組み合わせで新たなインタフェース型を宣言することもできる（インタフェースの埋め込み）。

```go
package io

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Reader interface {
    Read(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// 埋め込みを使ったインタフェース定義
type ReadWriter interface {
    Reader
    Writer
}

// ↑と一緒
type ReadWriter interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```
 
 # 7-3 - インタフェースを満足する

インタフェースが要求しているすべてのメソッドを型が保持していれば、その方はインタフェースを満足する。

- ex. *os.File は以下を満足する
  - io.Reader
  - io.Writer
  - io.Closer
  - io.ReadWriter

ある具象型が特定のインタフェースを満足しているという意味でその具象型がそのインタフェースであると表現することが多い。

- ex. *bytes.Buffer は io.Writer / *os.File は io.ReadWriter である

その型がインタフェースを満足していれば、そのインタフェースへ代入できる。

```go
var w io.Writer
w = os.Stdout
w = new(bytes.Buffer)
w = time.Second // コンパイルエラー
```

インタフェースは具象型とその型が保持する値を包んで隠すので、インタフェース型が公開しているメソッドしか呼び出せない。

```go
os.Stdout.Write([]byte("hello"))
os.Stdout.Close()

var w io.Writer
w = os.Stdout
w.Write([]byte("hello"))
w.Close() // コンパイルエラー（io.WriterインタフェースにCloseメソッドがないため）
```

空インタフェース型はそれを満足する方に対して何も要求していないため、すべての値を空インタフェースへ代入できる。

```go
var any interface{}
any = true
any = 123
any = "hello"
any = map[string]int{"one": 1}
any = new(bytes.Buffer)
```

具象型は関連のない多くのインタフェースを満足できる。

- ex. デジタルコンテンツを体系化するプログラム

```go
// Album, Book, Movie, Magazine, Podcast, TVEpisode, Track

type Artifact interface {
    Title() string
    Creators() []string
    Created() time.Time
}

type Text interface {
    Pages() int
    Words() int
    PageSize() int
}

type Audio interface {
    Stream() (io.ReadCloser, error)
    RunningTime() time.Duration
    Format() string
}

type Video interface {
    Stream() (io.ReadCloser, error)
    RunningTime() time.Duration
    Format() string
    Resolution() (x, y int)
}

type Streamer interface {
    Stream() (io.ReadCloser, error)
    RunningTime() time.Duration
    Format() string
}

type Audio interface {
    Streamer
}

type Video interface {
    Streamer
    Resolution() (x, y int)
}
```

 # 7-4 - flag.Value によるフラグの解析

指定された期間だけスリープするプログラム

```go
// デフォルトは1sec
var period = flag.Duration("period", 1*time.Second, "sleep period")

func main() {
	flag.Parse()
	fmt.Printf("Sleeping for %v...", *period)
	time.Sleep(*period)
	fmt.Println()
}
```

`flag.Duration` 関数は `time.Duration` 型のフラグ変数を生成する → ユーザが使いやすい形式で期間を指定可能になる。

```go
package flag

type Value interface {
    String() string
    Set(string) error
}
```

 flag.Value インタフェースを満足する型でフラグ表記を定義することもできる。

 ---

 気温を摂氏、華氏で指定することができるcelsiusFlag型を定義してみる。

 ```go
 package tempconv

import (
	"flag"
	"fmt"
)

type Celsius float64
type Fahrenheit float64

// celsiusFlag型
type celsiusFlag struct{ Celsius }

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

func (c Celsius) String() string {
	return fmt.Sprintf("%g ℃", c)
}

// FToC 華氏 → 摂氏
func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// celsiusFlag型に対するSetメソッド
func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // no error check needed
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}
```

`flag.CommandLine.Var` でグローバル変数である flag.CommandLine にフラグを追加できる。

# 7-5 - インタフェース値

インタフェースは

- 動的な型（dynamic type）
- 動的な値（dynamic value）

の2つの構成要素を持っている。Goは静的型付言語なので、型はコンパイル時に解釈され、値ではない。

```go
var w io.Writer
w = os.Stdout
w = new(bytes.Buffer)
w = nil
```

はじめの宣言時は type, value ともに nilが入る（nilインタフェース値）。メソッド呼び出しするとnil参照エラー。

```go
w.Write([]byte("hello")) // panic
```

2つめで*os.Fileの値を代入している。具象型からインタフェース型へ暗黙的に変換している。動的な型は *os.File (ポインタ型)で動的な値はos.Fileへの参照となる。

```go
w.Write([]byte("hello")) // "hello"
```

3つめで*bytes.Bufferの値をインタフェース値に代入している。動的な型は *bytes.Bufferで動的な値は bytes.Buffer への参照。

```go
w.Write([]byte("hello")) // bytes.Buffer へ "hello" を書き出す
```

4つめでnilを代入してDynamic Type, Dynamic Valueともにnilにしている。

---

```go
var x interface {} = []int{1, 2, 3}
fmt.Println(x == x)
```

インタフェース値は `==` , `!=` で比較できるが、型が比較可能でない場合はpanicになる。こわい。

# 7-6 - sort.Interface でのソート

ソートをするためには

- 列の長さ
- 要素比較の方法
- 要素を入れ替える方法

が必要。これがそのままインタフェースになる。

```go
package sort

type Interface interface {
    Len()
    Less(i, j int) bool
    Swap(i, j int)
}
```

文字列のスライス names をソートする例。

```go
type StringSlice []string

func (p StringSlice) Len() int { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int) { return p[i], p[j] = p[j], p[i] }

// これで sort.Sort(StringSlice(names)) できる
sort.Sort(StringSlice(names))

// もともとsortパッケージに提供されているStrings関数でもおなじことができる
sort.Strings(names)
```

プレイリストをソートする例

```go
package main

import "time"

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

// Track型のスライスへのポインタ
var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush()
}
```

- Reverse関数
  - Less関数に対する引数を逆にしているInterfaceを持つ構造体の型

```go
package sort

type reverse struct { Interface }
func (r reverse) Less(i, j, int) bool { return r.Interface.Less(j, i) }
func Reverse(data Interface) Interface { return reverse{data} }
```

- ソート条件を複数指定する

```go
type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

sort.Sort(customSort{tracks, func(x, y *Track) bool {
	if x.Title != y.Title {
		return x.Title < y.Title
	}
	if x.Year != y.Year {
		return x.Year < y.Year
	}
	if x.Length != y.Length {
		return x.Length < y.Length
	}
	return false
}})
```

- ソート済かを判定する

```go
values := []int{3, 1, 4, 1}
fmt.Println(sort.IntsAreSorted(values)) // ソートされてないのでfalse
sort.Ints(values) // 整数昇順にソート
fmt.Println(values) // "[1 1 3 4]"
fmt.Println(sort.IntsAreSorted(values)) // ソート済なのでtrue
sort.Sort(sort.Reverse(sort.IntSlice(values))) // Reverse関数で逆順に
fmt.Println(values) // "[4 3 1 1]"
fmt.Println(sort.IntsAreSorted(values)) // ソートされてない？結果はfalse
```
