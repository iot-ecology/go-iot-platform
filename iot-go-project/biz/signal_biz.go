package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type SignalBiz struct{}

var bizMqtt = MqttClientBiz{}

func (biz *SignalBiz) PageSignal(mqttClientId, ty string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var signals []models.Signal

	db := glob.GDb

	if mqttClientId != "" {
		db = db.Where("mqtt_client_id = ?", mqttClientId)
	}
	if ty != "" {
		db = db.Where("type = ?", ty)
	}
	db.Model(&models.Signal{}).Count(&pagination.Total) // 计算总记录数

	offset := (page - 1) * size
	db = db.Offset(offset).Limit(size).Find(&signals)

	for i, signal := range signals {
		id, err := bizMqtt.FindById(signal.MqttClientId)
		if err != nil {

		}
		if id != nil {
			signals[i].MqttClientName = id.ClientId
		}
	}

	pagination.Data = signals
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *SignalBiz) FindByIdForSignal(id int) (models.Signal, error) {
	var signal models.Signal

	result := glob.GDb.First(&signal, id)
	if result.Error != nil {

		return models.Signal{}, errors.New(result.Error.Error())

	}
	mqttClient, err := bizMqtt.FindById(signal.MqttClientId)
	if err != nil {
		return models.Signal{}, err
	}

	signal.MqttClientName = mqttClient.ClientId
	return signal, err

}

// FindByName 根据名称查询信号对象
func (biz *SignalBiz) FindByName(name string) (*models.Signal, error) {
	var signal models.Signal
	result := glob.GDb.First(&signal, "name = ?", name).Error
	if result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result
	}

	return &signal, nil
}

func (biz *SignalBiz) PageSignalWaringConfig(signalId int, mqttClientId string, page, size int) (*servlet.PaginationQ, error) {

	var pagination servlet.PaginationQ
	var signals []models.SignalWaringConfig

	db := glob.GDb

	if signalId != -1 {
		db = db.Where("signal_id = ?", signalId)
	}
	if mqttClientId != "" {
		db = db.Where("mqtt_client_id = ?", mqttClientId)

	}
	db.Model(&models.SignalWaringConfig{}).Count(&pagination.Total) // 计算总记录数

	offset := (page - 1) * size
	db = db.Offset(offset).Limit(size).Find(&signals)

	//for i, signal := range signals {
	//	forSignal, err := biz.FindByIdForSignal(signal.SignalId)
	//	if err != nil {
	//
	//	}
	//	signals[i].Signal = forSignal
	//}
	pagination.Data = signals
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func (biz *SignalBiz) SetSignalWaringCache(config models.SignalWaringConfig) {

	configBytes, _ := json.Marshal(config)
	glob.GRedis.LPush(context.Background(), fmt.Sprintf("waring:%d", config.SignalId), configBytes)
}
func (biz *SignalBiz) RemoveSignalWaringCache(config models.SignalWaringConfig) {
	configBytes, _ := json.Marshal(config)
	glob.GRedis.LRem(context.Background(), fmt.Sprintf("waring:%d", config.SignalId), 0, configBytes)
}

func (biz *SignalBiz) SetSignalCache(config *models.Signal) {
	configBytes, _ := json.Marshal(config)

	glob.GRedis.LPush(context.Background(), "signal:"+strconv.Itoa(config.MqttClientId), configBytes)
}

func (biz *SignalBiz) RemoveSignalCache(config *models.Signal) {
	configBytes, _ := json.Marshal(config)

	glob.GRedis.LRem(context.Background(), "signal:"+strconv.Itoa(config.MqttClientId), 0, configBytes)
}
