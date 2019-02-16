package fetcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// 请求限速 10毫秒一个
//var rateLimit = time.Tick(10 * time.Millisecond)

func Fetch(url string) ([]byte, error) {
	log.Print("fetching url ", url)
	//<-rateLimit
	request, e := http.NewRequest(http.MethodGet, url, nil)
	if e != nil {
		return nil, fmt.Errorf("make Request Fail url is %s", url)
	}
	// 没有user-agent 不能访问
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("send Requset Though"+
			" Fetch Error Code %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
