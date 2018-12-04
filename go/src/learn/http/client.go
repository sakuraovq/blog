package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {

	//resp, err := http.Get("http://sakuraus.cn")
	//if err != nil {
	//	panic(err)
	//}
	//

	request, _ := http.NewRequest(http.MethodGet, "http://google.cn", nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Mobile Safari/537.36")

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("Redircet:", req)
			return nil
		},
	}
	client.Timeout = time.Second * 2
	//resp, err := http.DefaultClient.Do(request)
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer resp.Body.Close()
	bytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", bytes)
}
