package repositories

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func Test_Get_Long_Urls(t *testing.T) {

	// setup
	type LongUrlsStruct struct {
		LongUrls map[string]string `json:"long_urls"`
	}

	longUrls := LongUrlsStruct{
		LongUrls: map[string]string{
			"https://facebook.com": "",
		},
	}

	data, _ := json.Marshal(longUrls)

	body := ioutil.NopCloser(bytes.NewReader([]byte(data)))

	// act
	urlsRepository := UrlsRepository{}

	LongUrls := urlsRepository.GetLongUrls(body)

	var LongUrlsType interface{} = LongUrls

	// assert
	_, isMapString := LongUrlsType.(map[string]string)

	if !isMapString {
		t.Errorf("Invalid data type. Data type must be map[string]string")
	}

	if _, key := LongUrls["https://facebook.com"]; !key {
		t.Errorf("Long urls [ https://facebook.com ] did not decoded")
	}
}

func Test_Get_Long_Urls_Throws_Error_If_Provided_Data_Will_Invalid(t *testing.T) {

	// assert
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Get long urls did not panic")
		}
	}()

	// setup
	data := "invalid data..."

	body := ioutil.NopCloser(bytes.NewReader([]byte(data)))

	// act
	urlsRepository := UrlsRepository{}

	urlsRepository.GetLongUrls(body)
}
