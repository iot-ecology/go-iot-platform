package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"igp/servlet/transmit"
	"iot-transmit/clickhouse"
	"iot-transmit/common"
	"strconv"
)

type ClickhouseTransmitBiz struct{}

func (biz *ClickhouseTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var ClickhouseTransmit []models.ClickhouseTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.InfluxdbTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&ClickhouseTransmit)

	pagination.Data = ClickhouseTransmit
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *ClickhouseTransmitBiz) Bind(req models.ClickhouseTransmitBind) {
	glob.GDb.Model(models.ClickhouseTransmitBind{}).Create(req)
	jsonData := biz.toByte(req)
	// 缓存构造
	glob.GRedis.LPush(context.Background(), "transmit:clickhouse:"+strconv.Itoa(req.MqttClientId), jsonData)
}

var clickhouseOp = clickhouse.ClickhouseOp{}

func (biz *ClickhouseTransmitBiz) toByte(req models.ClickhouseTransmitBind) []byte {
	var ref models.ClickhouseTransmit

	glob.GDb.First(&ref, req.ClickhouseTransmitId)

	v := transmit.ClickhouseTransmitCache{
		ID:       req.ID,
		Host:     ref.Host,
		Port:     ref.Port,
		Username: ref.Username,
		Password: ref.Password,
		Database: req.Database,
		Table:    req.Table,
		Script:   req.Script,
	}
	jsonData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// changeEnable 修改启用状态
func (biz *ClickhouseTransmitBiz) ChangeEnable(req models.ClickhouseTransmitBind) {
	glob.GRedis.LRem(context.Background(), "transmit:clickhouse:"+strconv.Itoa(req.MqttClientId), 1, biz.toByte(req))
}

// mockScript 模拟执行脚本
func (biz *ClickhouseTransmitBiz) MockScript(dataRowList []common.DataRowList,
	script string) [][]clickhouse.ClickhouseParam {
	return clickhouseOp.RunScript(dataRowList, script)

}
