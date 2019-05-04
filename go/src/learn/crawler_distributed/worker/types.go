package worker

import (
	"fmt"
	"learn/crawler/engine"
	"learn/crawler/zhenai/parser"
	"learn/crawler_distributed/config"
	"log"
)

type SerializedParser struct {
	FuncName string
	Args     interface{}
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

// 序列化req
func SerializedRequest(req engine.Request) Request {
	name, args := req.Parse.Serialize()
	return Request{
		Url: req.Url,
		Parser: SerializedParser{
			FuncName: name,
			Args:     args,
		},
	}
}

// 反序列化req
func UnSerializedRequest(serializedRequest Request) (engine.Request, error) {
	parser, err := GetSerializedRequestParser(serializedRequest.Parser)

	if err != nil {
		return engine.Request{}, err
	}
	req := engine.Request{
		Url:   serializedRequest.Url,
		Parse: parser,
	}

	return req, nil
}

// 获取能工作的parser
func GetSerializedRequestParser(serializedParser SerializedParser) (engine.Parser, error) {

	switch serializedParser.FuncName {
	case config.ParserCityList: // 获取城市列表
		return engine.NewParserFunc(config.ParserCityList, parser.GetCityList), nil
	case config.ParserCity: // 获取城市
		return engine.NewParserFunc(config.ParserCity, parser.GetCity), nil
	case config.ParserProfile: // 获取用户
		if gender, ok := serializedParser.Args.(string); ok {
			return parser.NewUserProfileParser(gender), nil
		}
		return nil, fmt.Errorf("invail parser UserProfileParser args %v", serializedParser.Args)
	case config.NilParser: // 空的parser
		return &engine.NilParser{}, nil
	default:
		return nil, fmt.Errorf("unkown parser %v", serializedParser)
	}
}

// 序列化result
func SerializedResult(parserResult engine.ParserResult) ParseResult {

	result := ParseResult{
		Items: parserResult.Items,
	}

	for _, req := range parserResult.Request {
		result.Requests = append(result.Requests, SerializedRequest(req))
	}
	return result
}

// 反序列化result
func UnSerializedResult(serializedResult ParseResult) engine.ParserResult {

	parserResult := engine.ParserResult{
		Items: serializedResult.Items,
	}

	for _, serializedReq := range serializedResult.Requests {
		// 一个解析失败了就算了 看情况
		request, err := UnSerializedRequest(serializedReq)
		if err != nil {
			log.Printf("UnSerializedRequest Fail %v", err)
			continue
		}

		parserResult.Request = append(parserResult.Request, request)
	}
	return parserResult
}
