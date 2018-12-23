package engine

import "log"

type Scheduler interface {
	Submit(Request)
	SendSchedulerChannelRequest(chan Request)
}

type ConcurrentEngine struct {
	Scheduler Scheduler
	WorkCount int
}

func NewConcurrentEngine(scheduler Scheduler, workCount int) ConcurrentEngine {
	return ConcurrentEngine{Scheduler: scheduler, WorkCount: workCount}
}

func (e ConcurrentEngine) Run(seed ...Request) {

	in := make(chan Request)
	out := make(chan ParserResult)
	// 设置Scheduler 工作chan 为 in 每次submit 到in chan 中
	e.Scheduler.SendSchedulerChannelRequest(in)
	// 初始化work
	for i := 0; i < e.WorkCount; i++ {
		createWorker(in, out)
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

func createWorker(in chan Request, out chan ParserResult) {
	go func() {
		for {
			req := <-in
			result, e := worker(req)
			if e != nil {
				continue
			}
			out <- result
		}
	}()
}
