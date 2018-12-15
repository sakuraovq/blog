package parser

import (
	"fmt"
	"learn/crawler/engine"
	"regexp"
)

var (
	hobbyRegx, profileRegx *regexp.Regexp
)

func init() {
	// 兴趣爱好匹配 示例这两种
	hobbyRegx = regexp.MustCompile(`<div class="m-btn pink" [^>]*>([^>]+)</div>`)
	// 用户基本信息匹配
	profileRegx = regexp.MustCompile(`<div class="m-btn purple" [^>]*>([^>]+)</div>`)
}
func UserProfile(contents []byte) engine.ParserResult {

	parserResult := engine.ParserResult{}

	hobbyMatch := hobbyRegx.FindAllSubmatch(contents, -1)
	for _, sub := range hobbyMatch {
		fmt.Println(string(sub[1]))
	}

	profileMatch := profileRegx.FindAllSubmatch(contents, -1)
	for _, sub := range profileMatch {
		fmt.Println(string(sub[1]))
	}

	parserResult.Request = append(parserResult.Request,
		engine.Request{
			Url:        "",
			ParserFunc: engine.NilParserResult,
		},
	)
	parserResult.Items = append(parserResult.Items, "UserProfile "+"")

	return parserResult
}
