package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type InfluxdbTransmitBiz struct{}

func (biz *InfluxdbTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var influxdbTransmits []models.InfluxdbTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.InfluxdbTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&influxdbTransmits)

	pagination.Data = influxdbTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *InfluxdbTransmitBiz) SetRedis(param models.InfluxdbTransmit) {
	jsonData, err := json.Marshal(param)

	if err != nil {
		panic(err)
	}
	glob.GRedis.HSet(context.Background(), "transmit:influxdb:"+strconv.Itoa(int(param.ID)), jsonData)
}

func (biz *InfluxdbTransmitBiz) DeleteRedis(param models.InfluxdbTransmit) {
	glob.GRedis.HDel(context.Background(), "transmit:influxdb:"+strconv.Itoa(int(param.ID)))
}

