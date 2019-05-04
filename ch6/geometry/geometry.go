package main

import (
	"fmt"
	"math"
)

// Point - 構造体
type Point struct {
	X, Y float64
}

// Path - 名前付きスライス
type Path []Point

// Distance - 普通の関数
func Distance(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// Distance - pはレシーバ、Point型のメソッド
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// Distance - スライスに対するメソッド
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}

func main() {
	a := Point{1.0, 2.0}
	b := Point{10.0, 20.0}

	// 関数呼び出し
	fmt.Printf("distance is %f\n", Distance(a, b))

	// メソッド呼び出し
	fmt.Printf("distance#2 is %f\n", a.Distance(b))

	// 異なる名前空間のメソッド
	perim := Path{
		{1, 1},
		{5, 1},
		{5, 4},
		{1, 1},
	}
	fmt.Printf("distance#3 is %f\n", perim.Distance())
}
