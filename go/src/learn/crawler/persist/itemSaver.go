package persist

import (
	"gopkg.in/olivere/elastic.v3"
	"log"
)

func GetItemSaver() chan interface{} {
	saver := make(chan interface{})

	go func() {
		itemCount := 0
		for {
			item := <-saver
			itemCount++
			log.Printf("Got count #%d item %+v ", itemCount, item)
			// TODO: need start elastic search server
			_, err := save(item)
			if err != nil {
				log.Printf("item error %v", err)
			}
		}
	}()
	return saver
}

const profileDatabase = "dating_profile"
const profileTable = "zhenai"

func save(item interface{}) (id string, err error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false))

	if err != nil {
		return "", err
	}

	response, err := client.Index().
		Index(profileDatabase).
		Type(profileTable).
		BodyJson(item).
		Do()

	if err != nil {
		return "", err
	}

	return response.Id, nil
}
