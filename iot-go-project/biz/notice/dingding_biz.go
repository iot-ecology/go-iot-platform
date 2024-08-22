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

func (biz *DingDingBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dingding []models.DingDing

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
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
