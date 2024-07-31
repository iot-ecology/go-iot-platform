package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type ClickhouseTransmitBiz struct{}

func (biz *ClickhouseTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var clickhouseTransmit []models.ClickhouseTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.ClickhouseTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&clickhouseTransmit)

	pagination.Data = clickhouseTransmit
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *ClickhouseTransmitBiz) SetRedis(param models.ClickhouseTransmit) {
	jsonData, err := json.Marshal(param)

	if err != nil {
		panic(err)
	}
	glob.GRedis.HSet(context.Background(), "transmit:clickhouse:"+strconv.Itoa(int(param.ID)), jsonData)
}

func (biz *ClickhouseTransmitBiz) DeleteRedis(param models.ClickhouseTransmit) {
	glob.GRedis.HDel(context.Background(), "transmit:clickhouse:"+strconv.Itoa(int(param.ID)))
}
