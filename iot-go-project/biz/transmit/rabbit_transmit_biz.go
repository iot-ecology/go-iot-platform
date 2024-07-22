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

type RabbitTransmitBiz struct{}

func (biz *RabbitTransmitBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var RabbitTransmits []models.RabbitmqTransmit

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.RabbitmqTransmit{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&RabbitTransmits)

	pagination.Data = RabbitTransmits
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *RabbitTransmitBiz) Bind(req models.RabbitmqTransmitBind) {
	glob.GDb.Model(models.RabbitmqTransmitBind{}).Create(req)
	jsonData := biz.toByte(req)
	// 缓存构造
	glob.GRedis.LPush(context.Background(), "transmit:Rabbit:"+strconv.Itoa(req.MqttClientId), jsonData)
}

func (biz *RabbitTransmitBiz) toByte(req models.RabbitmqTransmitBind) []byte {
	var RabbitInfo models.RabbitmqTransmit

	glob.GDb.First(&RabbitInfo, req.RabbitmqTransmitId)

	v := cache.RabbitTransmitCache{
		ID:     "Rabbit-" + strconv.Itoa(int(req.ID)),
		Host:   RabbitInfo.Host,
		Port:   RabbitInfo.Port,
		Script: req.Script,
		Exchange:  req.Exchange,
		RoutingKey:  req.RoutingKey,
	}
	jsonData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// ChangeEnable 修改启用状态
func (biz *RabbitTransmitBiz) ChangeEnable(req models.RabbitmqTransmitBind) {
	glob.GDb.Model(models.RabbitmqTransmitBind{}).Where("id = ?", req.ID).Update("enable", req.Enable)
	glob.GRedis.LRem(context.Background(), "transmit:Rabbit:"+strconv.Itoa(req.MqttClientId), 1, biz.toByte(req))
}

//var RabbitOp = Rabbit.RabbitOp{}
//
//// MockScript 模拟执行脚本
//func (biz *RabbitTransmitBiz) MockScript(dataRowList []common.DataRowList, script string) [][]Rabbit.RabbitParam {
//	return RabbitOp.RunScript(dataRowList, script)
//}
