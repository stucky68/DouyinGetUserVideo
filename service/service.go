package service

import (
	"DouyinDownload/model"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var httpData string

var UA string

func ParserConfig(data string) {
	var config model.DownloadConfig
	err := json.Unmarshal([]byte(data), &config)
	if err != nil {
		panic(err)
	}

	UA = config.UA

}



func HandleJson(data model.Data, count *int) bool {
	flag := false
	for _, item := range data.AwemeList {
		index := strings.Index(item.Video.Origin_cover.Uri, "_")
		videoTime, _ := strconv.ParseInt(item.Video.Origin_cover.Uri[index+1:], 10, 64)
		//fmt.Println(videoTime)

		if videoTime >= 1577808000 {
			*count++
			flag = true
		}
	}
	return flag
}

func GetData(url string) (tac, dytk string, err error) {
	client := &http.Client{}

	rand.Seed(time.Now().Unix())
	s := strconv.Itoa(rand.Intn(100000))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/"+s+".1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/"+s+".1")

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "07056837-9ac4-4d7a-bc47-3af9ffb58e40")
	req.Header.Add("Cookie", "odin_tt=cbdfbdf9bc6a050b5eb1847318a9632061e9f327b8a6ef1eda145d6c25a8bf4e863bffc7d109be0efe8188bf1a59755b1b804ef6bb62bc6d9eefdb57a8640553")
	req.Header.Add("Referer", url)
	req.Header.Add("Connection", "keep-alive")
	res, err := client.Do(req)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	result := string(b)

	var tacRegexp = regexp.MustCompile(`<script>(.*?)</script>`)
	tacs := tacRegexp.FindStringSubmatch(result)
	if len(tacs) > 1 {
		tac = tacs[1]
	} else {
		return tac, dytk, errors.New("查找tac失败")
	}

	var dytkRegexp = regexp.MustCompile(`dytk: '(.*?)'`)
	dytks := dytkRegexp.FindStringSubmatch(result)
	if len(dytks) > 1 {
		dytk = dytks[1]
	} else {
		return tac, dytk, errors.New("查找dytk失败")
	}
	return
}

func GetVideo(user_id, signature, dytk string, max_cursor int64) (error, model.Data) {
	client := &http.Client{}
	url := "https://www.iesdouyin.com/web/api/v2/aweme/post/?user_id=" + user_id + "&sec_uid=&count=21&max_cursor=" + strconv.FormatInt(max_cursor, 10) + "&aid=1128&_signature=" + signature + "&dytk=" + dytk
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, model.Data{}
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("cookie", "_ga=GA1.2.938284732.1578806304; _gid=GA1.2.1428838910.1578806305")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", UA)
	req.Header.Add("Host", "www.iesdouyin.com")
	res, err := client.Do(req)
	if err != nil {
		return err, model.Data{}
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err, model.Data{}
	}
	var data model.Data
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err, model.Data{}
	}
	return nil, data
}
