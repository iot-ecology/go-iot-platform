package models

import (
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// DingDing 钉钉通知渠道  加签模式
type DingDing struct {
	gorm.Model
	Name        string `json:"name" structs:"name"`
	AccessToken string `json:"access_token" structs:"access_token"`
	Secret      string `json:"secret" structs:"secret"`
	Content     string `json:"content" structs:"content"` // 模板内容
}

// DingDingBindProduct 钉钉通知渠道绑定产品
type DingDingBindProduct struct {
	gorm.Model
	DingDingId int `json:"ding_ding_id" structs:"ding_ding_id"`
	ProductId  int `json:"product_id" structs:"product_id"`
}

// MessageTemplate 推送模板
type MessageTemplate struct {
	gorm.Model
	Content      string  `json:"content" structs:"content"`
	DeviceName   string  `json:"device_name" structs:"device_name"`
	SignalId     int     `json:"signal_id" structs:"signal_id"`           // 信号表的外键ID
	MqttClientId int     `json:"mqtt_client_id" structs:"mqtt_client_id"` // MQTT客户端表的外键ID
	SignalName   string  `json:"signal_name" structs:"signal_name"`
	SignalValue  float64 `json:"signal_value" structs:"signal_value"`
	Min          float64 `json:"min" structs:"min"` // 范围,小值
	Max          float64 `json:"max" structs:"max"` // 范围,大值
	Unit         string  `json:"unit" structs:"unit"`
	InOrOut      int     `gorm:"in_or_out"  json:"in_or_out" structs:"in_or_out"` //范围内报警,范围外报警 1 范围内报警 0 范围外报警

}

func (template *MessageTemplate) Format() string {
	temp := template.Content

	temp = strings.Replace(temp, "{{device_name}}", template.DeviceName, -1)
	temp = strings.Replace(temp, "{{signal_name}}", template.SignalName, -1)
	// 最少保留2位小数
	temp = strings.Replace(temp, "{{min}}", strconv.FormatFloat(template.Min, 'f', -2, 64), -1)
	temp = strings.Replace(temp, "{{max}}", strconv.FormatFloat(template.Max, 'f', -2, 64), -1)
	temp = strings.Replace(temp, "{{signal_value}}", strconv.FormatFloat(template.SignalValue, 'f', -2, 64), -1)
	temp = strings.Replace(temp, "{{unit}}", template.Unit, -1)
	if template.InOrOut == 1 {
		temp = strings.Replace(temp, "{{in_or_out}}", "范围内报警", -1)
	} else if template.InOrOut == 0 {
		temp = strings.Replace(temp, "{{in_or_out}}", "范围外报警", -1)
	}
	return temp
}
