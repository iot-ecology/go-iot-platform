package biz

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"log"
)

type MqttClientBiz struct{}

func (s *MqttClientBiz) CreateMqtt(client models.MqttClient) *models.MqttClient {

	id, err2 := s.FindByClientId(client.ClientId)
	if id != nil {
		panic("client_id 已存在")
	}

	if err2 != nil {
		panic(err2)
	}

	// 查询
	var result *models.MqttClient
	err := glob.GDb.First(&result, "host = ? AND port = ? AND client_id = ?", client.Host, client.Port, client.ClientId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		glob.GDb.Model(&models.MqttClient{}).Create(&client)
		return &client
	} else if err != nil {
		log.Fatalf("Error occurred during query: %v", err)
		panic(err)
	} else {
		panic("此数据已存在")

	}

	return nil
}

func (s *MqttClientBiz) Start(id any) *models.MqttClient {
	var mqttClient models.MqttClient
	glob.GDb.First(&mqttClient, id)
	s.SetScriptRedis(mqttClient.ClientId, mqttClient.Script)
	return &mqttClient
}

// FindByClientId 查找并返回与特定ClientID关联的MqttClient对象
func (s *MqttClientBiz) FindByClientId(clientId string) (*models.MqttClient, error) {
	var mqttClient models.MqttClient
	result := glob.GDb.First(&mqttClient, "client_id = ?", clientId).Error
	if result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result
	}

	return &mqttClient, nil
}

func (s *MqttClientBiz) FindById(id any) (*models.MqttClient, error) {
	var mqttClient models.MqttClient

	result := glob.GDb.First(&mqttClient, "id = ?", id).Error
	if result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result
	}

	return &mqttClient, nil
}

func (s *MqttClientBiz) SetScriptRedis(mqttClientId string, script string) {
	glob.GRedis.HSet(context.Background(), "mqtt_script", mqttClientId, script)
}

func (s *MqttClientBiz) PageMqttData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var mqttClients []models.MqttClient

	db := glob.GDb

	if name != "" {
		db = db.Where("client_id like ?", "%"+name+"%")
	}

	db.Model(&models.MqttClient{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db = db.Offset(offset).Limit(size).Find(&mqttClients)

	background := context.Background()
	var cc []string

	for _, el := range glob.GRedis.Keys(background, "node_bind:*").Val() {
		val := glob.GRedis.SMembers(background, el).Val()
		cc = append(cc, val...)
	}
	for i, client := range mqttClients {
		mqttClients[i].Start = contains(cc, client.ClientId)
	}

	pagination.Data = mqttClients
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

func contains(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
