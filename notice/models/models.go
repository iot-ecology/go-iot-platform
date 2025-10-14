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

import (
	"strconv"
	"strings"
)

// MessageTemplate 推送模板
type MessageTemplate struct {
	Content      string  `json:"content" structs:"content"`
	DeviceUid string `json:"device_uid" structs:"device_uid"`

	GeneratorTime int64 `json:"generator_time" structs:"generator_time"` // 数据生产时间
	DeviceName   string  `json:"device_name" structs:"device_name"`
	SignalId     int     `json:"signal_id" structs:"signal_id"`           // 信号表的外键ID
	MqttClientId string     `json:"mqtt_client_id" structs:"mqtt_client_id"` // MQTT客户端表的外键ID
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
