package main

import (
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/rpcsupport"
	"log"
)

func main() {

	log.Fatal(RpcServer(
		fmt.Sprintf(":%d", config.SeverPort),
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
