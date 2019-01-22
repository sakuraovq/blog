package client

import (
	"fmt"
	"learn/crawler/engine"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/rpcsupport"
	"log"
)

func GetItemSaver() (chan engine.Item, error) {

	client, e := rpcsupport.NewClient(
		fmt.Sprintf(":%d", config.SeverPort))
	if e != nil {
		return nil, e
	}
	saver := make(chan engine.Item)

	//var itemCount int32

	for i := 0; i < config.SaverGoCount; i++ {
		go func() {
			for {
				item := <-saver
				result := ""
				err := client.Call(config.SaverService, item, &result)

				if err != nil {
					log.Printf("item error %v", err)
				}
				//atomic.AddInt32(&itemCount, 1)
				//log.Printf("Got count #%d item %+v saved result %s", itemCount, item, result)
			}
			client.Close()
		}()
	}
	return saver, nil
}
