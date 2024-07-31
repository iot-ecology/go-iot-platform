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

type KafkaTransmitBindBiz struct{}

func (biz *KafkaTransmitBindBiz) PageData(topic string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.KafkaTransmitBind

	db := glob.GDb

	if topic != "" {
		db = db.Where("topic like ?", "%"+topic+"%")
	}

	db.Model(&models.KafkaTransmitBind{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *KafkaTransmitBindBiz) Bind(req models.KafkaTransmitBind) {
	if req.Enable {
		jsonData := biz.toByte(req)
		// 缓存构造
		glob.GRedis.LPush(context.Background(), "transmit:Kafka:"+req.DeviceUid +":" +req.IdentificationCode, jsonData)
	}
}

func (biz *KafkaTransmitBindBiz) toByte(req models.KafkaTransmitBind) []byte {
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
func (biz *KafkaTransmitBindBiz) ChangeEnable(req models.KafkaTransmitBind) {
	glob.GDb.Model(models.KafkaTransmitBind{}).Where("id = ?", req.ID).Update("enable", req.Enable)
	biz.HandlerRedis(req)
}

func (biz *KafkaTransmitBindBiz) HandlerRedis(req models.KafkaTransmitBind) {
	if req.Enable {
		biz.Bind(req)
	} else {
		glob.GRedis.LRem(context.Background(), "transmit:Kafka:"+req.DeviceUid +":" +req.IdentificationCode, 1, biz.toByte(req))
	}
}
