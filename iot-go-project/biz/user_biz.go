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

type UserBiz struct{}

func (biz *UserBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.User

	db := glob.GDb

	if name != "" {
		db = db.Where("username like ?", "%"+name+"%")
	}

	db.Model(&models.User{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *UserBiz) FindUser(name string, password string) *models.User {
	var user models.User
	db := glob.GDb
	tx := db.Where("username  = ?", name).Where("password = ?", password).First(&user)
	if tx.Error != nil {
		return nil
	}
	return &user
}
