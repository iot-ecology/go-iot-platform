package main

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
)

func TestA(t *testing.T) {
	var config = RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Db:       1,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	}
	InitGlobalRedisClient(config)
	globalRedisClient.Set(context.Background(), "a", 1, 0)
	jsonData, _ := json.Marshal(config)

	globalRedisClient.HSet(context.Background(), "auth:http",strconv.Itoa(int(1)), jsonData)


}

func TestMqCustomer(t *testing.T) {

	var config = MQConfig{
		Host:     "127.0.0.1",
		Port:     5672,
		Password: "guest",
		Username: "guest",
	}
	InitRabbitCon(config)
	ch, err := GRabbitMq.Channel()
	failOnError(err, "Failed to open a channel")

	preHandlerMessage, err := ch.Consume("a", // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		true,  // no-wait
		nil,   // args
	)
	if err != nil {
		failOnError(err, "Failed to register a consumer")
	}
	go func() {
		for d := range preHandlerMessage {
			t.Log(string(d.Body))
		}
	}()
}
