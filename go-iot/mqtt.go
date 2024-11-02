package main

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"sync"
)

var clock sync.Mutex

// MqttConfig 定义了MQTT客户端配置的结构体
type MqttConfig struct {
	Broker   string `json:"broker"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SubTopic string `json:"sub_topic"`
	ClientId string `json:"client_id"`
}

// MQTTMessage 结构体
type MQTTMessage struct {
	MQTTClientID string `json:"mqtt_client_id"`
	Message      string `json:"message"`
}


func StopMqttClient(clientId string,config MqttConfig) {

	zap.S().Infof("StopMqttClient 开始, clientId = %v", clientId)


	marshal, _ := json.Marshal(config)

	globalRedisClient.HDel(context.Background(), "mqtt_config:use", clientId)
	globalRedisClient.SRem(context.Background(), "node_bind:"+globalConfig.NodeInfo.Name, 0, clientId)
	AddNoUseConfig(config, marshal)
}

func StopMqttClient2(clientId string) {
	zap.S().Errorf("StopMqttClient 开始, clientId = %v", clientId)



	globalRedisClient.HDel(context.Background(), "mqtt_config:no", clientId)
	globalRedisClient.HDel(context.Background(), "mqtt_config:use", clientId)
	globalRedisClient.SRem(context.Background(), "node_bind:"+globalConfig.NodeInfo.Name, 0, clientId)
}


func PushMqttMsg(clientId string, topic string, qos byte, retained bool, payload string) {
	//client := c[clientId]
	//if client != nil {
	//	(*client).Publish(topic, qos, retained, payload)
	//}
}

func CreateMqttClientMin(broker string, port int, username string, password string, subTopic string,
	clientId string) bool {



	config := MqttConfig{
		Broker:   broker,
		Port:     port,
		Username: username,
		Password: password,
		SubTopic: subTopic,
		ClientId: clientId,
	}
	client := NewMqttClient(clientId,config)
	err := client.Connect(broker, username, password, port)
	if err != nil {
		zap.S().Errorf("mqtt connect err = %v", err)
        return false
	}
	go client.Subscribe(subTopic)
	go client.HandlerMsg()
	//c[clientId] = &client.client

	return true

}

// CreateMqttClient 函数根据传入的MqttConfig配置创建一个MQTT客户端
// 参数：
//
//	config MqttConfig: MQTT客户端配置信息
//
// 返回值：
//
//	int64: 创建MQTT客户端后的总数，如果达到最大客户端数量则返回-1，如果MQTT客户端配置异常则返回-2
func CreateMqttClient(config MqttConfig) int64 {

	i := globalRedisClient.SCard(context.Background(), "node_bind:"+globalConfig.NodeInfo.Name).Val()

	if globalConfig.NodeInfo.Size > i {
		clientMin := CreateMqttClientMin(config.Broker, config.Port, config.Username, config.Password, config.SubTopic, config.ClientId)
		if !clientMin  {
			return -2

		}
		zap.S().Debugf("创建mqtt客户端成功")
		return i + 1

	} else {
		return -1

	}

}
