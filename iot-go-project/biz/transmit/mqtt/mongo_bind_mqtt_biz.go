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

type MongoTransmitBindBiz struct{}

func (biz *MongoTransmitBindBiz) PageData(collection string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.MongoTransmitBind

	db := glob.GDb

	if collection != "" {
		db = db.Where("collection like ?", "%"+collection+"%")
	}

	db.Model(&models.MongoTransmitBind{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}


func (biz *MongoTransmitBindBiz) Bind(req models.MongoTransmitBind) {
	if req.Enable {
		jsonData := biz.toByte(req)
		// 缓存构造
		glob.GRedis.LPush(context.Background(), "transmit:mongo:"+req.Protocol +":"+strconv.Itoa(req.DeviceUid) +":" +req.IdentificationCode, jsonData)
	}
}

func (biz *MongoTransmitBindBiz) toByte(req models.MongoTransmitBind) []byte {
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
func (biz *MongoTransmitBindBiz) ChangeEnable(req models.MongoTransmitBind) {
	glob.GDb.Model(models.MongoTransmitBind{}).Where("id = ?", req.ID).Update("enable", req.Enable)
	biz.HandlerRedis(req)
}

func (biz *MongoTransmitBindBiz) HandlerRedis(req models.MongoTransmitBind) {
	if req.Enable {
		biz.Bind(req)
	} else {
		glob.GRedis.LRem(context.Background(), "transmit:mongo:"+req.Protocol +":"+strconv.Itoa(req.DeviceUid) +":" +req.IdentificationCode, 1, biz.toByte(req))
	}
}

var mongoOp = mongo.MongoOp{}

// MockScript 模拟执行脚本
func (biz *MongoTransmitBindBiz) MockScript(dataRowList []common.DataRowList, script string) []map[string]interface{} {
	return mongoOp.RunScript(dataRowList, script)
}
