package main

import (
	"flag"
	"fmt"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/rpcsupport"
	"learn/crawler_distributed/worker"
	"log"
)

var workerPort = flag.Int("port",
	config.WorkPort0,
	"worker port default "+string(config.WorkPort0))

func main() {
	flag.Parse()

	log.Fatal(rpcsupport.RpcServer(
		fmt.Sprintf(":%d", *workerPort), worker.CrawService{}))
}
