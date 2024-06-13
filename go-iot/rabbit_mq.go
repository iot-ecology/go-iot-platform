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

	CreateRabbitQueue("pre_handler")

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
