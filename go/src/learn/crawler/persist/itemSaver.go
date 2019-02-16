package persist

import (
	"encoding/json"
	"errors"
	"gopkg.in/olivere/elastic.v3"
	"learn/crawler/engine"
	"log"
)

func GetItemSaver(index string) (chan engine.Item, error) {

	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	saver := make(chan engine.Item)

	go func() {
		itemCount := 0
		for {
			item := <-saver
			itemCount++
			log.Printf("Got count #%d item %+v ", itemCount, item)
			// TODO: need start elastic search server
			err := Save(item, client, index)
			if err != nil {
				log.Printf("item error %v", err)
			}
		}
	}()

	return saver, nil
}

func Save(item engine.Item, client *elastic.Client, index string) error {

	if item.Type == "" {
		return errors.New("Must exits Type !")
	}

	indexService := client.Index()
	indexService.Index(index).Type(item.Type)

	if item.Id != "" {
		indexService.Id(item.Id)
	}

	// es 3.0下 传递字符串靠谱
	bytes, _ := json.Marshal(item)

	_, err := indexService.BodyString(string(bytes)).Do()

	if err != nil {
		return err
	}

	return nil
}
