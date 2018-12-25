package scheduler

import "learn/crawler/engine"

type QueuedScheduler struct {
	requestChan chan engine.Request
	workerChan  chan chan engine.Request  // 每个 worker 创建一个 chan Request
}

// 提交一个 Request
func (queued *QueuedScheduler) Submit(req engine.Request) {
	queued.requestChan <- req
}

// 告诉 worker 已经 准备好了
func (queued *QueuedScheduler) WorkerReady(worker chan engine.Request) {
	queued.workerChan <- worker
}

// 获取 work chan 没一个worker, 分配一个 chan Request
func (queued *QueuedScheduler) GetWorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

func (queued *QueuedScheduler) Run() {
	// 初始化 chan
	queued.workerChan = make(chan chan engine.Request)
	queued.requestChan = make(chan engine.Request)
	go func() {
		var requestQueue []engine.Request
		var workerQueue []chan engine.Request

		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request

			// 有一个request 和一个 worker chan 才进行投递request
			if len(requestQueue) > 0 && len(workerQueue) > 0 {
				activeRequest = requestQueue[0]
				activeWorker = workerQueue[0]
			}
			select {
			case req := <-queued.requestChan: // 把 submit的request 投递到队列中等待有空闲的worker在消费
				requestQueue = append(requestQueue, req)
			case worker := <-queued.workerChan: // 把 worker 放在 worker队列 等待 request 一起消费
				workerQueue = append(workerQueue, worker)
			case activeWorker <- activeRequest: // 只有有活动的req和worker才能case 到, 默认定义activeWorker为nil
				requestQueue = requestQueue[1:] // 去除以投递的请求
				workerQueue = workerQueue[1:]
			}
		}
	}()
}
