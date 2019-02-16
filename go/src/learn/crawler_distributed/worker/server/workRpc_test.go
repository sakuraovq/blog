package main

import (
	"learn/crawler/engine"
	"learn/crawler/zhenai/parser"
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
	req := engine.Request{
		Url:   "http://album.zhenai.com/u/1816494626",
		Parse: parser.NewUserProfileParser("男"),
	}

	sReq := worker.SerializedRequest(req)
	var result worker.ParseResult
	e = client.Call(config.CrawService, sReq, &result)
	if e != nil {
		t.Errorf("call CrawService Fail %v", e)
	}
	log.Fatal("result", worker.UnSerializedResult(result))
}
