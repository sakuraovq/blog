package client

import (
	"fmt"
	"learn/crawler/engine"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/worker"
	"net/rpc"
)

func CreateWorkProcessor(cliChan chan *rpc.Client) engine.Processor {

	return func(request engine.Request) (engine.ParserResult, error) {
		// to rpc client connected pool get one client
		serializedRequest := worker.SerializedRequest(request)
		var serializedResult worker.ParseResult

		cli := <-cliChan
		err := cli.Call(config.CrawService,
			serializedRequest, &serializedResult)
		if err != nil {
			fmt.Println(err)
			return engine.ParserResult{}, err
		}
		return worker.UnSerializedResult(serializedResult), nil
	}

}
