package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type UserBiz struct{}

func (biz *UserBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.User

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.User{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *UserBiz) FindUser(name string, password string) *models.User {
	var user models.User
	db := glob.GDb
	tx := db.Where("username  = ?", name).Where("password = ?", password).First(&user)
	if tx.Error != nil {
		return nil
	}
	return &user
}
