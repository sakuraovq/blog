package engine

import "log"

// Request 调度器
type Scheduler interface {
	Submit(Request)
	GetWorkerChan() chan Request
	Run()
	WorkerNotifier
}

// worker 通知接口
type WorkerNotifier interface {
	WorkerReady(chan Request)
}

type ConcurrentEngine struct {
	Scheduler Scheduler
	WorkCount int
}

func (e ConcurrentEngine) Run(seed ...Request) {

	out := make(chan ParserResult)
	// 运行调度器
	e.Scheduler.Run()
	// 初始化work
	for i := 0; i < e.WorkCount; i++ {
		createWorker(e.Scheduler.GetWorkerChan(), out, e.Scheduler)
	}
	// 初始化种子
	for _, req := range seed {
		e.Scheduler.Submit(req)
	}

	gotCount := 0
	for {
		// 接收 parerResult
		parserResult := <-out
		for _, item := range parserResult.Items {
			log.Printf("Got item #%d val %v", gotCount, item)
			gotCount++
		}
		// 提交请求
		for _, req := range parserResult.Request {
			e.Scheduler.Submit(req)
		}
	}
}

func createWorker(in chan Request, out chan ParserResult, notify WorkerNotifier) {

	go func() {
		for {
			// tell worker is ready
			notify.WorkerReady(in)
			req := <-in
			result, e := worker(req)
			if e != nil {
				continue
			}
			out <- result
		}
	}()
}
