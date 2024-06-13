package main

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"strings"
)

// AddNoUseConfig 将不使用的MQTT配置信息存储到Redis中
//
// 参数：
// config MqttConfig - MQTT配置信息
// body []byte - 待存储的字节数组
//
// 返回值：
// 无返回值
func AddNoUseConfig(config MqttConfig, body []byte) {
	globalRedisClient.HSet(context.Background(), "mqtt_config:no", config.ClientId, body)
}

// GetNoUseConfig 从Redis中获取不使用的MQTT配置信息列表
//
// 返回值：
// []string - 包含不使用MQTT配置信息的字符串切片
func GetNoUseConfig() []string {
	var noUseConfigs []string
	for _, s2 := range globalRedisClient.HGetAll(context.Background(), "mqtt_config:no").Val() {

		noUseConfigs = append(noUseConfigs, s2)
	}
	return noUseConfigs
}

// RemoveNoUseConfig 从Redis中删除不使用的MQTT配置信息
//
// 参数：
//
//	config MqttConfig - MQTT配置信息
//
// 返回值：
//
//	无
func RemoveNoUseConfig(config MqttConfig) {
	globalRedisClient.HDel(context.Background(), "mqtt_config:no", config.ClientId)
}

// AddUseConfig 将MQTT配置信息存储到Redis中
//
// 参数：
// config MqttConfig - MQTT配置信息
//
// 返回值：
// 无
func AddUseConfig(config MqttConfig) {

	jsonStr, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	globalRedisClient.HSet(context.Background(), "mqtt_config:use", config.ClientId, jsonStr)
}

// GetUseConfig 函数从Redis中获取已使用的MQTT配置信息
//
// 参数：
//
//	client_id string - MQTT客户端ID
//
// 返回值：
//
//	string - 对应的MQTT配置信息
func GetUseConfig(clientId string) string {
	return globalRedisClient.HGet(context.Background(), "mqtt_config:use", clientId).Val()
}

// GetNoUseConfigById 根据MQTT客户端ID获取不使用的MQTT配置信息
//
// 参数：
//
//	client_id string - MQTT客户端ID
//
// 返回值：
//
//	string - 对应的不使用的MQTT配置信息
func GetNoUseConfigById(clientId string) string {
	return globalRedisClient.HGet(context.Background(), "mqtt_config:no", clientId).Val()
}

// RemoveUseConfig 从Redis中删除已使用的MQTT配置信息
//
// 参数：
// config MqttConfig - 待删除的MQTT配置信息
//
// 返回值：
// 无
func RemoveUseConfig(config MqttConfig) {
	globalRedisClient.HDel(context.Background(), "mqtt_config:use", config.ClientId)
}

// CheckHasConfig 检查Redis中是否存在已使用的MQTT配置信息
//
// 参数：
//
//	config MqttConfig - 待检查的MQTT配置信息
//
// 返回值：
//
//	bool - 若Redis中存在该配置信息则返回true，否则返回false
func CheckHasConfig(config MqttConfig) bool {
	return globalRedisClient.HExists(context.Background(), "mqtt_config:use", config.ClientId).Val()
}

// BindNode 函数用于将MQTT客户端绑定到指定节点
//
// 参数：
//
//	config MqttConfig - MQTT配置信息
//	node_name string - 节点名称
//
// 返回值：
//
//	无
func BindNode(config MqttConfig, nodeName string) {
	globalRedisClient.SAdd(context.Background(), "node_bind:"+nodeName, config.ClientId)

	RemoveNoUseConfig(config)
	AddUseConfig(config)

}

// RemoveBindNode 函数用于从指定节点解绑MQTT客户端
//
// 参数：
//
//	client_id string - MQTT客户端ID
//	node_name string - 节点名称
//
// 返回值：
//
//	无
func RemoveBindNode(clientId string, nodeName string) {
	globalRedisClient.SRem(context.Background(), "node_bind:"+nodeName, 0, clientId)

	configStr := GetUseConfig(clientId)
	var config MqttConfig
	var byt = []byte(configStr)

	err := json.Unmarshal(byt, &config)
	if err != nil {
		zap.S().Error("Error unmarshalling JSON:", err)
		return
	}

	RemoveUseConfig(config)
	AddNoUseConfig(config, byt)
}

// GetBindClientId 函数用于获取指定节点下所有绑定的MQTT客户端ID列表
//
// 参数：
// node_name string - 节点名称
//
// 返回值：
// []string - 包含所有绑定在该节点下的MQTT客户端ID的字符串切片
func GetBindClientId(nodeName string) []string {
	return globalRedisClient.SMembers(context.Background(), "node_bind:"+nodeName).Val()

}

// FindMqttClientId 函数用于查找MQTT客户端ID所在的节点名称
//
// 参数：
//
//	mqtt_client_id string：待查找的MQTT客户端ID
//
// 返回值：
//
//	string：MQTT客户端ID所在的节点名称，若未找到则返回空字符串
func FindMqttClientId(mqttClientId string) string {
	background := context.Background()
	result, err := globalRedisClient.Keys(background, "node_bind:*").Result()
	if err != nil {

	}
	for _, elm := range result {
		for _, storageMqttClientId := range globalRedisClient.SMembers(background, elm).Val() {
			if storageMqttClientId == mqttClientId {
				modifiedString := strings.Replace(elm, "node_bind:", "", -1)
				return modifiedString

			}

		}
	}
	return ""
}
