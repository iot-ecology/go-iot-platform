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
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var GRabbitMq *amqp.Connection

func CreateRabbitQueue(queueName string) {

	ch, err := GRabbitMq.Channel()
	if err != nil {
		zap.S().Fatalf("Failed to open a channel %v", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			zap.S().Errorf("Error: %+v", err)

		}
	}(ch)

	_, err = ch.QueueDeclare(queueName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		zap.S().Fatalf("创建queue异常 %s", queueName)
	}
}
func InitRabbitCon() {
	conn, err := amqp.Dial(genUrl())
	if err != nil {
		zap.S().Fatalf("Failed to connect to RabbitMQ  %v", err)
	}

	GRabbitMq = conn

	CreateRabbitQueue("pre_ws_handler")

}
func genUrl() string {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", globalConfig.MQConfig.Username, globalConfig.MQConfig.Password, globalConfig.MQConfig.Host, globalConfig.MQConfig.Port)
	return connStr
}

// PushToQueue 将消息推送到RabbitMQ队列中
//
// 参数：
// queue_name: string类型，目标队列的名称
// body: []byte类型，待发送的消息体
//
// 返回值：
// 无返回值
func PushToQueue(queueName string, body []byte) {

	ch, _ := GRabbitMq.Channel()
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			zap.S().Errorf("Error: %+v", err)

		}
	}(ch)

	_ = ch.PublishWithContext(context.Background(), "", queueName, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	zap.S().Infof(" [x] 发送到 %s 消息体 %s", queueName, body)

}
