package models

import "gorm.io/gorm"

type MySQLTransmit struct {
	gorm.Model `structs:"-"`
	Name       string `json:"name" gorm:"column:name;type:varchar(255);"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);"`
	Username   string `json:"username" gorm:"column:username;type:varchar(255);"`
	Password   string `json:"password" gorm:"column:password;type:varchar(255);"`
	Database   string `json:"database" gorm:"column:database;type:varchar(255);"`
}

type MySQLTransmitBind struct {
	gorm.Model      `structs:"-"`
	MqttClientId    int    `json:"mqtt_client_id"`                                                  // MQTT客户端表的外键ID
	MySQLTransmitId uint   `json:"mysql_transmit_id" gorm:"column:mysql_transmit_id;type:int(10);"` // MySQL传输表的外键ID
	Table           string `json:"table" gorm:"column:table;type:varchar(255);"`                    // 表
	Script          string `json:"script" gorm:"column:script"`                                     // 转换insert语句的脚本
}

type MongoTransmit struct {
	gorm.Model `structs:"-"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);"`
	Username   string `json:"username" gorm:"column:username;type:varchar(255);"`
	Password   string `json:"password" gorm:"column:password;type:varchar(255);"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);"`
}

type MongoTransmitBind struct {
	gorm.Model      `structs:"-"`
	MqttClientId    int    `json:"mqtt_client_id"`                                                              // MQTT客户端表的外键ID
	MongoTransmitId uint   `json:"mongo_transmit_id_transmit_id" gorm:"column:mysql_transmit_id;type:int(10);"` // Mongo传输表的外键ID
	Collection      string `json:"collection" gorm:"column:collection;type:varchar(255);"`                      // 集合表
	Database        string `json:"database" gorm:"column:database;type:varchar(255);"`
	Script          string `json:"script" gorm:"column:script"` // 转换insert语句的脚本
}

type InfluxdbTransmit struct {
	gorm.Model `structs:"-"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);"`
	Token      string `json:"token" gorm:"column:token;type:varchar(255);"`
}

type InfluxdbTransmitBind struct {
	gorm.Model         `structs:"-"`
	MqttClientId       int    `json:"mqtt_client_id"`                                                  // MQTT客户端表的外键ID
	InfluxdbTransmitId uint   `json:"mysql_transmit_id" gorm:"column:mysql_transmit_id;type:int(10);"` // 传输表的外键ID
	Bucket             string `json:"bucket" gorm:"column:bucket;type:varchar(255);"`                  // 桶
	Org                string `json:"org" gorm:"column:org;type:varchar(255);"`                        // 组织
	Measurement        string `json:"measurement" gorm:"column:measurement;type:varchar(255);"`        // 测量
	Script             string `json:"script" gorm:"column:script"`                                     // 转换insert语句的脚本
}

type ClickhouseTransmit struct {
	gorm.Model `structs:"-"`
	Host       string `json:"host" gorm:"column:host;type:varchar(255);"`
	Port       int    `json:"port" gorm:"column:port;type:int(10);"`
	Username   string `json:"username" gorm:"column:username;type:varchar(255);"`
	Password   string `json:"password" gorm:"column:password;type:varchar(255);"`
}

type ClickhouseTransmitBind struct {
	gorm.Model           `structs:"-"`
	MqttClientId         int    `json:"mqtt_client_id"`                                                  // MQTT客户端表的外键ID
	ClickhouseTransmitId uint   `json:"mysql_transmit_id" gorm:"column:mysql_transmit_id;type:int(10);"` // 传输表的外键ID
	Database             string `json:"database" gorm:"column:database;type:varchar(255);"`              // 数据库
	Script               string `json:"script" gorm:"column:script"`                                     // 转换insert语句的脚本
}
