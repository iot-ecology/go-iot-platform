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

package cache

type MySQLTransmitCache struct {
	ID       string `json:"ID"`
	Host     string `json:"host"`
	Port     int    `json:"port" `
	Username string `json:"username" `
	Password string `json:"password" `
	Database string `json:"database" `
	Table    string `json:"table"`
	Script   string `json:"script"`
}

type MongoTransmitCache struct {
	ID         string `json:"ID"`
	Host       string `json:"host"`
	Port       int    `json:"port" `
	Username   string `json:"username" `
	Password   string `json:"password" `
	Database   string `json:"database" `
	Collection string `json:"collection"`
	Script     string `json:"script"`
}

type InfluxTransmitCache struct {
	ID          string `json:"ID"`
	Host        string `json:"host"`
	Port        int    `json:"port" `
	Token       string `json:"token"`
	Bucket      string `json:"bucket"`
	Org         string `json:"org"`
	Measurement string `json:"measurement"`
	Script      string `json:"script"`
}

type ClickhouseTransmitCache struct {
	ID       string `json:"ID"`
	Host     string `json:"host" `
	Port     int    `json:"port" `
	Username string `json:"username" `
	Password string `json:"password" `
	Database string `json:"database" `
	Table    string `json:"table"`
	Script   string `json:"script"`
}

type CassandraTransmitCache struct {
	ID       string `json:"ID"`
	Host     string `json:"host" `
	Port     int    `json:"port" `
	Username string `json:"username" `
	Password string `json:"password" `
	Database string `json:"database" `
	Table    string `json:"table"`
	Script   string `json:"script"`
}

type KafkaTransmitCache struct {
	ID     string `json:"ID"`
	Host   string `json:"host"`
	Port   int    `json:"port" `
	Topic  string `json:"topic"`
	Script string `json:"script"`
}
type RabbitTransmitCache struct {
	ID         string `json:"ID"`
	Host       string `json:"host"`
	Port       int    `json:"port" `
	Script     string `json:"script"`
	Exchange   string `json:"exchange" `    // 交换机
	RoutingKey string `json:"routing_key" ` // 路由键

}
