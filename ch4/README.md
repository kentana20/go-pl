第四章: コンポジット型
---

基本型を組み合わせることで作成できるコンポジット型。

- 配列
- スライス
- マップ
- 構造体

の4種。配列と構造体は合成（Aggregate）型。メモリ内で別の値を連結したもの。

- 配列内の要素はすべて同じ型になるが、構造体は不均質
- 配列、構造体は固定サイズだが、スライス、マップは動的サイズ

## 4.1 配列

特定の型の0個以上の固定長列。Goでは配列が直接使われることは稀で、スライスを使うことが多いが、スライスを理解するためには配列の理解が必要。
配列内の個々の要素はインデックス表記でアクセスでき、インデックスは0からはじまる。

```go
var a [3]int
fmt.Println(a[0]) // 要素の1つ目
fmt.Println(a[len(a)-1]) // 最後の要素

// 配列リテラルで初期化できる
var q[3]int = [3]int{1,2,3}
var r[3]int = [3]int{1,2} // 3つ目の要素は0になる
```

配列は大きさを含めて型になる。

```go
q := [3]int{0,1,2}
q = [4]int{0,1,2,3} // コンパイルエラー
```

インデックスと値の組のリストを指定することも可能。

```go
type Currency int

const (
    USD Currency = iota
    EUR
    GBP
    RMB
)

symbol := [...]string{USD: "$", EUR: "€", GBP: "£", RMB: "¥"}
fmt.Println(RMB, symbol[RMB]) // 3 ¥
```

インデックスを省略しても良いが、指定されていない値は要素型のゼロ値になる。

```go
r := [...]int{99: -1} // 100個の要素のうち、最後の要素は-1でそれ以外はすべて0になる
```

- 比較

配列内の要素型が比較可能であれば、配列自体も比較可能となる。

```go
a := [2]int{1, 2}
b := [...]int{1, 2}
c := [2]int{1, 3}

fmt.Println(a == b, a == c, b == c) // true false true

d := [3]int{1, 2}

fmt.Println(a == d) // 配列の大きさは型の一部なので比較できずにコンパイルエラーになる
```

## 4.2 スライス

すべての要素が同じ型である**可変長列**。配列の要素の部分列へアクセスするデータ構造を取る。

- ポインタ: スライスを通して到達可能な配列の最初の要素
- 長さ: スライスの要素数。len関数で取得できる。
- 容量: スライスの開始から基底配列の終わりまでの要素数。cap関数で取得できる。

の3つの構成要素を持つ。

スライス演算子 s[i:j] は列sのiからj-1までの要素を参照する新たなスライスを生成する。
cap(s) を超えてスライスを作成するとpanicになるが、len(s) を超えてスライスを作成することができる。

```go
months := [...]string{
    1: "January",
    2: "February",
    3: "March",
    4:"April",
    5: "May",
    6: "June",
    7: "July",
    8: "August",
    9: "September",
    10: "October",
    11: "November",
    12: "December"}

Q2 := months[4:7] // len:3 (4, 5, 6) / cap:9
summer := months[6:9] // len:3 (6, 7, 8) / cap: 7

fmt.Println(Q2)
fmt.Println(summer)

for _, s := range summer {
    for _, q := range Q2 {
        if s == q {
            fmt.Println("%s appears in both\n", s) // June appears in both
        }
    }
}
```

- スライスリテラル
  - 暗黙的に正しい大きさの配列変数を生成して、その配列を指すスライスをつくる

```go
	s := []int{0, 1, 2, 3, 4, 5} // スライスリテラル
	reverse(s[:2])               // [0, 1] -> [1, 0]
	reverse(s[2:])               // [2, 3, 4, 5] -> [5, 4, 3, 2]
	reverse(s)                   // [2, 3, 4, 5, 0, 1]
    fmt.Println(s)               // [2, 3, 4, 5, 0, 1]
```

- make関数
  - 指定された型、長さ、容量のスライスを作成する

```go
make([]T, len)
make([]T, len, cap)
make([]T, cap[:len])
```

### 4.2.1 append関数

append関数はスライスに項目を追加できる。

```go
var runes []rune
for _, r := range "Hello, 世界" {
    runes = append(runes, r)
}
fmt.Println("%q\n", runes) // "['H', 'e', 'l', 'l', 'o', ',', ' ', '世', '界']"
```

appendの呼び出し結果はappendへ渡された値を持つ同じスライス変数に代入するのが普通（スライス変数の更新）。

```go
var x []int
x = append(x, 1)
x = append(x, 2, 3)
x = append(x, 4, 5, 6)
x = append(x, x...) // xを追加
fmt.Println(x) // [1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6]
```

`...`  ... スライスを追加することができる。

### 4-2-2 スライス内での技法

スタックの実装。

```go
// push
stack = append(stack, v)

// top
top := stack[len(stack) - 1]

// pop
stack = stack[:len(stack) - 1]

// remove
func remove(slice []int, i int) []int {
    copy(slice[i:], slice[i+1:])
    return slice[:len(slice)-1]
}

// remove(順番保証しないバージョン)
func remove(slice []int, i int) []int {
    slice[i] = slice[len(slice) - 1] // 最後の要素を移動
    return slice[:len(slice)-1]
}
```

## 4.3 マップ

ハッシュテーブルは順序付けのないキーと値の組のコレクション。Goのマップはハッシュテーブルへの参照。 `map[K]V` と書く。Kはキーの型、Vは値の型。

| # | # | # |
|---|---|---|
| K | キー | `==` で比較可能でなければならない |
| V | 値 | 制限なし |

```go
ages := map[string]int{
    "alice": 31,
    "charlie": 34
}

delete(ages, "alice") // キー"alice"の要素を削除
delete(ages, "kate") // キー"kate"はマップにないので無視される

ages["bob"] = ages["bob"] + 1 // 1

for age, name := range ages {
    fmt.Printf("%s\t%d\n", name, age)
}

var ages2 map[string]int
ages2["carol"] = 20 // マップが割り当てられていないとpanicになる

age, ok := ages["bob"] // 1, true "bob"が存在しなければokにはfalseが返る
```

`delete` 関数は要素がマップの中になくてもコケない。存在しないキーを指定すると `0` が返る。スライスと同じく `range` に基づく forループを使える。ただし、繰り返し順序は保証されないので、決まった順序で実行するにはキーをソートする必要がある。マップ同士の比較をする場合はスライスと同じくループで書く必要がある。

