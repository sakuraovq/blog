package engine

type ParserFunction func([]byte, string) ParserResult

type Parser interface {
	ParserFunc(contents []byte, url string) ParserResult
	// 序列化
	Serialize() (name string, args interface{})
}

type Request struct {
	Url   string
	Parse Parser
}

type ParserResult struct {
	Request []Request
	Items   []Item
}

type Item struct {
	Url     string
	Id      string
	Type    string
	Payload interface{}
}

type NilParser struct{}

func (parser *NilParser) ParserFunc([]byte, string) ParserResult {
	return ParserResult{}
}

func (parser *NilParser) Serialize() (name string, args interface{}) {
	return "NilParser", nil
}



type Parse struct {
	Parser   ParserFunction
	FuncName string
}

func NewParserFunc(name string, Parser ParserFunction) *Parse {
	return &Parse{
		FuncName: name,
		Parser:   Parser,
	}
}

func (p *Parse) ParserFunc(contents []byte, url string) ParserResult {
	return p.Parser(contents, url)
}

func (p *Parse) Serialize() (name string, args interface{}) {
	return p.FuncName, nil
}
