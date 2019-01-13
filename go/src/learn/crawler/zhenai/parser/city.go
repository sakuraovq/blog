package parser

import (
	"learn/crawler/engine"
	"regexp"
)

const getUserRule = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`
const getUserGenderRule = `<span class="grayL">性别：</span>([^>])`
const getNextPage = `href="(http://www.zhenai.com/zhenghun/[^"]+)"`

// 获取城市下面的用户
func GetCity(contents []byte, url string) engine.ParserResult {

	regx := regexp.MustCompile(getUserRule)
	genderRegx := regexp.MustCompile(getUserGenderRule)
	nextPage := regexp.MustCompile(getNextPage)

	// 性别
	genders := genderRegx.FindAllSubmatch(contents, -1)
	// 用户地址 姓名
	subMatch := regx.FindAllSubmatch(contents, -1)
	// 下一页城市/扩展链接
	nextPageCity := nextPage.FindAllSubmatch(contents, -1)

	parserResult := engine.ParserResult{}

	for _, city := range nextPageCity {
		parserResult.Request = append(parserResult.Request,
			engine.Request{
				Url:        string(city[1]),
				ParserFunc: GetCity,
			},
		)
	}

	for k, match := range subMatch {
		// match 作用域在for里面 由于函数不会被当前循环调用 所有需要吧name 拷贝
		parserResult.Request = append(parserResult.Request,
			engine.Request{
				Url:        string(match[1]),
				ParserFunc: userProfileFunc(string(genders[k][1])),
			},
		)
	}
	return parserResult
}

// 用户parserFunc
func userProfileFunc(gender string) engine.ParserFunction {
	return func(contents []byte, url string) engine.ParserResult {
		return UserProfile(contents, url, gender)
	}
}
