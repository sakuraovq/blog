package main

import "fmt"

func main() {

	var a, b = 3, 5

	defer func() {
		z := a + b
		fmt.Print(z)
	}()
	a = a + b
	b = a - b
}
