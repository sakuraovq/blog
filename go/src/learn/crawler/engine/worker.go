package engine

import (
	"learn/crawler/fetcher"
	"log"
)

func worker(req Request) (ParserResult, error) {
	//log.Printf("Fetching url %s", req.Url)
	body, e := fetcher.Fetch(req.Url)
	if e != nil {
		log.Printf("fetch error %v", e)
		return ParserResult{}, e
	}
	parserResult := req.ParserFunc(body, req.Url)
	return parserResult, nil
}
