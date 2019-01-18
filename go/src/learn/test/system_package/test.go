package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {

	str := `aaa"b"
1
	12qwer
	qwer`
	reader := strings.NewReader(str)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan(){
		fmt.Println(scanner.Text())
	}
}
