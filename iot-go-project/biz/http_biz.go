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
	"context"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type HttpHandlerBiz struct{}

func (biz *HttpHandlerBiz) ById(id uint) (*models.HttpHandler, error) {
	var HttpHandler models.HttpHandler

	result := glob.GDb.First(&HttpHandler, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &HttpHandler, nil
}

func (biz *HttpHandlerBiz) PageData(name, device_info_id string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var HttpHandlerList []models.HttpHandler

	db := glob.GDb
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if device_info_id != "" {
		db = db.Where("device_info_id = ?", device_info_id)
	}

	db.Model(&models.HttpHandler{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&HttpHandlerList)
	pagination.Data = HttpHandlerList
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *HttpHandlerBiz) SetRedis(data models.HttpHandler) {
	glob.GRedis.HSet(context.Background(), "struct:Http", strconv.Itoa(int(data.DeviceInfoId)), data.Script)
}
func (biz *HttpHandlerBiz) RemoveRedis(data models.HttpHandler) {
	glob.GRedis.HDel(context.Background(), "struct:Http", strconv.Itoa(int(data.DeviceInfoId)))
}
