package biz

import (
	"fmt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type SignalDelayWaringParamBiz struct{}

func (biz *SignalDelayWaringParamBiz) PageData(name, signalDelayWaringId string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var rules []models.SignalDelayWaringParam

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	if signalDelayWaringId != "" {
		db = db.Where("signal_delay_waring_id=?", signalDelayWaringId)
	}

	db.Model(&models.SignalDelayWaringParam{}).Count(&pagination.Total)

	offset := (page - 1) * size
	db = db.Offset(offset).Limit(size).Find(&rules)

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
