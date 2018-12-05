package parser

import (
	"learn/crawler/engine"
	"regexp"
)

const cityListParserRule = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^>]+)</a>`

// 获取城市列表
func GetCityList(contents []byte) engine.ParserResult {

	regx := regexp.MustCompile(cityListParserRule)
	subMatch := regx.FindAllSubmatch(contents, -1)
	parserResult := engine.ParserResult{}
	for _, m := range subMatch {
		// 生成Request
		parserResult.Request = append(parserResult.Request,
			engine.Request{
				Url:        string(m[1]),
				ParserFunc: GetCity,
			})
		parserResult.Items = append(parserResult.Items, "City "+string(m[2]))
	}
	return parserResult
}
