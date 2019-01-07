package parser

import (
	"fmt"
	"learn/crawler/engine"
	"learn/crawler/model"
	"regexp"
	"strconv"
	"strings"
)

var (
	hobbyRegx, profileRegx, guessLikeRegx, guessLikeUserNameRegx, IdRegx *regexp.Regexp
)

func init() {
	// 兴趣爱好匹配 示例这两种
	hobbyRegx = regexp.MustCompile(`<div class="m-btn pink" [^>]*>([^>]+)</div>`)
	// 用户基本信息匹配
	profileRegx = regexp.MustCompile(`<div class="m-btn purple" [^>]*>([^>]+)</div>`)
	guessLikeRegx = regexp.MustCompile(
		`https://photo.zastatic.com/images/photo/[^/]+/([\d]+)/`)
	guessLikeUserNameRegx = regexp.MustCompile(
		`<span class="nickname f-clamp1" [^>]*>([^<]+)</span>`)
	IdRegx = regexp.MustCompile(`<div class="id" [^>]*>([^<]+)</div>`)
}

func MatchAtoi(str, filter string) int {
	val := strings.Trim(str, filter)
	i, e := strconv.Atoi(val)
	if e != nil {
		fmt.Println(str+" atoi err "+filter, e.Error())
		return 0
	}
	return i
}

func UserProfile(contents []byte, name, gender, url string, ) engine.ParserResult {
	parserResult := engine.ParserResult{}
	profile := model.Profile{}
	hobbyMatch := hobbyRegx.FindAllSubmatch(contents, -1)

	profile.Name = name
	profile.Gender = gender

	for idx, sub := range hobbyMatch {
		subStr := string(sub[1])
		switch idx {
		case 1:
			profile.RegisteredResidence = strings.TrimPrefix(subStr, "籍贯:")
		case 5:
			profile.House = subStr
		case 6:
			profile.Car = subStr
		}
	}

	userId := string(IdRegx.FindAll(contents, -1)[0])
	userId = strings.TrimPrefix(userId, "ID：")
	profileMatch := profileRegx.FindAllSubmatch(contents, -1)
	for idx, sub := range profileMatch {
		subStr := string(sub[1])
		if strings.HasPrefix(subStr, "工作地") {
			continue
		}
		switch idx {
		case 0:
			profile.Marriage = subStr
		case 1:
			profile.Age = MatchAtoi(subStr, "岁")
		case 2:
			profile.Constellation = subStr
		case 3:
			profile.Height = subStr
		case 4:
			profile.Weight = subStr
		case 5:
			fallthrough
		case 6:
			if strings.HasPrefix(subStr, "月收入") {
				profile.Income = strings.TrimPrefix(subStr, "月收入:")
			}
		case 7:
			profile.Occupation = subStr
		case 8:
			profile.Education = subStr
		}
	}

	parserResult.Items = []engine.Item{
		{
			Url:     url,
			Type:    "zhenai",
			Id:      userId,
			Payload: profile,
		},
	}

	// 猜你喜欢用户分析
	guessLikeUserId := guessLikeRegx.FindAllSubmatch(contents, -1)
	// 猜你喜欢用户名称
	//guessLikeUserNameList := guessLikeUserNameRegx.FindAllSubmatch(contents, -1)

	idx := 0
	for _, guessId := range guessLikeUserId {
		thisGuessId := string(guessId[1])

		if thisGuessId == userId || thisGuessId == "" {
			continue
		}
		//guessLikeUserName := string(guessLikeUserNameList[idx][1]) // 猜你喜欢用户名
		guessLikeUserName := ""
		guessUserUrl := "http://album.zhenai.com/u/" + thisGuessId

		guessRequest := engine.Request{
			Url: guessUserUrl,
			ParserFunc: func(guessContents []byte) engine.ParserResult {
				return UserProfile(guessContents, guessLikeUserName, profile.Gender, guessUserUrl)
			},
		}
		idx++
		parserResult.Request = append(parserResult.Request, guessRequest)
	}

	return parserResult
}
