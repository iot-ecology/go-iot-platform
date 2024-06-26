package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type RoleBiz struct{}

func (biz *RoleBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.Role

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.Role{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *RoleBiz) FindByUserId(userId uint) []uint {
	var roles []models.UserRole
	db := glob.GDb
	db.Where("user_id = ?", userId).Find(&roles)
	var roleIds []uint
	for _, v := range roles {
		roleIds = append(roleIds, v.RoleId)
	}
	return roleIds

}
