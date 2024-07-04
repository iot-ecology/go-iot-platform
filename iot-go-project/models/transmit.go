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
	Script          string `json:"script" gorm:"column:script;type:varchar(255);"`                  // 转换insert语句的脚本
}

type MongoTransmit struct {
	gorm.Model `structs:"-"`
}
