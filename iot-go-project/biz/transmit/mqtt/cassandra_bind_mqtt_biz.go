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

type CassandraTransmitBindBiz struct{}

func (biz *CassandraTransmitBindBiz) PageData(table string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.CassandraTransmitBind

	db := glob.GDb

	if table != "" {
		db = db.Where("table like ?", "%"+table+"%")
	}

	db.Model(&models.CassandraTransmitBind{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *CassandraTransmitBindBiz) Bind(req models.CassandraTransmitBind) {
	if req.Enable == true {
		jsonData := biz.toByte(req)
		// 缓存构造
		glob.GRedis.LPush(context.Background(), "transmit:cassandra:"+req.DeviceUid +":" +req.IdentificationCode, jsonData)
	}
}

func (biz *CassandraTransmitBindBiz) toByte(req models.CassandraTransmitBind) []byte {
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
func (biz *CassandraTransmitBindBiz) ChangeEnable(req models.CassandraTransmitBind) {
	glob.GDb.Model(models.CassandraTransmitBind{}).Where("id = ?", req.ID).Update("enable", req.Enable)
	biz.HandlerRedis(req)
}

func (biz *CassandraTransmitBindBiz) HandlerRedis(req models.CassandraTransmitBind) {
	if req.Enable == false {
		glob.GRedis.LRem(context.Background(), "transmit:cassandra:"+req.DeviceUid +":" +req.IdentificationCode, 1, biz.toByte(req))
	} else {
		biz.Bind(req)
	}
}

var CassandraOp = cassandra.CassandraOp{}

// MockScript 模拟执行脚本
func (biz *CassandraTransmitBindBiz) MockScript(dataRowList []common.DataRowList, script string) [][]cassandra.CassandraParam {

	return CassandraOp.RunScript(dataRowList, script)
}
