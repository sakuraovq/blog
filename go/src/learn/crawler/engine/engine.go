package engine

import (
	"learn/crawler/fetcher"
	"log"
)

type SimpleEngine struct{}

func (SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	for _, req := range seeds {
		requests = append(requests, req)
	}

	for len(requests) > 0 {
		req := requests[0]
		requests = requests[1:]

		result, e := worker(req)
		if e != nil {
			continue
		}
		for _, items := range result.Items {
			// %v 打印原始数据
			log.Printf("Got item %v", items)
		}

		// 填充队列
		requests = append(requests, result.Request...)
	}

}

func worker(req Request) (ParserResult, error) {
	log.Printf("Fetching url %s", req.Url)
	body, e := fetcher.Fetch(req.Url)
	if e != nil {
		log.Printf("fetch error %v", e)
		return ParserResult{}, e
	}
	parserResult := req.ParserFunc(body)
	return parserResult, nil
}
