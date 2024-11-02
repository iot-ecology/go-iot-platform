package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type ShipmentRecordBiz struct{}

func (biz *ShipmentRecordBiz) PageData(customerName, status string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.ShipmentRecord

	db := glob.GDb

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if customerName != "" {
		db = db.Where("customer_name like ?", "%"+customerName+"%")
	}
	db.Model(&models.ShipmentRecord{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
