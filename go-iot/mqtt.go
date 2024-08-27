package main

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"log"
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

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	zap.S().Debugf("处理消息: %s  消息主题: %s\n", msg.Payload(), msg.Topic())

	reader := client.OptionsReader()
	id := reader.ClientID()

	// 创建 MQTTMessage 实例并序列化为 JSON
	mqttMsg := MQTTMessage{
		MQTTClientID: id,
		Message:      string(msg.Payload()),
	}
	jsonData, err := json.Marshal(mqttMsg)
	if err != nil {
		zap.S().Errorf("Error marshalling MQTT message to JSON: %v", err)
		return
	}
	//PushToQueue("pre_handler", jsonData)
	msgchan <- jsonData
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	zap.S().Debugf("MQTT客户端链接成功")
}

func handleMessage() {
	for {
		select {
		case msg := <-msgchan:
			// 处理消息
			PushToQueue("pre_handler", msg)
		}
	}
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {

	reader := client.OptionsReader()
	id := reader.ClientID()
	zap.S().Errorf("失去链接，id: %s ,error %+v：", id, err)
	StopMqttClient(id)
}

var c = make(map[string]*mqtt.Client)

func StopMqttClient(clientId string) {
	clock.Lock()
	defer clock.Unlock()
	zap.S().Errorf("StopMqttClient 开始, clientId = %v", clientId)
	delete(c, clientId)
	client := c[clientId]
	if client != nil {
		(*client).Disconnect(250)
	}
	config := configMap[clientId]
	marshal, _ := json.Marshal(config)

	globalRedisClient.HDel(context.Background(), "mqtt_config:use", clientId)
	globalRedisClient.SRem(context.Background(), "node_bind:"+globalConfig.NodeInfo.Name, 0, clientId)
	AddNoUseConfig(config, marshal)
}

var configMap map[string]MqttConfig
var msgchan = make(chan []byte, 1000)

func PushMqttMsg(clientId string, topic string, qos byte, retained bool, payload string) {
	client := c[clientId]
	if client != nil {
		(*client).Publish(topic, qos, retained, payload)
	}
}

func CreateMqttClientMin(broker string, port int, username string, password string, subTopic string, clientId string) mqtt.Client {
	clock.Lock()
	defer clock.Unlock()
	if configMap == nil {
		configMap = make(map[string]MqttConfig)
	}
	// 先判断 configMap 中是否有 clientId ， 如果有删除
	if _, ok := configMap[clientId]; ok {
		delete(configMap, clientId)
	}

	configMap[clientId] = MqttConfig{
		Broker:   broker,
		Port:     port,
		Username: username,
		Password: password,
		SubTopic: subTopic,
		ClientId: clientId,
	}
	mqtt.ERROR = log.New(getWriteSync(), "[ERROR] ", 0)
	//mqtt.CRITICAL = log.New(getWriteSync(), "[CRIT] ", 0)
	//mqtt.WARN = log.New(getWriteSync(), "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(getWriteSync(), "[DEBUG] ", 0)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetAutoReconnect(false)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		zap.S().Error("创建MQTT客户端异常", token.Error())
		return nil

	}
	sub(client, subTopic)

	c[clientId] = &client

	return client

}

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		reader := client.OptionsReader()
		id := reader.ClientID()

		// 创建 MQTTMessage 实例并序列化为 JSON
		mqttMsg := MQTTMessage{
			MQTTClientID: id,
			Message:      string(msg.Payload()),
		}
		jsonData, err := json.Marshal(mqttMsg)
		if err != nil {
			zap.S().Errorf("Error marshalling MQTT message to JSON: %v", err)
			return
		}
		PushToQueue("pre_handler", jsonData)
	})
	token.Wait()
	if token.Wait() && token.Error() != nil {
		zap.S().Error("订阅异常", token.Error())
	}
	zap.S().Debugf("订阅主题: %s", topic)
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
		if clientMin == nil {
			return -2

		}
		zap.S().Debugf("创建mqtt客户端成功")
		return i + 1

	} else {
		return -1

	}

}
