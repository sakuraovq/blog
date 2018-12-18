package parser

import (
	"learn/crawler/engine"
	"regexp"
)

const getUserRule = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`
const getUserGenderRule = `<span class="grayL">性别：</span>([^>])`

// 获取城市下面的用户
func GetCity(contents []byte) engine.ParserResult {

	regx := regexp.MustCompile(getUserRule)
	genderRegx := regexp.MustCompile(getUserGenderRule)
	// 性别
	genders := genderRegx.FindAllSubmatch(contents, -1)
	// 用户地址 姓名
	subMatch := regx.FindAllSubmatch(contents, -1)
	parserResult := engine.ParserResult{}
	for k, match := range subMatch {
		// match 作用域在for里面 由于函数不会被当前循环调用 所有需要吧name 拷贝
		name := string(match[2])		// 昵称
		gender := string(genders[k][1]) // 性别
		parserResult.Request = append(parserResult.Request,
			engine.Request{
				Url: string(match[1]),
				ParserFunc: func(contents []byte) engine.ParserResult {
					return UserProfile(contents, name, gender)
				},
			},
		)
		parserResult.Items = append(parserResult.Items, "User "+name)
	}
	return parserResult
}
