package main

import (
	"flag"
	"learn/crawler/engine"
	"learn/crawler/scheduler"
	"learn/crawler/zhenai/parser"
	"learn/crawler_distributed/client"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/rpcsupport"
	"log"
	"net/rpc"
	"strings"
)

var (
	workerHosts = flag.String("worker_hosts", "", "worker hosts, use , split")
	workerCount = flag.Int("worker_count", 15, "run worker count")

	esPort = flag.Int("es_port",
		config.SeverPort,
		"item saver port default"+string(config.SeverPort))
)
// gopm get -g -v golang.org/x/text
func main() {
	flag.Parse()

	if *workerHosts == "" {
		log.Print("must input worker_hosts")
		return
	}

	itemChan, err := client.GetItemSaver(*esPort)
	if err != nil {
		panic(err)
	}

	// 获取worker rpc连接池
	clientChan := getWorkerClientPool(*workerHosts)

	processor := client.CreateWorkProcessor(clientChan)

	engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		ItemSaver:        itemChan,
		WorkCount:        *workerCount,
		RequestProcessor: processor,
	}.Run(engine.Request{
		Url:   config.ZhenAiSeedUrl,
		Parse: engine.NewParserFunc(config.ParserCityList, parser.GetCityList),
	})
}

// 获取worker 连接池
func getWorkerClientPool(workerHosts string) chan *rpc.Client {
	hosts := strings.Split(workerHosts, ",")
	var workerClients []*rpc.Client
	for _, host := range hosts {
		// connect worker rpc client
		newClient, err := rpcsupport.NewClient(host)
		if err != nil {
			log.Printf("connecting worker server %s fail %v", host, err)
			continue
		}
		log.Printf("connected worker server %s success", host)
		workerClients = append(workerClients, newClient)
	}
	// build rpc client connected pool
	workerRpcClientChan := make(chan *rpc.Client)

	go func() {
		for {
			for _, workerClient := range workerClients {
				workerRpcClientChan <- workerClient
			}
		}
	}()

	return workerRpcClientChan
}
