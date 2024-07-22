package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"iot-transmit/cache"
	"strconv"
)

type KafkaTransmitBiz struct{}

func (biz *KafkaTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var KafkaTransmits []models.KafkaTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.KafkaTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&KafkaTransmits)

	pagination.Data = KafkaTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *KafkaTransmitBiz) Bind(req models.KafkaTransmitBind) {
	glob.GDb.Model(models.KafkaTransmitBind{}).Create(req)
	jsonData := biz.toByte(req)
	// 缓存构造
	glob.GRedis.LPush(context.Background(), "transmit:Kafka:"+strconv.Itoa(req.MqttClientId), jsonData)
}

func (biz *KafkaTransmitBiz) toByte(req models.KafkaTransmitBind) []byte {
	var KafkaInfo models.KafkaTransmit

	glob.GDb.First(&KafkaInfo, req.KafkaTransmitId)

	v := cache.KafkaTransmitCache{
		ID:     "Kafka-" + strconv.Itoa(int(req.ID)),
		Host:   KafkaInfo.Host,
		Port:   KafkaInfo.Port,
		Script: req.Script,
		Topic:  req.Topic,
	}
	jsonData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// ChangeEnable 修改启用状态
func (biz *KafkaTransmitBiz) ChangeEnable(req models.KafkaTransmitBind) {
	glob.GDb.Model(models.KafkaTransmitBind{}).Where("id = ?", req.ID).Update("enable", req.Enable)
	glob.GRedis.LRem(context.Background(), "transmit:Kafka:"+strconv.Itoa(req.MqttClientId), 1, biz.toByte(req))
}

//var KafkaOp = Kafka.KafkaOp{}
//
//// MockScript 模拟执行脚本
//func (biz *KafkaTransmitBiz) MockScript(dataRowList []common.DataRowList, script string) [][]Kafka.KafkaParam {
//	return KafkaOp.RunScript(dataRowList, script)
//}
