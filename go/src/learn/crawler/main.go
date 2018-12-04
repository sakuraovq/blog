package main

import (
	"learn/crawler/engine"
	"learn/crawler/zhenai/parser"
)

const CrawlerUrl = "http://www.zhenai.com/zhenghun"

// gopm get -g -v golang.org/x/text
func main() {

	engine.Run(engine.Request{
		Url:        CrawlerUrl,
		ParserFunc: parser.GetCityList,
	})

}
