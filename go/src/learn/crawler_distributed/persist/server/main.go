package main

import (
	"flag"
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/rpcsupport"
	"log"
)

var esPort = flag.Int("es_port",
	config.SeverPort,
	"elastic search port default "+string(config.SeverPort))

func main() {
	flag.Parse()

	log.Fatal(RpcServer(
		fmt.Sprintf(":%d", *esPort),
		config.ElasticIndex))
}

func RpcServer(host, index string) error {

	client, e := elastic.NewClient(
		elastic.SetSniff(false))
	if e != nil {
		return e
	}

	itemSaverService := rpcsupport.ItemSaverService{
		Client: client,
		Index:  index,
	}
	return rpcsupport.RpcServer(host, itemSaverService)
}
