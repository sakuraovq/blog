package persist

import (
	"encoding/json"
	"gopkg.in/olivere/elastic.v3"
	"learn/crawler/model"
	"testing"
)

func TestSaver(t *testing.T) {
	testItem := model.Profile{
		Name: "test",
	}
	id, e := save(testItem)
	if e != nil {
		panic(e)
	}

	client, err := elastic.NewClient(
		elastic.SetSniff(false))

	if err != nil {
		panic(err)
	}

	result, err := client.Get().
		Index(profileDatabase).
		Type(profileTable).
		Id(id).
		Do()

	if err != nil {
		panic(err)
	}
	var resultItem model.Profile

	err = json.Unmarshal(*result.Source, &resultItem)

	if err != nil {
		panic(err)
	}
	if testItem != resultItem {
		t.Errorf("elasticsearch data %+v expect %+v", testItem, resultItem)
	}
}
