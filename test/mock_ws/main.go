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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	input := auth()
	var data map[string]interface{}

	// 解析 JSON 字符串到 map 中
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// 提取 uid 的值
	uid, ok := data["uid"].(string)
	if ok {
		fmt.Println("uid:", uid)
	} else {
		fmt.Println("uid not found or not a string")
	}

	wsclient(uid)
}

func auth() string {
	url := "http://localhost:13332/auth"
	method := "GET"

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
		return ""
	}
	req.Header.Add("device_id", "7")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+auth)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	return string(body)
}

func wsclient(uid string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:13332/ws?id="+uid, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
