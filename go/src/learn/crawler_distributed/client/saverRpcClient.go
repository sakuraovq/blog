package client

import (
	"fmt"
	"learn/crawler/engine"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/rpcsupport"
	"log"
)

func GetItemSaver(esPort int) (chan engine.Item, error) {

	client, e := rpcsupport.NewClient(
		fmt.Sprintf(":%d", esPort))
	if e != nil {
		return nil, e
	}
	saver := make(chan engine.Item)

	for i := 0; i < config.SaverGoCount; i++ {
		go func() {
			for {
				item := <-saver
				// 不阻塞 异步saver
				go func() {
					result := ""
					err := client.Call(config.SaverService, item, &result)

					if err != nil {
						log.Printf("item error %v", err)
					}
				}()
			}		
		}()
	}
	return saver, nil
}
