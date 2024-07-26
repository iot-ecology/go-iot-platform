package biz

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type DeviceInfoBiz struct{}

var productBiz = ProductBiz{}

func (biz *DeviceInfoBiz) PageData(sn string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.DeviceInfo

	db := glob.GDb

	if sn != "" {
		db = db.Where("sn like ?", "%"+sn+"%")
	}

	db.Model(&models.DeviceInfo{}).Count(&pagination.Total)
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	for i, info := range dt {
		dt[i].ProductName = productBiz.FindById(info.ProductId).Name


	}
	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
//mqttClients[i].LastPushTime = glob.GRedis.Get(context.Background(), "last_push_time:"+strconv.Itoa(int(client.ID))).Val()


func (biz *DeviceInfoBiz) FindById(id uint) *models.DeviceInfo {
	redis := biz.FindByIdWithRedis(id)
	if redis != nil {
		return redis
	}

	var dt models.DeviceInfo
	db := glob.GDb
	db.Where("id = ?", id).Find(&dt)
	biz.SetRedis(dt)
	return &dt
}

func (biz *DeviceInfoBiz) FindByIdWithRedis(id uint) *models.DeviceInfo {
	val := glob.GRedis.HGet(context.Background(), "struct:device_info", strconv.Itoa(int(id))).Val()

	var res models.DeviceInfo
	err := json.Unmarshal([]byte(val), &res)
	if err != nil {
		return nil
	}
	return &res
}

func (biz *DeviceInfoBiz) SetRedis(newV models.DeviceInfo) {
	jsonData, _ := json.Marshal(newV)
	glob.GRedis.HSet(context.Background(), "struct:device_info", strconv.Itoa(int(newV.ID)), jsonData)
}
func (biz *DeviceInfoBiz) RemoveRedis(id uint) {
	glob.GRedis.HDel(context.Background(), "struct:device_info", strconv.Itoa(int(id)))
}
