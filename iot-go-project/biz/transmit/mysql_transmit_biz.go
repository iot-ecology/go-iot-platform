package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type MySQLTransmitBiz struct{}

func (biz *MySQLTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var mySQLTransmits []models.MySQLTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.MySQLTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&mySQLTransmits)

	pagination.Data = mySQLTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *MySQLTransmitBiz) SetRedis(param models.MySQLTransmit) {
	jsonData, err := json.Marshal(param)

	if err != nil {
		panic(err)
	}
	glob.GRedis.HSet(context.Background(), "transmit:mysql:"+strconv.Itoa(int(param.ID)), jsonData)
}

func (biz *MySQLTransmitBiz) DeleteRedis(param models.MySQLTransmit) {
	glob.GRedis.HDel(context.Background(), "transmit:mysql:"+strconv.Itoa(int(param.ID)))
}