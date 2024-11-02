package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type KafkaTransmitBiz struct{}

func (biz *KafkaTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var KafkaTransmits []models.KafkaTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.KafkaTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&KafkaTransmits)

	pagination.Data = KafkaTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}


func (biz *KafkaTransmitBiz) SetRedis(param models.KafkaTransmit) {
	jsonData, err := json.Marshal(param)

	if err != nil {
		panic(err)
	}
	glob.GRedis.HSet(context.Background(), "transmit:kafka:"+strconv.Itoa(int(param.ID)), jsonData)
}

func (biz *KafkaTransmitBiz) DeleteRedis(param models.KafkaTransmit) {
	glob.GRedis.HDel(context.Background(), "transmit:kafka:"+strconv.Itoa(int(param.ID)))
}

