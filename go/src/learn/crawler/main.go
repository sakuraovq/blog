package main

import (
	"learn/crawler/engine"
	"learn/crawler/scheduler"
	"learn/crawler/zhenai/parser"
	"learn/crawler_distributed/client"
	"learn/crawler_distributed/config"
)

// gopm get -g -v golang.org/x/text
func main() {
	itemChan, err := client.GetItemSaver()
	if err != nil {
		panic(err)
	}

	engine.ConcurrentEngine{
		Scheduler: &scheduler.QueuedScheduler{},
		ItemSaver: itemChan,
		WorkCount: 15,
	}.Run(engine.Request{
		Url:   config.ZhenAiSeedUrl,
		Parse: engine.NewParserFunc("GetCityList", parser.GetCityList),
	})

}
