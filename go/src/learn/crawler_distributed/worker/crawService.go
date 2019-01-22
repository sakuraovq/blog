package worker

import "learn/crawler/engine"

type CrawService struct{}

// 传输到engine->work 需要序列化
func (CrawService) Process(request Request, result *ParseResult) error {
	engineRequest, err := UnSerializedRequest(request)
	if err != nil {
		return err
	}
	parserResult, err := engine.Worker(engineRequest)
	if err != nil {
		return err
	}
	*result = SerializedResult(parserResult)
	return nil
}
