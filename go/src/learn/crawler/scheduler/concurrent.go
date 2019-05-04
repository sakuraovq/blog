package scheduler

import "learn/crawler/engine"

type ConcurrentScheduler struct {
	WorkerChan chan engine.Request // 所有worker 抢一个 chan
}

func (e *ConcurrentScheduler) GetWorkerChan() chan engine.Request {
	return e.WorkerChan
}

func (e *ConcurrentScheduler) WorkerReady(chan engine.Request) {
}

func (e *ConcurrentScheduler) Run() {
	e.WorkerChan = make(chan engine.Request)
}

// 多个go 提交到 workerChan 然后去抢 workerChan
func (e *ConcurrentScheduler) Submit(req engine.Request) {
	go func() { e.WorkerChan <- req }()
}
