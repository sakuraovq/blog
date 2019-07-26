package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func login() {
	////获取登陆界面的cookie
	//c := &http.Client{
	//	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	//		fmt.Println(req)
	//		fmt.Println(via)
	//		return errors.New("redirecting")
	//	},
	//}
	//
	//str := `{"login":"chjlu666","part_key":"","password":"hjlz2018321","error":"","success":"","isFetching":false,"loginType":"account","verifyRequestCode":"","verifyResponseCode":"","captchaCode":"","verifyType":null,"rohrToken":"eJzFlG9v2jAQxr+LJXgzi9hObBIkNEE7qqxraTvYWlUVCsGEjGBHiQP9o373nUPJ2mndXk6Kkp/Pj++eOOc8oSJcoN4TWmZRgnqUEOK5zxhtZYF6iHZIRyCMTAlTnHMvYH5Amc8xit/EuKAeRvPi2zHq3fqEYEHInQ1cwfiWCp9gSjiE9hx47A4zDy4rCkGDVsbkZc9xZB6VZa4L09nI1FSR6sR640RxrCtlnEqlJt3KTCep+jhPZqWuilj23XYpi20KtIvSTZS28ywyS11s+qwda2VSVcm+rWALdPaaN/m37FdhRypTPLQzuTSzuTZGb2ZZqtb9Fhu92AB6MVKmiaryljtqvLTcY7fFxIsfGO2rQejgCWIMhgdfMLTOWu4A0sL1J38Q3jI7d/AIbEtPc0gk783RKlJKZk01uQ/PrGZW5bP49/l2kSarf79eIWMN9//yftZzIi92C+i/v3066AWtEPTRZmL7yOUC0y6FRS73gdi71G10lkhN4jUJvyYPSNTkArk1UUx5UBM0NhcNeZa8AIg1RBsib8gjkIX4NUGMWJ1Hu3A4LDA4JXUEjolfixjHfvcdAI14X7NfLnBA7D6t7T7BM3q9X3g4nUzG500yHJ5fTCdWbg7yM/ghwLztKSD5eTfaxZ9Ok6PB2eB86Pjy+OoyXg6m9/7XkMbjxXC8WF07S/nQFTpxvujvYblmIYsey8Hq8gMP9eJ0PjzLrkfFyr0ZRD/MCTnhw5trKrJiOFZ8kmQkpkcPSxLJq2CNnn8CMAlbvg=="}`
	//
	//var loginUrl = "https://epassport.meituan.com/api/account/login?service=waimai&bg_source=3&loginContinue=http:%2F%2Fe.waimai.meituan.com%2Fv2%2Fepassport%2Fentry&loginType=account"
	//
	//req, err := http.NewRequest("POST", loginUrl, strings.NewReader(str))
	//
	//if err != nil {
	//	panic(err)
	//}
	//req.Header.Set("Cookie","uuid=9dc9dcffcd6dfd512859.1556095273.1.0.0;")
	//req.Header.Set("Accept", "application/json")
	//req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	//req.Header.Set("Referer", "https://epassport.meituan.com/account/unitivelogin?bg_source=3&service=waimai&platform=2&continue=http://e.waimai.meituan.com/v2/epassport/entry&left_bottom_link=%2Faccount%2Funitivesignup%3Fbg_source%3D3%26service%3Dwaimai%26platform%3D2%26continue%3Dhttp%3A%2F%2Fe.waimai.meituan.com%2Fv2%2Fepassport%2FsignUp%26extChannel%3Dwaimaie%26ext_sign_up_channel%3Dwaimaie&right_bottom_link=%2Faccount%2Funitiverecover%3Fbg_source%3D3%26service%3Dwaimai%26platform%3D2%26continue%3Dhttp%3A%2F%2Fe.waimai.meituan.com%2Fv2%2Fepassport%2FchangePwd")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	//req.Header.Set("x-requested-with", "XMLHttpRequest")
	//
	//res, err := c.Do(req)
	//
	//if err != nil {
	//	panic(err)
	//}
	//data, _ := ioutil.ReadAll(res.Body)
	//res.Body.Close()
	////{"status":{"code":0,"message":"成功"},"needChangePassword":null,"captchaToken":null,"maskMobile":null,"verifyTO":null,"needChangeTO":{"needChangePassword":false,"needChangeLogin":false,"needChangeName":false,"needChangeContact":false},"bizAccountInfo":{"login":null,"name":null,"contact":null},"bsid":"9rR4O3NTGPH8FUmCqefgq6dbzeo0WgqrytHqCTNh8QqMSeMph5aM_hngBYtkt_sDutUuRa0_yjjjDHeCSMT9ng","continue":"http://e.waimai.meituan.com/v2/epassport/entry"}
	//var response = make(map[string]interface{})
	//json.Unmarshal(data, &response)
	//
	////continueLoginUrl := fmt.Sprintf("http://e.waimai.meituan.com/v2/epassport/entry?BSID=%s&source=TOKEN_SOURCE_LOGIN", response["bsid"])
	////
	////request1, err := http.NewRequest(http.MethodGet, continueLoginUrl, nil)
	////if err != nil {
	////	panic(err)
	////}
	////
	////for _, v := range respCookies {
	////	request1.AddCookie(v)
	////}
	//
	//// 获取响应
	////res1 := client(request1)
	////res1body, err := ioutil.ReadAll(res1.Body)
	//fmt.Println(response)

	response := make(map[string]interface{})
	response["bsid"] = "Yzx7xkYc_LmwCvXaHLBDx94JYvyQ8TSqXPPVdMYZe7uNR2OGB32O3gB-dbUTWkMyhbtvbW5AImKiiKAu7A8_kQ"

	BSID := response["bsid"].(string)

	cli := &http.Client{}

	logonUrl := "http://e.waimai.meituan.com/v2/epassport/logon"
	request2, _ := http.NewRequest(http.MethodPost, logonUrl, strings.NewReader("BSID="+BSID))

	request2.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request2.Header.Set("Referer", "http://e.waimai.meituan.com/v2/epassport/entry?BSID="+BSID)

	request2.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	request2.Header.Set("Cookie", "device_uuid=!f71a3646-3e4d-4a1a-bedb-85f719ab8c89; uuid_update=true; uuid=f409d461be652f0d7dc3.1555479525.1.0.0; pushToken=0QWRBxRofwVF1CG5wv_u1xJpmFKot6G6N-KALYApEhwU*; _lxsdk_cuid=16a2a71e47fc8-077ce0772aa44b-36697e04-1aeaa0-16a2a71e47fc8; _lxsdk=16a2a71e47fc8-077ce0772aa44b-36697e04-1aeaa0-16a2a71e47fc8; _lxsdk_s=16a2a71e480-ead-1d-03e%7C%7C11; wpush_server_url=wss://wpush.meituan.com; shopCategory=food; JSESSIONID=suuqorkydlc29amuateb3frh")
	response3, _ := cli.Do(request2)

	resbody3, _ := ioutil.ReadAll(response3.Body)

	fmt.Println("token", string(resbody3))

}

func main() {
	login()
}

func client(req *http.Request) *http.Response {
	cli := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
			req.Header.Set("Host", "e.waimai.meituan.com")
			req.Header.Set("Upgrade-Insecure-Requests", "1")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")

			fmt.Println("CheckRedirect", req.URL)
			return nil
		},
	}
	response, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	return response
}
