package main

import (
	"fmt"
	"math"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az, judgea := corner(i+1, j)
			bx, by, bz, judgeb := corner(i, j)
			cx, cy, cz, judgec := corner(i, j+1)
			dx, dy, dz, judged := corner(i+1, j+1)

			if !judgea || !judgeb || !judgec || !judged {
				continue
			}
			fill := fill(az, bz, cz, dz)

			fmt.Printf("<polygon points='%g,%g,%g,%g,%g,%g,%g,%g' fill='%s' />\n", ax, ay, bx, by, cx, cy, dx, dy, fill)
		}
	}
	fmt.Println("</svg>")
}

func fill(a, b, c, d float64) string {
	if ((a + b + c + d) / 4) > 0 {
		return "#ff0000"
	} else {
		return "#0000ff"
	}
}

func corner(i, j int) (float64, float64, float64, bool) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z, judge := f(x, y)

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	return sx, sy, z, judge
}

func f(x, y float64) (float64, bool) {
	r := math.Hypot(x, y)

	if math.IsInf((math.Sin(r) / r), 0) {
		return 0, false
	}

	return math.Sin(r) / r, true
}
