package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	url := "http://msg.t.com/test/ti"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	for i := 0; i < 10000; i++ {
		go get(url, ctx)

	}

	time.Sleep(100 * time.Second)
	cancel()
}

func get(url string, ctx context.Context) {

	request, e := http.NewRequest(http.MethodGet, url, nil)
	if e != nil {
		fmt.Printf("make Request Fail url is %s", url)
		return
	}

	request.Header.Add("bid", "10006")

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Println("req error ", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("send Requset Though"+
			" Fetch Error Code %d", resp.StatusCode)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("erroor", err)
		return
	}

	fmt.Printf("res = %s", data)

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}
