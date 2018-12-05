package parser

import (
	"learn/crawler/engine"
	"regexp"
)

const getUserRule = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`

// 获取城市下面的用户
func GetCity(contents []byte) engine.ParserResult {

	regx := regexp.MustCompile(getUserRule)
	subMatch := regx.FindAllSubmatch(contents, -1)
	parserResult := engine.ParserResult{}
	for _, match := range subMatch {
		parserResult.Request = append(parserResult.Request,
			engine.Request{
				Url:        string(match[1]),
				ParserFunc: UserProfile,
			},
		)
		parserResult.Items = append(parserResult.Items, "User "+string(match[2]))
	}
	return parserResult
}
