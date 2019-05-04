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



# 7.2 - インタフェース