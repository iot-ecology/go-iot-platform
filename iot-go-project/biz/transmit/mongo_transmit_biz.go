package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type MongoTransmitBiz struct{}

func (biz *MongoTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var mongoTransmits []models.MongoTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.MongoTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&mongoTransmits)

	pagination.Data = mongoTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *MongoTransmitBiz) SetRedis(param models.MongoTransmit) {
	jsonData, err := json.Marshal(param)

	if err != nil {
		panic(err)
	}
	glob.GRedis.HSet(context.Background(), "transmit:mongo:"+strconv.Itoa(int(param.ID)), jsonData)
}

func (biz *MongoTransmitBiz) DeleteRedis(param models.MongoTransmit) {
	glob.GRedis.HDel(context.Background(), "transmit:mongo:"+strconv.Itoa(int(param.ID)))
}
