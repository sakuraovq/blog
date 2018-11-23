package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

func main() {
	//fmt.Println("hello world")
	//consts()
	//enums()
	//readfile()
	//fmt.Println(
	//switchs(10),
	//switchs(70),
	//switchs(619),
	//)

	//readFile("abc.txt")
	//fmt.Println(convertToBin(10), convertToBin(5))

	res := apply(pow, 1, 5)
	fmt.Println(res)
}

func pow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}

func apply(op func(int, int) int, a, b int) int {
	p := reflect.ValueOf(op).Pointer()
	opName := runtime.FuncForPC(p).Name()
	fmt.Println("call func name ", opName, "params", a, b, "res", op(a, b))
	return op(a, b)
}

func readFile(filename string) {
	buf, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scan := bufio.NewScanner(buf)
	for scan.Scan() {
		fmt.Println(scan.Text())
	}
}

func convertToBin(n int) (result string) {
	result = ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(lsb) + result
	}
	return
}

func switchs(score int) string {
	grand := ""
	switch {
	case score < 0 || score > 100:
		panic(fmt.Sprintf("invail score %d", score))
	case score < 60:
		grand = "F"
	case score < 70:
		grand = "b"
	}
	return grand

}

func readfile() {

	const filename = "abc.txt"
	if contenst, err := ioutil.ReadFile(filename); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", contenst)
	}
}

func enums() {
	const (
		golang = iota
		_
		php
		javascript
		c
	)

	const (
		b = 1 << (10 * iota)
		kb
		mb
		gb
		tb
	)

	fmt.Println(b, kb, mb, gb, tb)

	fmt.Println(golang, php, javascript, c)
}

func consts() {
	const filename = "abc.txt"
	const a, b = 3, 4

	fmt.Println(filename, math.Sqrt(a*a+b*b))
}
