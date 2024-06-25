package biz

import (
	"gorm.io/gorm"
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type ProductionPlanBiz struct{}

func (biz *ProductionPlanBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.ProductionPlan

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.ProductionPlan{}).Count(&pagination.Total)
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *ProductionPlanBiz) ChangeProductionPlanState(param servlet.ProductionPlanChangeParam) bool {
	var ProductionPlan models.ProductionPlan

	result := glob.GDb.First(&ProductionPlan, param.ID)

	if result.Error != nil {
		return false
	}

	if ProductionPlan.Status == "准备中" && param.Status == "已完成" {
		return false
	}
	if ProductionPlan.Status == "进行中" && param.Status == "准备中" {
		return false
	}
	tx := glob.GDb.Begin()
	if tx.Error != nil {
		return false
	}

	update := tx.Model(&ProductionPlan).Update("status", param.Status)
	if update.Error != nil {
		tx.Rollback()
		return false
	}

	if param.Status == "已完成" {

		var pp []models.ProductPlan

		tx.Where("production_plan_id = ?", ProductionPlan.ID).Find(pp)

		for _, plan := range pp {
			db := tx.Model(&models.Product{}).Where("id = ?", plan.ProductID).Update("quantity", gorm.Expr("quantity + ?", plan.Quantity))
			if db.Error != nil {
				tx.Rollback()
				return false
			}
		}

	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false
	}
	return true

}
