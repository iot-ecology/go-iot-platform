package notice

import (
	"go.uber.org/zap"
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type FeiShuBiz struct{}

func (biz *FeiShuBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var FeiShu []models.FeiShu

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.FeiShu{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&FeiShu)

	pagination.Data = FeiShu
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *FeiShuBiz) Bind(req []models.FeiShuBindProduct) bool {
	param := req[0]
	db := glob.GDb

	tx := db.Begin()
	if tx.Error != nil {
		return false

	}

	result := tx.Where("product_id = ?", param.ProductId).Delete(models.FeiShuBindProduct{})
	if result.Error != nil {
		// 如果出现错误，回滚事务
		tx.Rollback()
		return false

	}
	result = tx.Model(&models.FeiShuBindProduct{}).CreateInBatches(req, len(req))

	if result.Error != nil {
		tx.Rollback()
		zap.S().Infoln("Error occurred during creation:", result.Error)
		return false
	}
	if err := tx.Commit().Error; err != nil {
		return false

	}
	return true

}
