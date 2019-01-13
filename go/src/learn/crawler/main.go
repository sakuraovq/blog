package main

import (
	"learn/crawler/engine"
	"learn/crawler/persist"
	"learn/crawler/scheduler"
	"learn/crawler/zhenai/parser"
)

const CrawlerUrl = "http://www.zhenai.com/zhenghun"
const profileDatabase = "dating_profile"

// gopm get -g -v golang.org/x/text
func main() {
	itemChan, err := persist.GetItemSaver(profileDatabase)
	if err != nil {
		panic(err)
	}

	engine.ConcurrentEngine{
		Scheduler: &scheduler.QueuedScheduler{},
		ItemSaver: itemChan,
		WorkCount: 15,
	}.Run(engine.Request{
		Url:        CrawlerUrl,
		ParserFunc: parser.GetCityList,
	})

}
