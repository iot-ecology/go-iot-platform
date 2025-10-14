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
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
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


func TestMongo(t *testing.T){

	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%d", url.QueryEscape("admin"),
		url.QueryEscape("admin"), "127.0.0.1", 27017)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)

	db := client.Database("iot")

	prefix := ""

	// 构建正则表达式，匹配以prefix开头的集合名称
	regex := primitive.Regex{Pattern: "^" + prefix, Options: "i"} // 'i' 表示不区分大小写

	// 构建查询条件
	filter := bson.M{"name": regex}
	// 检查集合是否存在
	collectionNames, err := db.ListCollectionNames(context.TODO(),filter)
	if err != nil {
		log.Fatal(err)
	}

	collectionExists := false
	for _, name := range collectionNames {
		if name == "你的集合名" {
			collectionExists = true
			break
		}
		log.Println(name)
	}
	log.Println(collectionExists)

	db.CreateCollection(context.TODO(), "aaaaaaaaa")
}

func TestInfluxDbCreateBucket(t *testing.T){
	var config = InfluxConfig{
		Host:     "127.0.0.1",
		Port:   8086,
		Token:  "mytoken",
		Org:    "myorg",
		Bucket: "mybucket",
	}
	InitInfluxDbClient(config)
	name, err := GlobalInfluxDbClient.OrganizationsAPI().FindOrganizationByName(context.Background(), config.Org)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(name)
	id, err := GlobalInfluxDbClient.BucketsAPI().FindBucketsByOrgID(context.Background(), *name.Id)
	if err != nil {
		log.Fatal(err)
	}

	// 假设id是一个存储桶ID的切片
	for _, bucketID := range *id {
		log.Println(bucketID.Name)
	}

	GlobalInfluxDbClient.BucketsAPI().CreateBucketWithName(context.Background(), name, "cli_create")




}