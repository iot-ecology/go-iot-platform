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

type CalcParamBiz struct{}

func (biz *CalcParamBiz) PageData(name, mqttClientId, signalName, ruleId string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.CalcParam

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	if ruleId != "" {
		db = db.Where("calc_rule_id = ?", ruleId)

	}

	if signalName != "" {
		db = db.Where("signal_id = ?", signalName)
	}
	db.Model(&models.CalcParam{}).Count(&pagination.Total)

	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	for i, calcParam := range dt {
		if calcParam.Protocol == "MQTT" {

			id, err := bizMqtt.FindById(strconv.Itoa(calcParam.DeviceUid))
			if err != nil {
				return nil, err
			}
			if id == nil {
				return nil, fmt.Errorf("no client found for ID: %s", strconv.Itoa(calcParam.DeviceUid))
			}
			dt[i].MqttClientName = id.ClientId
		}

	}
	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
