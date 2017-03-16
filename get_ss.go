// Created by ev3rs0u1 on 2017/3/15.
package main

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
)

const (
	APIKEY  = "20141110217674"
	APIURL  = "http://api.wwei.cn/dewwei.html"
	ImgURL  = "http://free.shadowsocks8.cc/images/server02.png"
	CfgPath = "./gui-config.json"
)

// JSON struct
type JSON struct {
	Data struct {
		RawData   string `json:"raw_data"`
		RawFormat string `json:"raw_format"`
		RawText   string `json:"raw_text"`
		RawType   string `json:"raw_type"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Status int64  `json:"status"`
}

// Profile struct
type Profile struct {
	Configs []struct {
		Method   string `json:"method"`
		Password string `json:"password"`
		Server   string `json:"server"`
		Port     int64  `json:"server_port"`
	} `json:"configs"`
}

func catch() {
	if err := recover(); err != nil {
		fmt.Println("[-]", err)
	}
}

// APIRequest ...
func APIRequest(url string, params map[string]string) string {
	timeout := time.Duration(3 * time.Second)
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", url, nil)
	query := req.URL.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	req.URL.RawQuery = query.Encode()
	res, err := client.Do(req)
	if err != nil {
		panic("URL connection failed")
	}
	defer res.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	return string(bodyBytes)
}

func decodeQRCode(str string) string {
	fmt.Println("[+]", "GET shadowsocks account...")
	JSONStruct := new(JSON)
	json.Unmarshal([]byte(str), &JSONStruct)
	if JSONStruct.Status != 1 {
		panic("GET shadowsocks account failed")
	}
	data := JSONStruct.Data.RawText[5:]
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic("Base64 decode failed")
	}
	return string(decoded)
}

func processProfile(str string) {
	fmt.Println("[+]", "Read config file...")
	fileBytes, err := ioutil.ReadFile(CfgPath)
	if err != nil {
		panic("Read config file failed")
	}
	JSONStruct := new(Profile)
	err = json.Unmarshal([]byte(string(fileBytes)), &JSONStruct)
	if err != nil {
		panic("JSON parse failed")
	}
	configs := JSONStruct.Configs
	ary0 := strings.Split(str, ":")
	method, aryStr, portStr := ary0[0], ary0[1], ary0[2]
	ary1 := strings.Split(aryStr, "@")
	pwd, server := ary1[0], ary1[1]
	port, _ := strconv.ParseInt(portStr[:len(portStr)-1], 10, 64)
	configs[0].Server = server
	configs[0].Method = method
	configs[0].Port = port
	configs[0].Password = pwd
	fmt.Printf(" [*] Server: %-15s\n", server)
	fmt.Printf(" [*] Method: %-15s\n", method)
	fmt.Printf(" [*] Port:   %-15d\n", port)
	fmt.Printf(" [*] Passwd: %-15s\n", pwd)
	body, _ := json.Marshal(JSONStruct)
	err = ioutil.WriteFile(CfgPath, body, 0644)
	if err != nil {
		panic("Write config file failed")
	}
	fmt.Println("[+]", "Success!")
}

func main() {
	defer catch()
	params := map[string]string{"data": ImgURL, "apikey": APIKEY}
	body := APIRequest(APIURL, params)
	processProfile(decodeQRCode(body))
}
