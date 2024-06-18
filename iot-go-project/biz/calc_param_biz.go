package biz

import (
	"fmt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type CalcParamBiz struct{}

func (biz *CalcParamBiz) PageData(name, mqttClientId, signalName, ruleId string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var rules []models.CalcParam

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	if mqttClientId != "" {
		db = db.Where("mqtt_client_id = ?", mqttClientId)
	}
	if ruleId != "" {
		db = db.Where("calc_rule_id = ?", ruleId)

	}

	if signalName != "" {
		db = db.Where("signal_id = ?", signalName)
	}
	db.Model(&models.CalcParam{}).Count(&pagination.Total)

	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&rules)

	for i, rule := range rules {
		id, err := bizMqtt.FindById(strconv.Itoa(rule.MqttClientId))
		if err != nil {
			return nil, err
		}
		if id == nil {
			return nil, fmt.Errorf("no client found for ID: %s", strconv.Itoa(rule.MqttClientId))
		}
		rules[i].MqttClientName = id.ClientId
	}
	pagination.Data = rules
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
