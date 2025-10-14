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
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"iot-transmit/cache"
	"iot-transmit/common"
)

var transmitCacheBiz = cache.TransmitCacheBiz{}

func HandlerTransmit(messages <-chan amqp.Delivery) {

	go func() {

		for d := range messages {
			var data []common.DataRowList

			err := json.Unmarshal(d.Body, &data)

			if err != nil {
				zap.S().Error("处理cassandra数据失败", zap.Error(err))
			}

			transmitCacheBiz.Run(globalRedisClient, data)
			err = d.Ack(false)
			if err != nil {
				zap.S().Errorf("消息确认异常：%+v", err)

			}
		}
	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}
