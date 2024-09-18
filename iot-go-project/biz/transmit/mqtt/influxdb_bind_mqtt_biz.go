package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"iot-transmit/cache"
	"iot-transmit/common"
	"iot-transmit/influxdb2"
	"strconv"
)

type InfluxdbTransmitBindBiz struct{}

func (biz *InfluxdbTransmitBindBiz) PageData(org string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.InfluxdbTransmitBind

	db := glob.GDb

	if org != "" {
		db = db.Where("org like ?", "%"+org+"%")
	}

	db.Model(&models.InfluxdbTransmitBind{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}


func (biz *InfluxdbTransmitBindBiz) Bind(req models.InfluxdbTransmitBind) {
	if req.Enable {
		jsonData := biz.toByte(req)
		// 缓存构造
		glob.GRedis.LPush(context.Background(), "transmit:influxdb:"+ req.Protocol +":"+string(req.DeviceUid) +":" +req.
			IdentificationCode, jsonData)
	}
}

func (biz *InfluxdbTransmitBindBiz) toByte(req models.InfluxdbTransmitBind) []byte {
	var ref models.InfluxdbTransmit

	glob.GDb.First(&ref, req.InfluxdbTransmitId)

	v := cache.InfluxTransmitCache{
		ID:          "influxdb-" + strconv.Itoa(int(req.ID)),
		Host:        ref.Host,
		Port:        ref.Port,
		Token:       ref.Token,
		Bucket:      req.Bucket,
		Org:         req.Org,
		Measurement: req.Measurement,
		Script:      req.Script,
	}
	jsonData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// ChangeEnable 修改启用状态
func (biz *InfluxdbTransmitBindBiz) ChangeEnable(req models.InfluxdbTransmitBind) {
	glob.GDb.Model(models.InfluxdbTransmitBind{}).Where("id = ?", req.ID).Update("enable", req.Enable)
	biz.HandlerRedis(req)
}

func (biz *InfluxdbTransmitBindBiz) HandlerRedis(req models.InfluxdbTransmitBind) {
	if req.Enable {
		biz.Bind(req)
	} else {
		glob.GRedis.LRem(context.Background(), "transmit:influxdb:"+req.Protocol +":"+string(req.DeviceUid) +":" +req.IdentificationCode, 1, biz.toByte(req))

	}
}

var InfluxdbOp = influxdb2.InfluxDbOp{}

// MockScript 模拟执行脚本
func (biz *InfluxdbTransmitBindBiz) MockScript(dataRowList []common.DataRowList, script string) []common.DataRowList {
	return InfluxdbOp.RunScript(dataRowList, script)

}