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

package biz

import (
	"fmt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type SignalDelayWaringParamBiz struct{}

func (biz *SignalDelayWaringParamBiz) PageData(name, signalDelayWaringId string, page, size int, device_uid, signal_id string) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.SignalDelayWaringParam

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	if signalDelayWaringId != "" {
		db = db.Where("signal_delay_waring_id=?", signalDelayWaringId)
	}

	if device_uid !=""{
		db = db.Where("device_uid = ? ",device_uid)
	}
	if signal_id !=""{
		db = db.Where("signal_id = ? ",signal_id)
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
