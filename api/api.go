package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jdApi/conf"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func Api(express string) (res conf.TraceResStruct, err error) {
	jsonStr := formatOrderInfo(express)
	params := genDigest(jsonStr)
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, _ := client.Get(conf.URI + params)
	if resp != nil {
		defer resp.Body.Close()
		pageBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(pageBytes, &res)
	}

	return

}

func genDigest(dataStr string) string {

	params := url.Values{
		"360buy_param_json": []string{dataStr},
		"access_token":      []string{conf.ApiConf["access_token"]},
		"app_key":           []string{conf.ApiConf["app_key"]},
		"method":            []string{conf.ApiConf["method"]},
		"timestamp":         []string{timeString()},
		"v":                 []string{conf.ApiConf["v"]},
	}
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	str := conf.ApiConf["app_secret"]
	for _, k := range keys {
		str += k
		str += params[k][0]
	}
	str += conf.ApiConf["app_secret"]

	data := []byte(str)
	hash := md5.Sum(data)
	md5str := fmt.Sprintf("%x", hash)
	md5str = strings.ToUpper(md5str)
	params["sign"] = []string{md5str}

	return params.Encode()
}

func timeString() string {
	now := time.Now()

	loc, _ := time.LoadLocation("Asia/Shanghai")

	return now.In(loc).Format("2006-01-02 15:04:05")
}

func formatOrderInfo(express string) (jsonStr string) {
	t := conf.TraceStruct{express}
	jsonByte, _ := json.Marshal(t)
	jsonStr = string(jsonByte)
	return
}