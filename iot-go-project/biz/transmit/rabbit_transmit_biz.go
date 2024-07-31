package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"iot-transmit/cache"
	"strconv"
)

type RabbitTransmitBiz struct{}

func (biz *RabbitTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var RabbitTransmits []models.RabbitmqTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.RabbitmqTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&RabbitTransmits)

	pagination.Data = RabbitTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *RabbitTransmitBiz) SetRedis(param models.RabbitmqTransmit) {
	jsonData, err := json.Marshal(param)

	if err != nil {
		panic(err)
	}
	glob.GRedis.HSet(context.Background(), "transmit:rabbit:"+strconv.Itoa(int(param.ID)), jsonData)
}

func (biz *RabbitTransmitBiz) DeleteRedis(param models.RabbitmqTransmit) {
	glob.GRedis.HDel(context.Background(), "transmit:rabbit:"+strconv.Itoa(int(param.ID)))
}