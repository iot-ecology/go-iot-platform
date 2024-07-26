package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestInfluxdbQueryM(t *testing.T) {
	var config = InfluxConfig{
		Host:     "127.0.0.1",
		Port:   8086,
		Token:  "mytoken",
		Org:    "myorg",
		Bucket: "mybucket",
	}
	InitInfluxDbClient(config)
	query := `from(bucket: "mybucket")
              |> range(start: -1h, stop: now())
		      |> filter(fn: (r) => r._measurement =~ /^d/)
              |> keep(columns: ["_measurement"])
              |> group()
              |> distinct(column: "_measurement")
              |> limit(n: 200)
              |> sort()`

	queryAPI := GlobalInfluxDbClient.QueryAPI("myorg")

	// 执行查询
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 处理查询结果
	for result.Next() {
		record := result.Record()
		value := record.Value()
		fmt.Printf("Result: %v\n", value)
	}

}
func TestA(t *testing.T) {
	var config = RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Db:       10,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	}
	InitGlobalRedisClient(config)
	globalRedisClient.Set(context.Background(), "a", 1, 0)
	jsonData, _ := json.Marshal(config)

	type Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
		DeviceId string `json:"device_id"`
	}
	auth := Auth{
		Username: "admin",
		Password: "admin",
		DeviceId: "1234567890",
	}
	jsonData, _ = json.Marshal(auth)

	globalRedisClient.HSet(context.Background(), "auth:coap", "1234567890", jsonData)

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
