package biz

import (
	"context"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type TcpHandlerBiz struct{}

func (biz *TcpHandlerBiz) ById(id uint) (*models.TcpHandler, error) {
	var TcpHandler models.TcpHandler

	result := glob.GDb.First(&TcpHandler, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &TcpHandler, nil
}

func (biz *TcpHandlerBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var TcpHandlerList []models.TcpHandler

	db := glob.GDb
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	db.Model(&models.TcpHandler{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&TcpHandlerList)
	pagination.Data = TcpHandlerList
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *TcpHandlerBiz) SetRedis(data models.TcpHandler) {
	glob.GRedis.HSet(context.Background(), "struct:tcp", strconv.Itoa(int(data.DeviceInfoId)), data.Script)
}
func (biz *TcpHandlerBiz) RemoveRedis(data models.TcpHandler) {
	glob.GRedis.HDel(context.Background(), "struct:tcp", strconv.Itoa(int(data.DeviceInfoId)))
}
