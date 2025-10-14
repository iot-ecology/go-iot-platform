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

type TcpHandlerBiz struct{}

func (biz *TcpHandlerBiz) ById(id uint) (*models.TcpHandler, error) {
	var TcpHandler models.TcpHandler

	result := glob.GDb.First(&TcpHandler, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &TcpHandler, nil
}

func (biz *TcpHandlerBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var TcpHandlerList []models.TcpHandler

	db := glob.GDb
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	db.Model(&models.TcpHandler{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&TcpHandlerList)
	pagination.Data = TcpHandlerList
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *TcpHandlerBiz) SetRedis(data models.TcpHandler) {
	glob.GRedis.HSet(context.Background(), "struct:tcp", strconv.Itoa(int(data.DeviceInfoId)), data.Script)
}
func (biz *TcpHandlerBiz) RemoveRedis(data models.TcpHandler) {
	glob.GRedis.HDel(context.Background(), "struct:tcp", strconv.Itoa(int(data.DeviceInfoId)))
}
