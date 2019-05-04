package main

import (
	"fmt"
	"learn/interface/faker"
)

// 接口定义又用户决定 实现方法就行
type Retriever interface {
	Get(url string) string
}

type Poster interface {
	Post(post map[string]string) string
}

// 组合接口
type RetrieverPoster interface {
	Retriever
	Poster
}

type Post struct {
	content map[string]string
}

func (post *Post) session(p RetrieverPoster) string {
	p.Post(post.content)
	return p.Get("http://sakuraus.cn")
}

func main() {
	var p RetrieverPoster
	// 指针实现的接口只能用指针
	p = &faker.Post{}
	post := Post{map[string]string{
		"contents": "test",
	}}

	// 解析接口指针指向的值方式:1
	switch p.(type) {
	case *faker.Post:
		fmt.Println("switch case faker")
	}

	// 方式:2
	if faker, ok := p.(*faker.Post); ok {
		fmt.Println("if ok faker", faker)
	}
	// 接口中有两个值 一个指针一个值
	fmt.Printf("%T %v", p, p)
	fmt.Println()
	post.session(p)
}
