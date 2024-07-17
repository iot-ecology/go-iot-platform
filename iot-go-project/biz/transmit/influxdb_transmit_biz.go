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

type InfluxdbTransmitBiz struct{}

func (biz *InfluxdbTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var InfluxdbTransmits []models.InfluxdbTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.InfluxdbTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&InfluxdbTransmits)

	pagination.Data = InfluxdbTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *InfluxdbTransmitBiz) Bind(req models.InfluxdbTransmitBind) {
	glob.GDb.Model(models.InfluxdbTransmitBind{}).Create(req)
	jsonData := biz.toByte(req)
	// 缓存构造
	glob.GRedis.LPush(context.Background(), "transmit:influxdb:"+strconv.Itoa(req.MqttClientId), jsonData)
}

func (biz *InfluxdbTransmitBiz) toByte(req models.InfluxdbTransmitBind) []byte {
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
func (biz *InfluxdbTransmitBiz) ChangeEnable(req models.InfluxdbTransmitBind) {
	glob.GRedis.LRem(context.Background(), "transmit:influxdb:"+strconv.Itoa(req.MqttClientId), 1, biz.toByte(req))
}

var InfluxdbOp = influxdb2.InfluxDbOp{}

// MockScript 模拟执行脚本
func (biz *InfluxdbTransmitBiz) MockScript(dataRowList []common.DataRowList, script string) []common.DataRowList {
	return InfluxdbOp.RunScript(dataRowList, script)

}
