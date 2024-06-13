package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type DashboardBiz struct{}

func (biz *DashboardBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.Dashboard

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.Dashboard{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db = db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
