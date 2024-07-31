package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type CassandraTransmitBiz struct{}

func (biz *CassandraTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var cassandraTransmit []models.CassandraTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.CassandraTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&cassandraTransmit)

	pagination.Data = cassandraTransmit
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *CassandraTransmitBiz) SetRedis(param models.CassandraTransmit) {
	jsonData, err := json.Marshal(param)

	if err != nil {
		panic(err)
	}
	glob.GRedis.HSet(context.Background(), "transmit:cassandra:"+strconv.Itoa(int(param.ID)), jsonData)
}

func (biz *CassandraTransmitBiz) DeleteRedis(param models.CassandraTransmit) {
	glob.GRedis.HDel(context.Background(), "transmit:cassandra:"+strconv.Itoa(int(param.ID)))
}
