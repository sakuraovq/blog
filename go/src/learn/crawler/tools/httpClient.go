package tools

import (
	"fmt"
	"net/http"
	"time"
)

func GetUserProfile(url string, uid int) {

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	request.Header.Add("Cookie", "ipCityCode=10127002; sid=eCb1RlnRnnR9sIApVbCJ; clientp=42836")
	request.Header.Add("Host", "album.zhenai.com")
	request.Header.Add("Referer", "http://album.zhenai.com/u/1690779978")
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

}
