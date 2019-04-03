package main

import "fmt"

func main() {
	data := []string{"hoge", "", "fuga", "", "piyo", ""}
	fmt.Printf("nonempty: %s\n", nonempty(data))

	data2 := []string{"hoge", "", "fuga", "", "piyo", ""}
	fmt.Printf("nonempty2: %s\n", nonempty2(data2))
}

func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	fmt.Println(strings) // [hoge fuga piyo "" piyo ""]
	return strings[:i]
}

func nonempty2(strings []string) []string {
	out := strings[:0]

	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
