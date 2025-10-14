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

type DeptBiz struct{}

func (biz *DeptBiz) ById(id uint) (*models.Dept, error) {
	var Dept models.Dept

	result := glob.GDb.First(&Dept, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Dept, nil
}
func (biz *DeptBiz) PageData(name, pid string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var deptList []models.Dept

	db := glob.GDb
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if pid != "" {
		db = db.Where("parent_id = ?", pid)
	}
	db.Model(&models.Dept{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&deptList)
	for i, dept := range deptList {
		id, err := biz.ById(dept.ParentId)
		if err == nil {
			deptList[i].ParentName = id.Name
		}


	}
	pagination.Data = deptList
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
