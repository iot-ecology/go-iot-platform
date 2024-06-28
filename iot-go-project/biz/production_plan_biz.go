package biz

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
	"time"
)

type ProductionPlanBiz struct{}

// BeforeCreate 创建生产计划后的任务
func (biz *ProductionPlanBiz) BeforeCreate(param models.ProductionPlan) {
	// 根据结束时间发送 计划到期通知
	set_expire_key(glob.StartNotification, param)
	set_expire_key(glob.DueSoonNotification, param)
	set_expire_key(glob.DueNotification, param)
}

// set_expire_key 根据传入的消息类型和参数设置 Redis 中的键的过期时间
//
// 参数：
//
//	mt: glob.MessageType 类型的消息类型
//	param: models.ProductionPlan 类型的生产计划参数
//
// 返回值：
//
//	无返回值
func set_expire_key(mt glob.MessageType, param models.ProductionPlan) {
	id := strconv.Itoa(int(param.ID))
	key := mt.String() + ":" + id
	glob.GRedis.SetNX(context.Background(), key, id, 0)
	switch mt {
	case glob.StartNotification:
		glob.GRedis.ExpireAt(context.Background(), key, param.StartDate)
	case glob.DueSoonNotification:
		glob.GRedis.ExpireAt(context.Background(), key, param.EndDate.AddDate(0, 0, -3))
	case glob.DueNotification:
		glob.GRedis.ExpireAt(context.Background(), key, param.EndDate)
	default:
		panic("unhandled default case")
	}
}

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
	now := time.Now()
	if param.Status == "已完成" {

		var pp []models.ProductPlan

		tx.Where("production_plan_id = ?", ProductionPlan.ID).Find(pp)

		for _, plan := range pp {
			// 更新产品库存
			db := tx.Model(&models.Product{}).Where("id = ?", plan.ProductID).Update("quantity", gorm.Expr("quantity + ?", plan.Quantity))
			if db.Error != nil {
				zap.S().Errorf("更新产品库存失败 %+v", db.Error)
				tx.Rollback()
				return false
			}

			var product models.Product

			find := tx.Model(&models.Product{}).Where("id = ?", plan.ProductID).Find(product)
			if find.Error != nil {
				tx.Rollback()
				return false
			}

			// 创建新的设备
			for range plan.Quantity {

				// 创建新的设备
				info := models.DeviceInfo{
					ProductId:         plan.ProductID,
					SN:                uuid.New().String(), // fixme： 生成设备SN
					ManufacturingDate: now,
					Source:            1,
					WarrantyExpiry:    now.AddDate(0, 0, product.WarrantyPeriod),
				}

				create := tx.Model(&models.DeviceInfo{}).Create(&info)
				if create.Error != nil {
					zap.S().Errorf("创建设备失败 %+v", create.Error)
					tx.Rollback()
					return false
				}
			}
		}

	}

	if err := tx.Commit().Error; err != nil {
		zap.S().Errorf("更新生产计划状态失败 %+v", err)
		tx.Rollback()
		return false
	}
	return true

}
