package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"iot-transmit/cache"
	"iot-transmit/common"
	"iot-transmit/mongo"
	"strconv"
)

type MongoTransmitBiz struct{}

func (biz *MongoTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var MongoTransmits []models.MongoTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.MongoTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&MongoTransmits)

	pagination.Data = MongoTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}
func (biz *MongoTransmitBiz) Bind(req models.MongoTransmitBind) {
	glob.GDb.Model(models.MongoTransmitBind{}).Create(req)
	jsonData := biz.toByte(req)
	// 缓存构造
	glob.GRedis.LPush(context.Background(), "transmit:mongo:"+strconv.Itoa(req.MqttClientId), jsonData)
}

func (biz *MongoTransmitBiz) toByte(req models.MongoTransmitBind) []byte {
	var ref models.MongoTransmit

	glob.GDb.First(&ref, req.MongoTransmitId)

	v := cache.MongoTransmitCache{
		ID: "mongo-" + strconv.Itoa(int(req.ID)),

		Host:       ref.Host,
		Port:       ref.Port,
		Username:   ref.Username,
		Password:   ref.Password,
		Database:   req.Database,
		Collection: req.Collection,
		Script:     req.Script,
	}
	jsonData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// ChangeEnable 修改启用状态
func (biz *MongoTransmitBiz) ChangeEnable(req models.MongoTransmitBind) {
	glob.GRedis.LRem(context.Background(), "transmit:mongo:"+strconv.Itoa(req.MqttClientId), 1, biz.toByte(req))
}

var mongoOp = mongo.MongoOp{}

// MockScript 模拟执行脚本
func (biz *MongoTransmitBiz) MockScript(dataRowList []common.DataRowList, script string) []map[string]interface{} {
	return mongoOp.RunScript(dataRowList, script)
}
