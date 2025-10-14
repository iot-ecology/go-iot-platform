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

type RoleBiz struct{}

func (biz *RoleBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.Role

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.Role{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *RoleBiz) FindByUserId(userId uint) []uint {
	var roles []models.UserRole
	db := glob.GDb
	db.Where("user_id = ?", userId).Find(&roles)
	var roleIds []uint
	for _, v := range roles {
		roleIds = append(roleIds, v.RoleId)
	}
	return roleIds

}
