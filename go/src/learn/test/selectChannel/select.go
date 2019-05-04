package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(
				time.Duration(rand.Intn(1500)) * time.Millisecond)
			out <- i
			i++
		}
	}()
	return out
}

func createWorker(id int) chan<- int {
	c := make(chan int)
	go func() {
		for ch := range c {
			time.Sleep(1 * time.Second)
			fmt.Printf("worker recevice id=%d chan=%d\n", id, ch)
		}
	}()

	return c
}

func main() {
	var values []int // 排队消费
	c1, c2 := generator(), generator()
	w := createWorker(0)

	timer := time.After(10 * time.Second) // 10秒后执行一次

	tick := time.Tick(time.Second) // 每秒执行定时器
	for {
		var (
			// (如果数据没有准备好,就用)nil chan 是阻塞等待的就不会被select 调度
			activeWorker chan<- int
			activeValue  int
		)
		if len(values) > 0 {
			activeWorker = w
			activeValue = values[0]
		}
		select {
		case n := <-c1:
			values = append(values, n)
		case n := <-c2:
			values = append(values, n)
		case activeWorker <- activeValue:
			values = values[1:]
		case <-tick:
			fmt.Println("values lens ", len(values))
		case <-time.After(800 * time.Millisecond):
			fmt.Println("timeout")
		case <-timer:
			fmt.Println("Bye")
			return
		}
	}
}
