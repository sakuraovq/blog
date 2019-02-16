package main

import (
	"learn/crawler/engine"
	"learn/crawler/model"
	"log"
	"net"
	"net/rpc/jsonrpc"
	"testing"
	"time"
)

func TestSever(t *testing.T) {
	const host = ":1234"

	// start sever
	go func() {
		error := RpcServer(host, "test")
		if error != nil {
			log.Fatal(error)
		}
	}()
	// wait server started
	time.Sleep(time.Second)

	testItem := engine.Item{
		Url:  "sakuraus.cn",
		Id:   "1",
		Type: "zhenai",
		Payload: model.Profile{
			Name: "test",
		},
	}
	conn, e := net.Dial("tcp", host)
	if e != nil {
		log.Fatal(e)
	}
	client := jsonrpc.NewClient(conn)

	var resultStr string
	e = client.Call("ItemSaverService.Saver", testItem, &resultStr)
	if e != nil {
		t.Errorf("call ItemSaverService.Saver error is %v", e)
	}
	if resultStr != "ok" {
		t.Errorf("call ItemSaverService.Saver Result not expected %s", resultStr)
	}
}
