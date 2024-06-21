package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type RepairRecordBiz struct{}

func (biz *RepairRecordBiz) PageData(sn string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.RepairRecord

	db := glob.GDb

	db.Model(&models.Product{}).Count(&pagination.Total)
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
