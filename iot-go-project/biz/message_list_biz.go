package biz

import (
	"igp/glob"
	"igp/models"
	"igp/servlet"
)

type MessageListBiz struct{}

func (biz *MessageListBiz) PageData(messageTypeId string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var MessageLists []models.MessageList

	db := glob.GDb

	if messageTypeId != "" {
		db = db.Where("message_type_id = ?", "%"+messageTypeId+"%")
	}

	db.Model(&models.MessageList{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&MessageLists)

	pagination.Data = MessageLists
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
