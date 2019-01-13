package persist

import (
	"encoding/json"
	"gopkg.in/olivere/elastic.v3"
	"learn/crawler/engine"
	"learn/crawler/model"
	"testing"
)

func TestSaver(t *testing.T) {
	testItem := engine.Item{
		Url:  "sakuraus.cn",
		Id:   "1",
		Type: "zhenai",
		Payload: model.Profile{
			Name: "test",
		},
	}
	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	const index = "test_profile"

	e := save(testItem, client, index)
	if e != nil {
		panic(e)
	}

	result, err := client.Get().
		Index(index).
		Type(testItem.Type).
		Id(testItem.Id).
		Do()

	if err != nil {
		panic(err)
	}
	var resultItem engine.Item

	err = json.Unmarshal(*result.Source, &resultItem)

	if err != nil {
		panic(err)
	}
	if testItem != resultItem {
		t.Errorf("elasticsearch data %+v expect %+v", testItem, resultItem)
	}
}
