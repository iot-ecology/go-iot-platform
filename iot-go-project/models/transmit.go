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

package models

import "gorm.io/gorm"

/****

fixme: MqttClientId 字段语义存在歧义，应该修订为两个字段
*/

type MySQLTransmit struct {
	gorm.Model `structs:"-"`
	Name       string `json:"name" gorm:"column:name;type:varchar(255);" structs:"name"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);" structs:"host"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);" structs:"port"`
	Username   string `json:"username" gorm:"column:username;type:varchar(255);" structs:"username"`
	Password   string `json:"password" gorm:"column:password;type:varchar(255);" structs:"password"`
	Database   string `json:"database" gorm:"column:database;type:varchar(255);" structs:"database"`
}

type MySQLTransmitBind struct {
	gorm.Model         `structs:"-"`
	Protocol           string `json:"protocol"`                                                                                    // 协议
	DeviceUid          int `json:"device_uid"`                                                                                  // 设备UID
	IdentificationCode string `json:"identification_code"`                                                                         // 设备标识码
	MySQLTransmitId    uint   `struct:"mysql_transmit_id" json:"mysql_transmit_id" gorm:"column:mysql_transmit_id;type:int(10);" ` // MySQL传输表的外键ID
	Table              string `struct:"table" json:"table" gorm:"column:table_name;type:varchar(
255);"`                                 // 表
	Script             string `struct:"script" json:"script" gorm:"column:script"`                                                 // 转换insert
	// 语句的脚本
	Enable bool `structs:"enable" json:"enable" gorm:"column:enable;type:tinyint(1);"` // 是否启用
}

type MongoTransmit struct {
	gorm.Model `structs:"-"`
	Name       string `structs:"name" json:"name" gorm:"column:name;type:varchar(255);"`
	Host       string `structs:"host" json:"host" gorm:"column:host;type:varchar(255);"`
	Username   string `structs:"username" json:"username" gorm:"column:username;type:varchar(255);"`
	Password   string `structs:"password" json:"password" gorm:"column:password;type:varchar(255);"`
	Port       int    `structs:"port" json:"port" gorm:"column:port;type:int(10);"`
}

type MongoTransmitBind struct {
	gorm.Model `structs:"-"`
	DeviceUid  int    `json:"device_uid"` // 设备UID
	Protocol   string `json:"protocol"`   // 协议

	IdentificationCode string `json:"identification_code"` // 设备标识码

	MongoTransmitId uint   `json:"mongo_transmit_id_transmit_id" gorm:"column:mysql_transmit_id;type:int(10);"` // Mongo传输表的外键ID
	Collection      string `json:"collection" gorm:"column:collection;type:varchar(255);"`                      // 集合表
	Database        string `json:"database" gorm:"column:database;type:varchar(255);"`
	Script          string `json:"script" gorm:"column:script"`                  // 转换insert语句的脚本
	Enable          bool   `json:"enable" gorm:"column:enable;type:tinyint(1);"` // 是否启用

}

type InfluxdbTransmit struct {
	gorm.Model `structs:"-"`
	Name       string `structs:"name" json:"name" gorm:"column:name;type:varchar(255);"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);"`
	Token      string `json:"token" gorm:"column:token;type:varchar(255);"`
}

type InfluxdbTransmitBind struct {
	gorm.Model `structs:"-"`
	DeviceUid  int    `json:"device_uid"` // 设备UID
	Protocol   string `json:"protocol"`   // 协议

	IdentificationCode string `json:"identification_code"` // 设备标识码

	InfluxdbTransmitId uint   `json:"influxdb_transmit_id" gorm:"column:influxdb_transmit_id;type:int(10);"` // 传输表的外键ID
	Bucket             string `json:"bucket" gorm:"column:bucket;type:varchar(255);"`                        // 桶
	Org                string `json:"org" gorm:"column:org;type:varchar(255);"`                              // 组织
	Measurement        string `json:"measurement" gorm:"column:measurement;type:varchar(255);"`              // 测量
	Script             string `json:"script" gorm:"column:script"`                                           // 转换insert语句的脚本
	Enable             bool   `json:"enable" gorm:"column:enable;type:tinyint(1);"`                          // 是否启用

}

type ClickhouseTransmit struct {
	gorm.Model `structs:"-"`
	Name       string `structs:"name" json:"name" gorm:"column:name;type:varchar(255);"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);"`
	Username   string `json:"username" gorm:"column:username;type:varchar(255);"`
	Password   string `json:"password" gorm:"column:password;type:varchar(255);"`
}

type ClickhouseTransmitBind struct {
	gorm.Model `structs:"-"`
	DeviceUid  int    `json:"device_uid"` // 设备UID
	Protocol   string `json:"protocol"`   // 协议

	IdentificationCode string `json:"identification_code"` // 设备标识码

	ClickhouseTransmitId uint   `json:"clickhouse_transmit_id" gorm:"column:clickhouse_transmit_id;type:int(10);"` // 传输表的外键ID
	Database             string `json:"database" gorm:"column:database;type:varchar(255);"`                        // 数据库
	Script               string `json:"script" gorm:"column:script"`                                               // 转换insert语句的脚本
	Table                string `json:"table" gorm:"column:table_name;type:varchar(255);"`                              // 表
	Enable               bool   `json:"enable" gorm:"column:enable;type:tinyint(1);" `                             // 是否启用

}

type CassandraTransmit struct {
	gorm.Model `structs:"-"`
	Name       string `structs:"name" json:"name" gorm:"column:name;type:varchar(255);"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);"`
	Username   string `json:"username" gorm:"column:username;type:varchar(255);"`
	Password   string `json:"password" gorm:"column:password;type:varchar(255);"`
}

type CassandraTransmitBind struct {
	gorm.Model `structs:"-"`
	DeviceUid  int    `json:"device_uid"` // 设备UID
	Protocol   string `json:"protocol"`   // 协议

	IdentificationCode string `json:"identification_code"` // 设备标识码

	CassandraTransmitId uint   `json:"cassandra_transmit_id" gorm:"column:cassandra_transmit_id;type:int(10);"` // 传输表的外键ID
	Database            string `json:"database" gorm:"column:database;type:varchar(255);"`                      // 数据库
	Table               string `json:"table" gorm:"column:table_name;type:varchar(255);"`                            // 表
	Script              string `json:"script" gorm:"column:script"`                                             // 转换insert语句的脚本
	Enable              bool   `json:"enable" gorm:"column:enable;type:tinyint(1);" `                           // 是否启用

}

type RabbitmqTransmit struct {
	gorm.Model `structs:"-"`
	Name       string `structs:"name" json:"name" gorm:"column:name;type:varchar(255);"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);"`
	Username   string `json:"username" gorm:"column:username;type:varchar(255);"`
	Password   string `json:"password" gorm:"column:password;type:varchar(255);"`
}

type RabbitmqTransmitBind struct {
	gorm.Model `structs:"-"`
	DeviceUid  int    `json:"device_uid"` // 设备UID
	Protocol   string `json:"protocol"`   // 协议

	IdentificationCode string `json:"identification_code"` // 设备标识码

	RabbitmqTransmitId uint   `json:"rabbitmq_transmit_id" gorm:"column:rabbitmq_transmit_id;type:int(10);"` // 传输表
	Exchange           string `json:"exchange" gorm:"column:exchange;type:varchar(255);"`                    // 交换机
	RoutingKey         string `json:"routing_key" gorm:"column:routing_key;type:varchar(255);"`              // 路由键
	Script             string `json:"script" gorm:"column:script"`                                           // 转换insert语句的脚本
	Enable             bool   `json:"enable" gorm:"column:enable;type:tinyint(1);" `                         // 是否启用

}

type KafkaTransmit struct {
	gorm.Model `structs:"-"`
	Name       string `structs:"name" json:"name" gorm:"column:name;type:varchar(255);"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);"`
}
type KafkaTransmitBind struct {
	gorm.Model `structs:"-"`
	DeviceUid  int    `json:"device_uid"` // 设备UID
	Protocol   string `json:"protocol"`   // 协议

	IdentificationCode string `json:"identification_code"` // 设备标识码

	KafkaTransmitId uint   `json:"kafka_transmit_id" gorm:"column:kafka_transmit_id;type:int(10);"` // 传输表
	Topic           string `json:"topic" gorm:"column:topic;type:varchar(255);"`                    // topic
	Script          string `json:"script" gorm:"column:script"`                                     // 转换insert语句的脚本
	Enable          bool   `json:"enable" gorm:"column:enable;type:tinyint(1);" `                   // 是否启用

}
