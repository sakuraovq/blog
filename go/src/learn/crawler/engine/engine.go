package engine

import (
	"fmt"
	"learn/crawler/fetcher"
	"log"
)

func Run(seeds ...Request) {
	// request queue
	var requests []Request
	for _, req := range seeds {
		requests = append(requests, req)
	}
	for len(requests) > 0 {
		req := requests[0]
		requests = requests[1:]
		log.Printf("fetching url %s", req.Url)

		body, e := fetcher.Fetch(req.Url)
		if e != nil {
			fmt.Println("fetch error", e.Error())
			continue
		}
		parserResult := req.ParserFunc(body)
		// 填充队列
		requests = append(requests, parserResult.Request...)

		for _, items := range parserResult.Items {
			// %v 打印原始数据
			log.Printf("Got item %v", items)
		}

	}
}
