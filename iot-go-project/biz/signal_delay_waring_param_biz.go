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
	var dt []models.SignalDelayWaringParam

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	if signalDelayWaringId != "" {
		db = db.Where("signal_delay_waring_id=?", signalDelayWaringId)
	}

	db.Model(&models.SignalDelayWaringParam{}).Count(&pagination.Total)

	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	for i, rule := range dt {
		id, err := bizMqtt.FindById(strconv.Itoa(rule.DeviceUid))
		if err != nil {
			return nil, err
		}
		if id == nil {
			return nil, fmt.Errorf("no client found for ID: %s", strconv.Itoa(rule.DeviceUid))
		}
		dt[i].MqttClientName = id.ClientId
	}
	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
