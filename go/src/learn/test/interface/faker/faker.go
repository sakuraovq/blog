package faker

import "fmt"

type Post struct {
	content string
}

func (p *Post) Post(ma map[string]string) string {
	p.content = ma["contents"]
	return ""
}

func (p *Post) Get(url string) string {
	fmt.Println(p.content, url)
	return "ok"
}
