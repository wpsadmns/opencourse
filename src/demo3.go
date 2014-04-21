package main

import (
	"fmt"
)

func main() {
	input1 := "input1"
	input2 := "input2"
	r1 := []rune(input1)
	r2 := []rune(input2)
	r1 = append(r1, r2[0:len(r2)]...)
	fmt.Println(string(r1))
}
