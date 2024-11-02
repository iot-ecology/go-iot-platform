package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type ScriptListBiz struct{}

func (biz *ScriptListBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var feishu []models.ScriptList

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.FeiShu{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&feishu)

	pagination.Data = feishu
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
