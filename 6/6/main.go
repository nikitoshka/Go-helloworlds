package main

import (
	"fmt"
	"the_go_progr_lang/6_6/point"
)

func main() {
	p := point.Point{}
	p.Set(2, 4)

	q := point.Point{}
	q.Set(1, 1)

	fmt.Println(p)

	fmt.Println(p.Difference(q))
	fmt.Println(point.Difference(p, q))
}
