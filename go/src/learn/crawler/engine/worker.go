package engine

import (
	"learn/crawler/fetcher"
	"log"
)

func Worker(req Request) (ParserResult, error) {
	body, e := fetcher.Fetch(req.Url)
	if e != nil {
		log.Printf("fetch error %v", e)
		return ParserResult{}, e
	}

	parserResult := req.Parse.ParserFunc(body, req.Url)
	return parserResult, nil
}
