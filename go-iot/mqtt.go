package main

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"time"
)

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
	PushToQueue("pre_handler", jsonData)

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	zap.S().Debugf("MQTT客户端链接成功")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {

	zap.S().Errorf("失去链接: %+v", err)
	reader := client.OptionsReader()
	id := reader.ClientID()
	StopMqttClient(id)
	config := configMap[id]

	jsonData, err := json.Marshal(config)
	if err != nil {
		zap.S().Errorf("to json error ,%+v", err)
	}
	PubCreateMqttClientOp(string(jsonData))
}

var c map[string]mqtt.Client

func StopMqttClient(clientId string) {
	client := c[clientId]
	if client != nil {
		client.Disconnect(0)

		// 删除在 no use 中的配置
		globalRedisClient.HDel(context.Background(), "mqtt_config:no", clientId)
		globalRedisClient.HDel(context.Background(), "mqtt_config:use", clientId)
		globalRedisClient.SRem(context.Background(), "node_bind:"+globalConfig.NodeInfo.Name, 0, clientId)
		delete(c, clientId)
	}

}

var configMap map[string]MqttConfig

func PushMqttMsg(clientId string, topic string, qos byte, retained bool, payload string) {
	client := c[clientId]
	client.Publish(topic, qos, retained, payload)
}

func CreateMqttClientMin(broker string, port int, username string, password string, subTopic string, clientId string) mqtt.Client {
	if configMap == nil {
		configMap = make(map[string]MqttConfig)
	}
	configMap[clientId] = MqttConfig{
		Broker:   broker,
		Port:     port,
		Username: username,
		Password: password,
		SubTopic: subTopic,
		ClientId: clientId,
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetPingTimeout(10 * time.Second)
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		zap.S().Error("创建MQTT客户端异常", token.Error())
		return nil

	}
	sub(client, subTopic)

	if c == nil {
		c = make(map[string]mqtt.Client)
	}
	c[clientId] = client

	return client

}

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	zap.S().Debugf("订阅主题: %s", topic)
}

// CreateMqttClient 创建MQTT客户端
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

	return -1

}
