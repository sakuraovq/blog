package engine

type ParserFunction func([]byte, string) ParserResult

type Request struct {
	Url        string
	ParserFunc ParserFunction
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
