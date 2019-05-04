package intset

import (
	"bytes"
	"fmt"
)

// IntSet - xx
type IntSet struct {
	words []uint64
}

// Has - xx
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add - xx
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// AddAll - 値のリストを追加する
func (s *IntSet) AddAll(list ...int) {
	for _, val := range list {
		s.Add(val)
	}
}

// UnionWith - sとtの和集合をsに設定する
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len - 要素数を返す
func (s *IntSet) Len() int {
	count := 0
	for _, word := range s.words {
		for word != 0 {
			count++
			word &= word - 1
		}
	}
	return count
}

// Remove - セットから x を取り除く
func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		return
	}
	s.words[word] &= ^(1 << bit)
}

// Clear - セットから全て尾の要素を取り除きます
func (s *IntSet) Clear() {
	s.words = []uint64{}
}

// Copy - セットのコピーを返します
func (s *IntSet) Copy() *IntSet {
	var c IntSet
	c.words = append(c.words, s.words...)
	return &c
}
