package biz

import (
	"context"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type HttpHandlerBiz struct{}

func (biz *HttpHandlerBiz) ById(id uint) (*models.HttpHandler, error) {
	var HttpHandler models.HttpHandler

	result := glob.GDb.First(&HttpHandler, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &HttpHandler, nil
}

func (biz *HttpHandlerBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var HttpHandlerList []models.HttpHandler

	db := glob.GDb
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	db.Model(&models.HttpHandler{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&HttpHandlerList)
	pagination.Data = HttpHandlerList
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *HttpHandlerBiz) SetRedis(data models.HttpHandler) {
	glob.GRedis.HSet(context.Background(), "struct:Http", strconv.Itoa(int(data.ID)), data.Script)
}
func (biz *HttpHandlerBiz) RemoveRedis(data models.HttpHandler) {
	glob.GRedis.HDel(context.Background(), "struct:Http", strconv.Itoa(int(data.ID)))
}
