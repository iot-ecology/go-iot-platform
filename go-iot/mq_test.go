package main

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestMq(t *testing.T) {
	conn, err := rabbitmq.NewConn("amqp://guest:guest@localhost", rabbitmq.WithConnectionOptionsLogging)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	publisher, err := rabbitmq.NewPublisher(conn, rabbitmq.WithPublisherOptionsLogging, rabbitmq.WithPublisherOptionsExchangeName("pre_handler"), rabbitmq.WithPublisherOptionsExchangeDeclare)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	err = publisher.Publish([]byte("测试数据"), []string{"pre_handler"}, rabbitmq.WithPublishOptionsContentType("application/json"), rabbitmq.WithPublishOptionsExchange("pre_handler"))

	if err != nil {
		log.Println(err)
	}
}

func TestPushMsg(t *testing.T) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
		}
	}(ch)

	mqttMsg := MQTTMessage{
		MQTTClientID: "1",
		Message:      "测试消息",
	}
	jsonData, err := json.Marshal(mqttMsg)
	if err != nil {
		zap.S().Errorf("Error marshalling MQTT message to JSON: %v", err)
		return
	}
	for true {
		err = ch.PublishWithContext(context.Background(), "", "a", // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        jsonData,
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", jsonData)
		//time.Sleep(1 * time.Second)
	}

}

func TestCreateQueueAndEx(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
		}
	}(ch)

	_, err = ch.QueueDeclare("hello", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func TestPushMsg2(t *testing.T) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
		}
	}(ch)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c []CalcParamCache
	c = append(c, CalcParamCache{
		MqttClientId: 2,
		Name:         "p1",
		SignalName:   "Temperature",
		Reduce:       "mean",
		CalcRuleId:   1,
	})

	mqttMsg := CalcCache{
		ID:     1,
		Param:  c,
		Cron:   "0/1 * * * *",
		Script: "function main(map){\n\n\treturn map;\n}",
		Offset: 1000,
	}

	jsonData, err := json.Marshal(mqttMsg)
	if err != nil {
		zap.S().Errorf("Error marshalling MQTT message to JSON: %v", err)
		return
	}
	initGlobalRedisClient(RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Db:       0,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})
	globalRedisClient.HSet(context.Background(), "calc_cache", strconv.Itoa(int(mqttMsg.ID)), jsonData)

	myMap := make(map[string]uint)

	myMap["id"] = mqttMsg.ID

	jsonData, err = json.Marshal(myMap)
	if err != nil {
		zap.S().Errorf("数据异常 %s", err)
		return
	}

	err = ch.PublishWithContext(ctx, "", "calc_queue", // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonData,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", jsonData)

}

type CalcCache struct {
	ID     uint             `json:"id"`
	Param  []CalcParamCache `json:"param"`
	Cron   string           `json:"cron"`
	Script string           `json:"script"`
	Offset int64            `json:"offset"`
}

type CalcParamCache struct {
	MqttClientId int    `json:"mqtt_client_id"`                                        // MQTT客户端表的外键ID
	Name         string `json:"name"`                                                  // 参数名称
	SignalName   string `gorm:"signal_name"  json:"signal_name" structs:"signal_name"` // 信号表 name
	Reduce       string `json:"reduce"`                                                // 数据聚合方式 1. 求和 2. 平均值 3. 最大值 4. 最小值
	CalcRuleId   int    `json:"calc_rule_id"`                                          // CalcRule 主键
}
