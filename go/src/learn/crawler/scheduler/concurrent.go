package scheduler

import "learn/crawler/engine"

type ConcurrentScheduler struct {
	WorkerChan chan engine.Request
}

// 并发提交
func (e *ConcurrentScheduler) Submit(req engine.Request) {
	go func() { e.WorkerChan <- req }()
}

func (e *ConcurrentScheduler) SendSchedulerChannelRequest(c chan engine.Request) {
	e.WorkerChan = c
}
