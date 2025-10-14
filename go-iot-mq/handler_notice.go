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
	"github.com/CatchZeng/feishu/pkg/feishu"
	"github.com/blinkbean/dingtalk"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"iot-notice/models"
)

func HandlerNotice(messages <-chan amqp.Delivery) {

	go func() {

		for d := range messages {
			var data models.MessageTemplate

			err := json.Unmarshal(d.Body, &data)

			if err != nil {
				zap.S().Error("处理通知数据失败", zap.Error(err))
			}

			val := globalRedisClient.LRange(context.Background(), "mqtt_client_id_bind_product:"+data.DeviceUid, 0,
				-1).Val()

			// 检索钉钉通道

			dingding_channel := globalRedisClient.HMGet(context.Background(), "message_channel:dingding", val...).Val()
			for _, str := range dingding_channel {
				ding := str.(DingDing)
				data.Content = ding.Content

				cli := dingtalk.InitDingTalkWithSecret(ding.AccessToken, ding.Secret)

				cli.SendTextMessage(data.Format())
			}

			// 检索飞书通道

			val = globalRedisClient.LRange(context.Background(), "message_channel_bind:feishu:"+data.DeviceUid, 0, -1).Val()

			dingding_channel = globalRedisClient.HMGet(context.Background(), "message_channel:feishu", val...).Val()
			for _, str := range dingding_channel {
				ding := str.(FeiShu)
				data.Content = ding.Content

				client := feishu.NewClient(ding.AccessToken, ding.Secret)
				text := feishu.NewText(data.Format())
				line := []feishu.PostItem{text}
				msg := feishu.NewPostMessage()

				msg.SetZHTitle("告警通知").
					AppendZHContent(line)
				client.Send(msg)
			}

			err = d.Ack(false)
			if err != nil {
				zap.S().Errorf("消息确认异常：%+v", err)

			}
		}
	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}
