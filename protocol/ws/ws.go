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
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
)

var GWs map[string]*websocket.Conn
var GWsMirror map[*websocket.Conn]string

func InitWebSocket(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			log.Println("升级协议", r.Header["User-Agent"])
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		og := GWsMirror[conn]
		globalRedisClient.LRem(context.Background(), "ws_uid:"+globalConfig.NodeInfo.Name, 1, og)
		delete(GWs, og)
		delete(GWsMirror, conn)
		if err != nil {
			zap.S().Error(err)
		}
	}(conn)
	userId := c.Query("id")
	log.Println("用户id:", userId)

	val := globalRedisClient.LRange(context.Background(), "ws_uid:"+globalConfig.NodeInfo.Name, 0, -1).Val()

	exists := false
	for _, el := range val {
		if el == userId {
			exists = true
			break
		}
	}
	if !exists {
		// 不存在关闭链接
		err := conn.Close()
		if err != nil {
			zap.S().Error(err)
			return
		}

	}
	if GWs == nil {
		GWs = make(map[string]*websocket.Conn)
	}
	if GWsMirror == nil {
		GWsMirror = make(map[*websocket.Conn]string)
	}
	GWs[userId] = conn
	GWsMirror[conn] = userId



	for {
		mt, message, err := conn.ReadMessage()

		zap.S().Infof("消息类型: %d", mt)
		if err != nil {
			log.Println(err)
			break
		}

		err2 := conn.WriteMessage(websocket.TextMessage, []byte(HandlerClientMessage(message, userId)))
		if err2 != nil {
			log.Println("write:", err2)
			break
		}
	}

}
type WsMessage struct {
	Uid     string `json:"uid"`
	Message string `json:"message"`
}

func HandlerClientMessage(p []byte, userId string) string {
	zap.S().Infof("收到客户端消息: %s", p)
	zap.S().Infof("ws客户端id: %s", userId)


	split := strings.Split(userId, "@")
	deviceId := split[0]
	coapMsg := WsMessage{
		Uid:     deviceId,
		Message: string(p),
	}
	jsonData, err := json.Marshal(coapMsg)

	if err != nil {
		zap.S().Error(err)
		return err.Error()
	}

	SetLastOpTime(userId)
	PushToQueue("pre_ws_handler", jsonData)

	return "接收websocket原始数据成功"
}
