package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func errPanic(w http.ResponseWriter, request *http.Request) error {
	panic(123)
}


func erroNotFound(w http.ResponseWriter, request *http.Request) error {
	return  os.ErrNotExist
}

func errorUnknown(w http.ResponseWriter, request *http.Request) error {
	return  errors.New("unknown")
}
func errNoPermission(w http.ResponseWriter, request *http.Request) error {
	return  os.ErrPermission
}

func TestErrWrapper(t *testing.T) {
	tests := []struct {
		h       appHandler
		code    int
		message string
	}{
		{errPanic, 500, "server error"},
		{errNoPermission, 403, "not permission"},
		{erroNotFound, 404, "not found"},

	}
	for _, tt := range tests {
		f := errWrapper(tt.h)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "https://sakuraus.cn", nil)
		f(response, request)
		bytes, _ := ioutil.ReadAll(request.Body)
		body := string(bytes)
		body = strings.Trim(body, "\n")
		if response.Code != tt.code || body != tt.message {
			t.Errorf("expect(%d ,%s);got(%d,%s)",
				tt.code, tt.message, response.Code, body)
		}
	}
}
