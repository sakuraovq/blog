package parser

import (
	"learn/crawler/engine"
	"learn/crawler_distributed/config"
	"regexp"
)

const cityListParserRule = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^>]+)</a>`

// 获取城市列表
func GetCityList(contents []byte, url string) engine.ParserResult {

	regx := regexp.MustCompile(cityListParserRule)
	subMatch := regx.FindAllSubmatch(contents, -1)
	parserResult := engine.ParserResult{}
	for _, m := range subMatch {
		// 生成Request
		parserResult.Request = append(parserResult.Request,
			engine.Request{
				Url:   string(m[1]),
				Parse: engine.NewParserFunc(config.ParserCity, GetCity),
			})
	}

	return parserResult
}
