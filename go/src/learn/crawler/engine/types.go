package engine

type Request struct {
	Url        string
	ParserFunc func([]byte) ParserResult
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
