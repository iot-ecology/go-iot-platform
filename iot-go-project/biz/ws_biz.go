package biz

import (
	"context"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type WebsocketHandlerBiz struct{}

func (biz *WebsocketHandlerBiz) ById(id uint) (*models.WebsocketHandler, error) {
	var WebsocketHandler models.WebsocketHandler

	result := glob.GDb.First(&WebsocketHandler, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &WebsocketHandler, nil
}

func (biz *WebsocketHandlerBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var WebsocketHandlerList []models.WebsocketHandler

	db := glob.GDb
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	db.Model(&models.WebsocketHandler{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&WebsocketHandlerList)
	pagination.Data = WebsocketHandlerList
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *WebsocketHandlerBiz) SetRedis(data models.WebsocketHandler) {
	glob.GRedis.HSet(context.Background(), "struct:Websocket", strconv.Itoa(int(data.ID)), data.Script)
}
func (biz *WebsocketHandlerBiz) RemoveRedis(data models.WebsocketHandler) {
	glob.GRedis.HDel(context.Background(), "struct:Websocket", strconv.Itoa(int(data.ID)))
}
