/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

		time.Sleep(10*time.Second)

	}

}
func send() {
	url := "http://localhost:8888/handler"
	method := "POST"

	// 手动设置用户名和密码
	username := "afff"
	password := "fff"
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
	req.Header.Add("device_id", "6")
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