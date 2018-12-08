package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"learn/crawler/engine"
	"learn/crawler/fetcher"
	"strings"
)

func init() {

}
func UserProfile(c []byte) engine.ParserResult {

	userResult := engine.ParserResult{}
	contents, e := fetcher.Fetch("http://album.zhenai.com/u/109816882")
	fmt.Println("content",string(contents),e)
	document, e := goquery.NewDocumentFromReader(strings.NewReader(string(contents)))
	if e != nil {
		fmt.Println("goquery error", e,c)
	}
	document.Find(".pink-btns").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
	return userResult
}
