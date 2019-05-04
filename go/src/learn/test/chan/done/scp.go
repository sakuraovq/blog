package main

import (
	"fmt"
	"sync"
)

type worker struct {
	in   chan int
	done func()
}

func main() {
	chanDemo()
}

func doWorker(id int, ch worker) {
	go func() {
		for c := range ch.in {
			fmt.Printf("worker recevice id=%d chan=%c\n", id, c)
			ch.done()
		}
	}()
}

func createWorker(id int, wg *sync.WaitGroup) worker {
	w := worker{
		in: make(chan int),
		done: func() {
			wg.Done()
		},
	}

	go doWorker(id, w)
	return w
}

func chanDemo() {
	var waitGroup sync.WaitGroup
	var channels [10] worker
	for i := 0; i < 10; i++ {
		channels[i] = createWorker(i, &waitGroup)
	}

	for i, worker := range channels {
		waitGroup.Add(1)
		worker.in <- 'a' + i
	}
	for i, worker := range channels {
		waitGroup.Add(1)
		worker.in <- 'A' + i
	}

	waitGroup.Wait()

}
