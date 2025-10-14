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

package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type InfluxdbTransmitBiz struct{}

func (biz *InfluxdbTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var influxdbTransmits []models.InfluxdbTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.InfluxdbTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&influxdbTransmits)

	pagination.Data = influxdbTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *InfluxdbTransmitBiz) SetRedis(param models.InfluxdbTransmit) {
	jsonData, err := json.Marshal(param)

	if err != nil {
		panic(err)
	}
	glob.GRedis.HSet(context.Background(), "transmit:influxdb:"+strconv.Itoa(int(param.ID)), jsonData)
}

func (biz *InfluxdbTransmitBiz) DeleteRedis(param models.InfluxdbTransmit) {
	glob.GRedis.HDel(context.Background(), "transmit:influxdb:"+strconv.Itoa(int(param.ID)))
}

