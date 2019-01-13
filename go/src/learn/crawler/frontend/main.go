package main

import (
	"learn/crawler/frontend/controller"
	"net/http"
)

func main() {

	http.Handle("/search", controller.CreateSearchResultHandler(
		"./view/search.html"))

	http.Handle("/", http.FileServer(http.Dir("./view/")))

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
