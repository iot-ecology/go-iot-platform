package biz

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"igp/ut"
	"strconv"
)

type DeviceInfoBiz struct{}

var productBiz = ProductBiz{}

func (biz *DeviceInfoBiz) PageData(sn, protocol string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.DeviceInfo

	db := glob.GDb

	if sn != "" {
		db = db.Where("sn like ?", "%"+sn+"%")
	}
	if protocol !=""{
		db = db.Where("sn = ?", protocol)

	}

	db.Model(&models.DeviceInfo{}).Count(&pagination.Total)
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	var resp []servlet.DeviceInfoRes

	for _, info := range dt {
		ProductName := productBiz.FindById(info.ProductId).Name
		resp = append(resp, servlet.DeviceInfoRes{
			ProductId:         info.ProductId,
			SN:                info.SN,
			ManufacturingDate: ut.LocalTime(info.ManufacturingDate),
			ProcurementDate:   ut.LocalTime(info.ProcurementDate),
			Source:            info.Source,
			WarrantyExpiry:    ut.LocalTime(info.WarrantyExpiry),
			PushInterval:      info.PushInterval,
			ErrorRate:         info.ErrorRate,
			Model:             info.Model,
			ProductName:       ProductName,
			Protocol: info.Protocol,
		})

	}
	pagination.Data = resp
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

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
func (biz *DeviceInfoBiz) FindBySn(sn string) *models.DeviceInfo {

	var dt models.DeviceInfo
	db := glob.GDb
	db.Where("sn = ?", sn).Find(&dt)
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
	background := context.Background()
	glob.GRedis.HSet(background, "struct:device_info", strconv.Itoa(int(newV.ID)), jsonData)

	glob.GRedis.Set(background,"share:device_info:" + newV.Protocol +":"+strconv.Itoa(int(newV.
		DeviceUid)) +":"+newV.IdentificationCode,
		jsonData,0)
}
func (biz *DeviceInfoBiz) RemoveRedis(id uint) {
	glob.GRedis.HDel(context.Background(), "struct:device_info", strconv.Itoa(int(id)))
}
