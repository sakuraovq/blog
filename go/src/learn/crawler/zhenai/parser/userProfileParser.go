package parser

import (
	"learn/crawler/engine"
	"learn/crawler_distributed/config"
)

type UserProfileParser struct {
	UserGender string
}

// 实例化一个指针 user parser
func NewUserProfileParser(gender string) *UserProfileParser {
	return &UserProfileParser{
		UserGender: gender,
	}
}

func (userParser *UserProfileParser) ParserFunc(contents []byte, url string) engine.ParserResult {
	return UserProfile(contents, url, userParser.UserGender)
}

func (userParser *UserProfileParser) Serialize() (name string, args interface{}) {
	return config.ParserProfile, userParser.UserGender
}
