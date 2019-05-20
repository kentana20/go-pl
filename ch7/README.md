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



## 7-9 - 式評価器



## 7-10 - 


## 7-11 - 


## 7-12 - 


## 7-13 - 


## 7-14 - 


