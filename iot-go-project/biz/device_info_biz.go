package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type DeviceInfoBiz struct{}

var productBiz = ProductBiz{}

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

	for i, info := range dt {
		dt[i].ProductName = productBiz.FindById(info.ProductId).Name
	}
	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *DeviceInfoBiz) FindById(id uint) *models.DeviceInfo {
	var dt models.DeviceInfo
	db := glob.GDb
	db.Where("id = ?", id).Find(&dt)
	return &dt
}
