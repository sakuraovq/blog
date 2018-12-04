package engine

type Request struct {
	Url        string
	ParserFunc func([]byte) ParserResult
}

type ParserResult struct {
	Request []Request
	Items   []interface{}
}

func NilParserResult(contents []byte) ParserResult {
	return ParserResult{}
}
