package main

import (
	"learn/crawler/engine"
	"learn/crawler/persist"
	"learn/crawler/scheduler"
	"learn/crawler/zhenai/parser"
)

const CrawlerUrl = "http://www.zhenai.com/zhenghun"

// gopm get -g -v golang.org/x/text
func main() {

	engine.ConcurrentEngine{
		Scheduler: &scheduler.QueuedScheduler{},
		ItemSaver: persist.GetItemSaver(),
		WorkCount: 15,
	}.Run(engine.Request{
		Url:        CrawlerUrl,
		ParserFunc: parser.GetCityList,
	})

}
