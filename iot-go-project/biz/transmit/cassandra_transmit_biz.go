package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"iot-transmit/cache"
	"iot-transmit/cassandra"
	"iot-transmit/common"
	"strconv"
)

type CassandraTransmitBiz struct{}

func (biz *CassandraTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var CassandraTransmit []models.CassandraTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.CassandraTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&CassandraTransmit)

	pagination.Data = CassandraTransmit
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *CassandraTransmitBiz) Bind(req models.CassandraTransmitBind) {
	glob.GDb.Model(models.CassandraTransmitBind{}).Create(req)
	jsonData := biz.toByte(req)
	// 缓存构造
	glob.GRedis.LPush(context.Background(), "transmit:cassandra:"+strconv.Itoa(req.MqttClientId), jsonData)
}

func (biz *CassandraTransmitBiz) toByte(req models.CassandraTransmitBind) []byte {
	var ref models.CassandraTransmit

	glob.GDb.First(&ref, req.CassandraTransmitId)

	v := cache.CassandraTransmitCache{
		ID:       "cassandra-" + strconv.Itoa(int(req.ID)),
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

// ChangeEnable 修改启用状态
func (biz *CassandraTransmitBiz) ChangeEnable(req models.CassandraTransmitBind) {
	glob.GDb.Model(models.CassandraTransmitBind{}).Where("id = ?", req.ID).Update("enable", req.Enable)

	glob.GRedis.LRem(context.Background(), "transmit:cassandra:"+strconv.Itoa(req.MqttClientId), 1, biz.toByte(req))
}

var CassandraOp = cassandra.CassandraOp{}

// MockScript 模拟执行脚本
func (biz *CassandraTransmitBiz) MockScript(dataRowList []common.DataRowList, script string) [][]cassandra.CassandraParam {

	return CassandraOp.RunScript(dataRowList, script)
}
