package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type DeviceInstallRecordBiz struct{}

func (biz *DeviceInstallRecordBiz) PageData(sn string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.DeviceInstallRecord

	db := glob.GDb

	db.Model(&models.DeviceInstallRecord{}).Count(&pagination.Total)
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
