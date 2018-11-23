package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func errPanic(writer http.ResponseWriter, request *http.Request) error {
	panic(123)
}

func TestErrWrapper(t *testing.T) {
	tests := []struct {
		h       appHandler
		code    int
		message string
	}{
		{errPanic, 500, "server error"},
	}
	for _, tt := range tests {
		f := errWrapper(tt.h)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "https://sakuraus.cn", nil)
		f(response, request)
		bytes, _ := ioutil.ReadAll(request.Body)
		body := string(bytes)
		if response.Code != tt.code || body != tt.message {
			t.Errorf("expect(%d ,%s);got(%d,%s)",
				tt.code, tt.message, response.Code, body)
		}
	}
}
