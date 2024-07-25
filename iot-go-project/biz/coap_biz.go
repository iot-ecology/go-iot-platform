package biz

import (
	"context"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type CoapHandlerBiz struct{}

func (biz *CoapHandlerBiz) ById(id uint) (*models.CoapHandler, error) {
	var CoapHandler models.CoapHandler

	result := glob.GDb.First(&CoapHandler, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &CoapHandler, nil
}

func (biz *CoapHandlerBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var CoapHandlerList []models.CoapHandler

	db := glob.GDb
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	db.Model(&models.CoapHandler{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&CoapHandlerList)
	pagination.Data = CoapHandlerList
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *CoapHandlerBiz) SetRedis(data models.CoapHandler) {
	glob.GRedis.HSet(context.Background(), "struct:Coap", strconv.Itoa(int(data.ID)), data.Script)
}
func (biz *CoapHandlerBiz) RemoveRedis(data models.CoapHandler) {
	glob.GRedis.HDel(context.Background(), "struct:Coap", strconv.Itoa(int(data.ID)))
}
