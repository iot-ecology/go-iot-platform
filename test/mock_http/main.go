package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	for {
		send()

		time.Sleep(30*time.Second)

	}

}
func send() {
	url := "http://192.168.3.101:8888/handler"
	method := "POST"

	// 手动设置用户名和密码
	username := "admin"
	password := "admin"
	authStr := username + ":" + password
	auth := base64.StdEncoding.EncodeToString([]byte(authStr))

	payload := strings.NewReader(`{
    "data":"2000"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("device_id", "aaa")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+auth)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}