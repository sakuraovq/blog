package persist

import (
	"errors"
	"gopkg.in/olivere/elastic.v3"
	"learn/crawler/engine"
	"log"
)

func GetItemSaver() chan engine.Item {
	saver := make(chan engine.Item)

	go func() {
		itemCount := 0
		for {
			item := <-saver
			itemCount++
			log.Printf("Got count #%d item %+v ", itemCount, item)
			// TODO: need start elastic search server
			err := save(item)
			if err != nil {
				log.Printf("item error %v", err)
			}
		}
	}()
	return saver
}

const profileDatabase = "dating_profile"

func save(item engine.Item) error {
	client, err := elastic.NewClient(
		elastic.SetSniff(false))

	if err != nil {
		return err
	}
	if item.Type == "" {
		return  errors.New("Must exits Type !")
	}

	indexService := client.Index()
	indexService.Index(profileDatabase).Type(item.Type)

	if item.Id != "" {
		indexService.Id(item.Id)
	}
	_, err = indexService.BodyJson(item.Payload).Do()

	if err != nil {
		return err
	}

	return nil
}
