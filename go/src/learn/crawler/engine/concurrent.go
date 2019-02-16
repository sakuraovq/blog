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
	Scheduler        Scheduler
	WorkCount        int
	ItemSaver        chan Item
	RequestProcessor Processor
}

type Processor func(Request) (ParserResult, error)

func (e ConcurrentEngine) Run(seed ...Request) {

	out := make(chan ParserResult)
	// 运行调度器
	e.Scheduler.Run()
	// 初始化work
	for i := 0; i < e.WorkCount; i++ {
		e.createWorker(e.Scheduler.GetWorkerChan(), out, e.Scheduler)
	}
	// 初始化种子
	for _, req := range seed {
		e.Scheduler.Submit(req)
	}

	for {
		// 接收 parerResult
		parserResult := <-out
		for _, item := range parserResult.Items {
			// 消耗小 不必用队列
			go func() { e.ItemSaver <- item }()
		}
		// 提交请求
		for _, req := range parserResult.Request {
			if checkDuplicate(req.Url) {
				e.Scheduler.Submit(req)
			}
		}
	}
}

func (e ConcurrentEngine) createWorker(in chan Request, out chan ParserResult, notify WorkerNotifier) {

	go func() {
		for {
			// tell worker is ready
			notify.WorkerReady(in)
			req := <-in
			result, e := e.RequestProcessor(req)
			if e != nil {
				log.Print("work error ", e)
				continue
			}
			out <- result
		}
	}()
}

var visitsMaps = make(map[string]bool)

// 验证url是否重复
func checkDuplicate(url string) bool {
	// 不存在会返回 false
	if visitsMaps[url] {
		return false
	}
	visitsMaps[url] = true
	return true
}
