package model

import "learn/crawler/engine"

// 搜索结果
type SearchResult struct {
	Hits  int
	Start int
	Prev  int
	Query string
	Next  int
	Items []engine.Item
}
