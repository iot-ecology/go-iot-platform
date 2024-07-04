package transmit

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type MySQLTransmitBiz struct{}

func (biz *MySQLTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var MySQLTransmits []models.MySQLTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.MySQLTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&MySQLTransmits)

	pagination.Data = MySQLTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
