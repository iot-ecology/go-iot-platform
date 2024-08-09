package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"igp/ut"
	"strconv"
)

type SignalBiz struct{}

var bizMqtt = MqttClientBiz{}

func (biz *SignalBiz) PageSignal(deviceUid ,protocol ,ty string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var signals []models.Signal

	db := glob.GDb

	if ty!=""{
		db = db.Where("type = ?", ty)
	}
	if deviceUid != "" {
		db = db.Where("device_uid = ?", deviceUid)
	}
	if protocol !=""{
		db = db.Where("protocol = ?", protocol)
	}
	db.Model(&models.Signal{}).Count(&pagination.Total) // 计算总记录数

	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&signals)

	for i, signal := range signals {
		id, err := bizMqtt.FindById(signal.DeviceUid)
		if err != nil {
			zap.S().Errorf("error %+v", err)
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
	mqttClient, err := bizMqtt.FindById(signal.DeviceUid)
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

func (biz *SignalBiz) PageSignalWaringConfig(signalId int, device_uid, code, protocol string, page, size int) (*servlet.PaginationQ, error) {

	var pagination servlet.PaginationQ
	var dt []models.SignalWaringConfig

	db := glob.GDb

	if signalId != -1 {
		db = db.Where("signal_id = ?", signalId)
	}
	if device_uid != "" {
		db = db.Where("device_uid = ?", device_uid)

	}
	if code != "" {
		db = db.Where("identification_code = ?", code)

	}
	if protocol != "" {
		db = db.Where("protocol = ?", protocol)

	}

	db.Model(&models.SignalWaringConfig{}).Count(&pagination.Total) // 计算总记录数

	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)


	pagination.Data = dt
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

	glob.GRedis.LPush(context.Background(), "signal:"+strconv.Itoa(config.DeviceUid) +":"+  config.
		IdentificationCode, configBytes)
}

func (biz *SignalBiz) RemoveSignalCache(config *models.Signal) {
	configBytes, _ := json.Marshal(config)

	glob.GRedis.LRem(context.Background(), "signal:"+strconv.Itoa(config.DeviceUid) +":" + config.
		IdentificationCode, 0, configBytes)
}



func (b SignalBiz) InitMongoCollection(m *models.SignalWaringConfig) {
	name := ut.CalcCollectionName(glob.GConfig.MongoConfig.WaringCollection, m.ID)
	ut.CheckCollectionAndCreate(glob.GConfig.MongoConfig.WaringCollection, name)
}