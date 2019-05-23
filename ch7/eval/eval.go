package eval

import (
	"fmt"
	"math"
)

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
    Check(vars map[Var]bool) error
}

//Eval - 環境の検索
func (v Var) Eval(env Env) float64 {
	return env[v]
}

//Eval - リテラル値の返却
func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

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
