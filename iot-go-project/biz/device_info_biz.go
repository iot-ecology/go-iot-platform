package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type DeviceInfoBiz struct{}

func (biz *DeviceInfoBiz) PageData(sn string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.DeviceInfo

	db := glob.GDb

	if sn != "" {
		db = db.Where("sn like ?", "%"+sn+"%")
	}

	db.Model(&models.DeviceInfo{}).Count(&pagination.Total)
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
