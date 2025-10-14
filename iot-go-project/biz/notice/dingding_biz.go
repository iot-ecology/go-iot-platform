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

package notice

import (
	"context"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"

	"go.uber.org/zap"
)

type DingDingBiz struct{}

func (biz *DingDingBiz) PageData(name, token, cot string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dingding []models.DingDing

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	if token != "" {
		db = db.Where("token like ?", "%"+token +"%")
	}
	if cot != "" {
		db = db.Where("content like ?", "%"+cot +"%")
	}

	db.Model(&models.DingDing{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dingding)

	pagination.Data = dingding
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *DingDingBiz) Bind(req []models.DingDingBindProduct) bool {
	param := req[0]
	db := glob.GDb

	tx := db.Begin()
	if tx.Error != nil {
		return false

	}

	var toDelete []models.DingDingBindProduct
	tx.Where("product_id = ?", param.ProductId).Find(&toDelete)
	result := tx.Where("product_id = ?", param.ProductId).Delete(models.DingDingBindProduct{})
	if result.Error != nil {
		// 如果出现错误，回滚事务
		tx.Rollback()
		return false

	}
	result = tx.Model(&models.DingDingBindProduct{}).CreateInBatches(req, len(req))

	if result.Error != nil {
		tx.Rollback()
		zap.S().Infoln("Error occurred during creation:", result.Error)
		return false
	}
	if err := tx.Commit().Error; err != nil {
		return false

	}

	for _, product := range toDelete {
		glob.GRedis.Del(context.Background(), "message_channel_bind:dingding:"+strconv.Itoa(product.ProductId))
	}

	for _, product := range req {
		glob.GRedis.LPush(context.Background(), "message_channel_bind:dingding:"+strconv.Itoa(product.ProductId),
			product.DingDingId)
	}
	return true

}
