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
	"iot-transmit/cache"
	"strconv"
)

type RabbitmqTransmitBindBiz struct{}

func (biz *RabbitmqTransmitBindBiz) PageData(exchange string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.RabbitmqTransmitBind

	db := glob.GDb

	if exchange != "" {
		db = db.Where("exchange like ?", "%"+exchange+"%")
	}

	db.Model(&models.RabbitmqTransmitBind{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *RabbitmqTransmitBindBiz) Bind(req models.RabbitmqTransmitBind) {
	if req.Enable {
		jsonData := biz.toByte(req)
		// 缓存构造
		glob.GRedis.LPush(context.Background(), "transmit:Rabbit:" + req.Protocol +":"+strconv.Itoa(req.DeviceUid) +":" +req.
			IdentificationCode,
			jsonData)
	}
}

func (biz *RabbitmqTransmitBindBiz) toByte(req models.RabbitmqTransmitBind) []byte {
	var RabbitInfo models.RabbitmqTransmit

	glob.GDb.First(&RabbitInfo, req.RabbitmqTransmitId)

	v := cache.RabbitTransmitCache{
		ID:         "Rabbit-" + strconv.Itoa(int(req.ID)),
		Host:       RabbitInfo.Host,
		Port:       RabbitInfo.Port,
		Script:     req.Script,
		Exchange:   req.Exchange,
		RoutingKey: req.RoutingKey,
	}
	jsonData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// ChangeEnable 修改启用状态
func (biz *RabbitmqTransmitBindBiz) ChangeEnable(req models.RabbitmqTransmitBind) {
	glob.GDb.Model(models.RabbitmqTransmitBind{}).Where("id = ?", req.ID).Update("enable", req.Enable)
	biz.HandlerRedis(req)
}

func (biz *RabbitmqTransmitBindBiz) HandlerRedis(req models.RabbitmqTransmitBind) {
	if req.Enable {
		biz.Bind(req)
	} else {
		glob.GRedis.LRem(context.Background(), "transmit:Rabbit:"+ req.Protocol+":"+strconv.Itoa(req.DeviceUid) +":" +req.
			IdentificationCode, 1,
			biz.toByte(req))
	}
}
