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
