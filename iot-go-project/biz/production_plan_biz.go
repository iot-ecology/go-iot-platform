/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package biz

import (
	"context"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductionPlanBiz struct{}

// BeforeCreate 创建生产计划后的任务
func (biz *ProductionPlanBiz) BeforeCreate(param models.ProductionPlan) {
	// 根据结束时间发送 计划到期通知
	setExpireKeyWithProductionPlan(glob.StartNotification, param)
	setExpireKeyWithProductionPlan(glob.DueSoonNotification, param)
	setExpireKeyWithProductionPlan(glob.DueNotification, param)
}

// setExpireKeyWithProductionPlan 根据传入的消息类型和参数设置 Redis 中的键的过期时间
//
// 参数：
//
//	mt: glob.MessageType 类型的消息类型
//	param: models.ProductionPlan 类型的生产计划参数
//
// 返回值：
//
//	无返回值
func setExpireKeyWithProductionPlan(mt glob.MessageType, param models.ProductionPlan) {
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
	var productionPlan models.ProductionPlan

	result := glob.GDb.First(&productionPlan, param.ID)

	if result.Error != nil {
		return false
	}

	if productionPlan.Status == "1" && param.Status == "3" {
		return false
	}
	if productionPlan.Status == "2" && param.Status == "1" {
		return false
	}
	tx := glob.GDb.Begin()
	if tx.Error != nil {
		return false
	}

	update := tx.Model(&productionPlan).Update("status", param.Status)
	if update.Error != nil {
		tx.Rollback()
		return false
	}
	now := time.Now()
	if param.Status == "3" {

		var pp []models.ProductPlan

		tx.Where("production_plan_id = ?", productionPlan.ID).Find(&pp)

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
				date := now.AddDate(0, 0, product.WarrantyPeriod)
				// 创建新的设备
				info := models.DeviceInfo{
					ProductId:         plan.ProductID,
					SN:                uuid.New().String(), // fixme： 生成设备SN
					ManufacturingDate: now,
					Source:            1,
					WarrantyExpiry:    date,
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
