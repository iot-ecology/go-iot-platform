package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type CalcRuleBiz struct{}

func (biz *CalcRuleBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var rules []models.CalcRule

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	db.Model(&models.CalcRule{}).Count(&pagination.Total) // 计算总记录数

	offset := (page - 1) * size
	db = db.Offset(offset).Limit(size).Find(&rules)

	pagination.Data = rules
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
