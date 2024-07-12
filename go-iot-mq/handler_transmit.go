package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func HandlerTransmit(messages <-chan amqp.Delivery) {

	go func() {

		for d := range messages {
			err := d.Ack(false)
			if err != nil {
				zap.S().Errorf("消息确认异常：%+v", err)

			}
		}
	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}
