package main

import (
	"fmt"
	"sync"
	"time"
)

type atomic struct {
	lock  sync.Mutex
	value int
}

func (atomic *atomic) inc() {
	fmt.Println("safe")
	func() {  // 部分锁
		atomic.lock.Lock()
		defer atomic.lock.Unlock()
		atomic.value++
	}()

}

func (atomic *atomic) get() int {

	atomic.lock.Lock()
	defer atomic.lock.Unlock()
	return atomic.value
}

func main() {

	atomic := atomic{}
	go func() {
		atomic.inc()
	}()

	go func() {
		atomic.inc()
	}()
	time.Sleep(time.Second)
	fmt.Println(atomic.get())
}
