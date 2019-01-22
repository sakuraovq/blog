package main

import (
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/rpcsupport"
	"learn/crawler_distributed/worker"
	"log"
	"testing"
	"time"
)

func TestWorkRpc(t *testing.T) {
	const host = ":1235"
	go func() {
		log.Fatal(rpcsupport.RpcServer(host, worker.CrawService{}))
	}()
	time.Sleep(time.Second)

	client, e := rpcsupport.NewClient(host)

	if e != nil {
		t.Errorf("init craw rpc client fail %v", e)
	}
	req := worker.Request{
		Url: "http://album.zhenai.com/u/1437200884",
		Parser: worker.SerializedParser{
			FuncName: config.ParserProfile,
			Args:     "测试下下",
		},
	}
	var result worker.ParseResult
	e = client.Call(config.CrawService, req, &result)
	if e != nil {
		t.Errorf("call CrawService Fail %v", e)
	}
	t.Log("success", result)
}
