package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func min(params ...int) int {
	min := 100

	for _, param := range params {
		if param < min {
			min = param
		}
	}

	return min
}

func max(params ...int) int {
	max := 0

	for _, param := range params {
		if param > max {
			max = param
		}
	}

	return max
}

func atoiM(str []string) []int {
	var ints []int

	for _, s := range str {
		if i, err := strconv.Atoi(s); err == nil {
			ints = append(ints, i)
		}
	}

	return ints
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("wrong number of arguments")
	}

	ints := atoiM(os.Args[2:])
	if len(ints) == 0 {
		log.Fatalln("failed to get intes from the arguments")
	}

	switch os.Args[1] {
	case "-min":
		fmt.Println(min(ints...))

	case "-max":
		fmt.Println(max(ints...))
	}
}
