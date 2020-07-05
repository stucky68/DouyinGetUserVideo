package main

import (
	"DouyinDownload/service"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type GetSignaturePara struct {
	Uid string `json:"uid"`
	Tac string `json:"tac"`
	UA  string `json:"ua"`
}

type GetSignatureResult struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func GetSignature(userID, ua string) (string, string) {
	tac := ""
	dytk := ""

	tac, dytk, err := service.GetData("https://www.iesdouyin.com/share/user/" + userID)
	for err != nil {
		//fmt.Println(err)
		tac, dytk, err = service.GetData("https://www.iesdouyin.com/share/user/" + userID)
	}

	client := &http.Client{}
	data, _ := json.Marshal(GetSignaturePara{
		Uid: userID,
		Tac: base64.StdEncoding.EncodeToString([]byte(tac)),
		UA:  ua,
	})
	req, err := http.NewRequest("GET", "http://127.0.0.1:3000", bytes.NewBuffer(data))
	if err != nil {
		return "", ""
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", ""
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", ""
	}
	result := &GetSignatureResult{}
	err = json.Unmarshal(b, result)
	if err != nil {
		return "", ""
	}
	return result.Result, dytk
}

func read3(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}


func main() {
	fmt.Println("当前版本2020-06-13")
	//time.Sleep(time.Second * 3)
	service.ParserConfig(read3("./config.json"))
	//user := read3("./user.txt")
	//uids := strings.Split(user, "\r\n")
	uids := []string{"62743508192"}
	for _, userID := range uids {
		os.MkdirAll("download/"+userID, os.ModePerm)
		signature, dytk := GetSignature(userID, service.UA)
		max_cursor := int64(0)
		count := 0
		log.Println(userID + "正在查询")
		for {
			err, d := service.GetVideo(userID, signature, dytk, max_cursor)
			if err != nil {
				continue
			}
			flag := service.HandleJson(d, &count)

			if d.HasMore {
				//签名失效 重新获取
				if d.MinCursor == 0 && d.MaxCursor == 0 || len(d.AwemeList) == 0{
					//signature, dytk = GetSignature(userID, service.UA)
					log.Println("签名失效", count)
					time.Sleep(time.Second * 1)
					continue
				} else {
					max_cursor = d.MaxCursor
				}
			} else {
				break
			}

			if !flag {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		log.Println(userID + "查询完毕 共" + strconv.FormatInt(int64(count), 10) + "条")
	}
}
