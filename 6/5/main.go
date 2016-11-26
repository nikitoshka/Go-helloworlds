package main

import (
	"fmt"

	"./intSet"
)

func main() {
	var s intSet.IntSet

	s.Add(10)
	s.Add(5)

	fmt.Println(&s)
	fmt.Println(s.Len())
	fmt.Println(s.Has(10))
	fmt.Println(s.Has(8))

	// fmt.Println()

	// s.Remove(5)
	// fmt.Println(&s)
	// fmt.Println(s.Len())
	// fmt.Println(s.Has(10))
	// fmt.Println(s.Has(5))

	// fmt.Println()

	// s.Clear()
	// fmt.Println(&s)
	// fmt.Println(s.Len())
	// fmt.Println(s.Has(10))
	// fmt.Println(s.Has(5))

	t := s.Copy()

	fmt.Println()

	fmt.Println(t)
	fmt.Println(t.Len())
	fmt.Println(t.Has(10))
	fmt.Println(t.Has(5))

	t.Clear()

	fmt.Println()

	fmt.Println(t)
	fmt.Println(t.Len())

	fmt.Println()

	fmt.Println(&s)
	fmt.Println(s.Len())
}
