package main

import (
	"fmt"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/rpcsupport"
	"learn/crawler_distributed/worker"
	"log"
)

func main() {

	log.Fatal(rpcsupport.RpcServer(
		fmt.Sprintf(":%d", config.WorkPort0), worker.CrawService{}))
}
