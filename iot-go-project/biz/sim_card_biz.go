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
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type SimCardBiz struct{}

func (biz *SimCardBiz) PageData(accessNumber ,iccid string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.SimCard

	db := glob.GDb
	if accessNumber != "" {
		db = db.Where("access_number like ?", "%"+accessNumber+"%")
	}
	if iccid !=""{
		db = db.Where("iccid like ?" ,"%" + iccid +"%")
	}
	db.Model(&models.SimCard{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

var deviceInfoBiz = DeviceInfoBiz{}

// func (biz *SimCardBiz) beforeCreate(card models.SimCard) {
// 	id := strconv.Itoa(int(card.ID))
// 	key := glob.SimCardExpireTime.String() + ":" + id
// 	glob.GRedis.SetNX(context.Background(), key, id, 0)
// 	glob.GRedis.ExpireAt(context.Background(), key, card.Expiration.AddDate(0, 0, -3))

// }

func (biz *SimCardBiz) PageHistory(simCardId string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var useHistories []models.SimUseHistory

	db := glob.GDb

	if simCardId != "" {
		db = db.Where("sim_id = ?", simCardId)
	}

	db.Model(&models.SimUseHistory{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&useHistories)

	var resp []servlet.SimUseHistoryResp

	for _, dashboard := range useHistories {
		var id = deviceInfoBiz.FindById(dashboard.DeviceInfoId)

		historyResp := servlet.SimUseHistoryResp{
			ID:           dashboard.ID,
			SimId:        dashboard.SimId,
			DeviceInfoId: dashboard.DeviceInfoId,
			Description:  dashboard.Description,
		}
		if id != nil {
			historyResp.SN = id.SN
		}

		resp = append(resp, historyResp)
	}

	pagination.Data = resp
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
