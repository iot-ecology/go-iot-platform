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
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func SetLastOpTime(uid string) {
	zap.S().Infof("SetLastOpTime 开始, uid = %v", uid)
	// 设置24小时有效时间
	globalRedisClient.Set(context.Background(), "ws:last:"+uid, time.Now().Unix(), 24*time.Hour)
	//globalRedisClient.Set(context.Background(), "ws:last:"+uid, uid, 10*time.Second)

}


func ListenerWs() {
	client := globalRedisClient

	// 配置Redis以启用过期事件的通知。
	client.ConfigSet(context.Background(), "notify-keyspace-events", "Ex")

	// 订阅Redis的过期事件。
	pubsub := client.Subscribe(context.Background(), "__keyevent@"+strconv.Itoa(globalConfig.RedisConfig.Db)+"__:expired")
	// 确保订阅在函数结束时关闭。
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
	}(pubsub)

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			// 如果接收消息时出现错误，记录错误并退出。
			zap.S().Fatalf("Error %s", err)
			return
		}

		// 如果消息负载以 "ws:last:" 开头，则进一步处理。
		if strings.HasPrefix(msg.Payload, "ws:last:") {
			// 分割负载字符串以获取最后一个元素。
			parts := strings.Split(msg.Payload, ":")
			lastElement := parts[len(parts)-1]
			zap.S().Infof("Delete ws conn %s", lastElement)

			if GWs != nil {

				conn := GWs[lastElement]
				if conn != nil {
					conn.Close()
					og := GWsMirror[conn]
					delete(GWs, og)
					delete(GWsMirror, conn)
					globalRedisClient.LRem(context.Background(), "ws_uid:"+globalConfig.NodeInfo.Name, 1, og)
				}
			}

		}
	}
}
