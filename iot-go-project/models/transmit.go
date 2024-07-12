package models

import "gorm.io/gorm"

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
	gorm.Model      `structs:"-"`
	MqttClientId    int  `structs:"mqtt_client_id" json:"mqtt_client_id" ` // 客户端表的外键ID
	MySQLTransmitId uint `struct:"mysql_transmit_id" json:"mysql_transmit_id" gorm:"column:mysql_transmit_id;type:int(
10);" ` // MySQL传输表的外键ID
	Table  string `struct:"table" json:"table" gorm:"column:table;type:varchar(255);"` // 表
	Script string `struct:"script" json:"script" gorm:"column:script"`                 // 转换insert
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
	gorm.Model      `structs:"-"`
	MqttClientId    int    `json:"mqtt_client_id"`                                                              // MQTT客户端表的外键ID
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
	gorm.Model         `structs:"-"`
	MqttClientId       int    `json:"mqtt_client_id"`                                                        // MQTT客户端表的外键ID
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
	gorm.Model           `structs:"-"`
	MqttClientId         int    `json:"mqtt_client_id"`                                                            // MQTT客户端表的外键ID
	ClickhouseTransmitId uint   `json:"clickhouse_transmit_id" gorm:"column:clickhouse_transmit_id;type:int(10);"` // 传输表的外键ID
	Database             string `json:"database" gorm:"column:database;type:varchar(255);"`                        // 数据库
	Script               string `json:"script" gorm:"column:script"`                                               // 转换insert语句的脚本
	Table                string `json:"table" gorm:"column:table;type:varchar(255);"`                              // 表
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
	gorm.Model          `structs:"-"`
	MqttClientId        int    `json:"mqtt_client_id"`                                                          // MQTT客户端表的外键ID
	CassandraTransmitId uint   `json:"cassandra_transmit_id" gorm:"column:cassandra_transmit_id;type:int(10);"` // 传输表的外键ID
	Database            string `json:"database" gorm:"column:database;type:varchar(255);"`                      // 数据库
	Table               string `json:"table" gorm:"column:table;type:varchar(255);"`                            // 表
	Script              string `json:"script" gorm:"column:script"`                                             // 转换insert语句的脚本
	Enable              bool   `json:"enable" gorm:"column:enable;type:tinyint(1);" `                           // 是否启用

}
