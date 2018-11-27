package main

import (
	"fmt"
	"time"
)

func main() {

	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("Hello i", i)
		}()



	}
	//fmt.Println(i)
	time.Sleep(time.Hour)
}
