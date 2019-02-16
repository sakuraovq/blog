package controller

import (
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"learn/crawler/engine"
	"learn/crawler/frontend/model"
	"learn/crawler/frontend/view"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type SearchResultHandler struct {
	view   view.SearchResultView
	client *elastic.Client
}

func (search SearchResultHandler) ServeHTTP(
	writer http.ResponseWriter, req *http.Request) {
	q := strings.TrimSpace(req.FormValue("q"))
	from, err := strconv.Atoi(req.FormValue("from"))
	if err != nil {
		from = 0
	}
	result, err := search.getSearchResult(q, from)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	err = search.view.Render(writer, result)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func CreateSearchResultHandler(temp string) SearchResultHandler {
	esClient, err := elastic.NewClient(
		elastic.SetSniff(false))

	if err != nil {
		panic(err)
	}
	return SearchResultHandler{
		view:   view.CreateSearchResultView(temp),
		client: esClient,
	}
}

func (search SearchResultHandler) getSearchResult(
	q string, from int) (model.SearchResult, error) {
	var searchResult model.SearchResult
	// 获取替换的字段名
	regx := regexp.MustCompile(`([A-Za-z]*):`)
	searchString := regx.ReplaceAllString(q, "Payload.$1:")

	resp, err := search.client.Search("dating_profile").
		Query(elastic.NewQueryStringQuery(searchString)).
		From(from).
		Do()
	if err != nil {
		return searchResult, err
	}
	searchResult.Hits = int(resp.TotalHits())
	searchResult.Start = from
	searchResult.Query = q
    // 上一页
	searchResult.Prev = from - 10
    // 下一页
	searchResult.Next = from + 10

	for _, eachItem := range resp.Each(reflect.TypeOf(engine.Item{})) {
		item, isOk := eachItem.(engine.Item)
		if !isOk {
			log.Println("error type", item)
			return searchResult, fmt.Errorf("undefinded item %v", item)
		}
		searchResult.Items = append(searchResult.Items, item)
	}
	return searchResult, nil
}
