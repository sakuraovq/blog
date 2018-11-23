package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type appHandler func(writer http.ResponseWriter, request *http.Request) error

const prefix = "/list/"

type UserErr interface {
	error
	Message() string
}

type UserReportError string

func (e UserReportError) Error() string {
	return e.Message()
}

func (e UserReportError) Message() string {
	return string(e)
}

func errWrapper(handler appHandler) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		defer func() {
			// 保护自己保护一次
			if r := recover(); r != nil {
				http.Error(writer,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		err := handler(writer, request)
		if err != nil {
			log.Println(err.Error())
			// 返回用户能看见的错误信息
			if err, ok := err.(UserErr); ok {
				http.Error(writer, err.Message(), http.StatusBadRequest)
				return
			}

			code := http.StatusOK
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			http.Error(writer, http.StatusText(code), code)
		}
	}
}

func fileList(writer http.ResponseWriter, request *http.Request) error {
	url := request.URL.Path
	if !strings.HasPrefix(url, prefix) {
		return UserReportError("path must start with " + prefix)
	}
	path := request.URL.Path[len(prefix):]
	file, e := os.Open(path)
	if e != nil {
		return e
	}
	defer file.Close()

	bytes, e := ioutil.ReadAll(file)
	if e != nil {
		return e
	}
	writer.Write(bytes)
	return nil
}

func main() {
	http.HandleFunc("/", errWrapper(fileList))
	serve := http.ListenAndServe(":80", nil)
	if serve != nil {
		panic(serve)
	}
}
