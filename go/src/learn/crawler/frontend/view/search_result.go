package view

import (
	"html/template"
	"io"
	"learn/crawler/frontend/model"
)

type SearchResultView struct {
	template *template.Template
}

func CreateSearchResultView(fileName string) SearchResultView {
	return SearchResultView{
		template: template.Must(template.ParseFiles(fileName)),
	}
}

func (view SearchResultView) Render(io io.Writer, data model.SearchResult) error {
	return view.template.Execute(io, data)
}
