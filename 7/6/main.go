package main

import (
	"flag"
	"fmt"

	"./tempFlag"
)

var temp = tempFlag.TempFlag("temp", 25, "set a temperature")

func main() {
	flag.Parse()

	fmt.Println(*temp)
}
