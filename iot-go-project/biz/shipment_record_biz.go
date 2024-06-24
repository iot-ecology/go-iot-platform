package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type ShipmentRecordBiz struct{}

func (biz *ShipmentRecordBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.ShipmentRecord

	db := glob.GDb

	db.Model(&models.ShipmentRecord{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
