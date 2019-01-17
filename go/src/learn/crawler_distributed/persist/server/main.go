package main

import (
	"gopkg.in/olivere/elastic.v3"
	"learn/crawler_distributed/rpcsupport"
	"log"
)

func main() {

	log.Fatal(RpcServer(":1234", "dating_profile"))
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
