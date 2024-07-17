package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type DeptBiz struct{}

func (biz *DeptBiz) ById(id uint) (*models.Dept, error) {
	var Dept models.Dept

	result := glob.GDb.First(&Dept, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Dept, nil
}
func (biz *DeptBiz) PageData(name, pid string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.Dept

	db := glob.GDb
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if pid != "" {
		db = db.Where("parent_id = ?", pid)
	}
	db.Model(&models.Dept{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)
	for i, dashboard := range dashboards {
		id, err := biz.ById(dashboard.ParentId)
		if err == nil {
			dashboards[i].ParentName = id.Name
		}
	}
	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
