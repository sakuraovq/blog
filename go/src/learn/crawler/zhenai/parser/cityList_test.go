package parser

import (
	"learn/crawler/fetcher"
	"testing"
)

func TestGetCityList(t *testing.T) {

	bytes, e := fetcher.Fetch("http://www.zhenai.com/zhenghun")
	if e != nil {
		panic(e)
	}
	parserResult := GetCityList(bytes, "http://www.zhenai.com/zhenghun/aba")
	expectedUrls := []string{
		"http://www.zhenai.com/zhenghun/aba",
		"http://www.zhenai.com/zhenghun/akesu",
	}
	//expectedCites := []string{
	//	"City 阿坝", "City 阿克苏",
	//}
	const resultSize = 470

	if len(parserResult.Request) != resultSize {
		t.Errorf("result should have %d requests"+
			";but had %d", len(parserResult.Request), resultSize)
	}

	//if len(parserResult.Items) != resultSize {
	//	t.Errorf("result should have %d items"+
	//		";but had %d", len(parserResult.Items), resultSize)
	//}

	for i, url := range expectedUrls {
		if parserResult.Request[i].Url != url {
			t.Errorf(" expected url#%d:%s ;but had %s", i, url, parserResult.Request[i].Url)
		}
	}

	//for i, city := range expectedCites {
	//	if parserResult.Items[i].(string) != city{
	//		t.Errorf(" expected city#%d:%s ;but city %s", i, city, parserResult.Items[i].(string))
	//	}
	//}
}
