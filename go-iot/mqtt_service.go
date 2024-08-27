package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

// MqttInterface 定义了MQTT客户端的基本接口
type MqttInterface struct {
	client mqtt.Client
	Id     string
}

// NewMqttClient 初始化并返回一个新的MqttInterface实例
func NewMqttClient(id string) *MqttInterface {
	return &MqttInterface{
		Id: id,
	}
}

// Connect 连接到MQTT服务器
func (m *MqttInterface) Connect(host, username, password string, port int) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	opts.SetUsername(username)
	opts.SetAutoReconnect(false)
	opts.SetPassword(password)
	opts.SetClientID(m.Id)
	opts.SetDefaultPublishHandler(m.messageHandler)
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		zap.S().Error("mqtt connection lost", zap.Error(err))
		StopMqttClient(m.Id)
	}
	opts.SetAutoReconnect(false)
	// 创建并启动客户端
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	m.client = client
	return nil
}

// messageHandler 处理接收到的消息
func (m *MqttInterface) messageHandler(client mqtt.Client, msg mqtt.Message) {
	mqttMsg := MQTTMessage{
		MQTTClientID: m.Id,
		Message:      string(msg.Payload()),
	}
	jsonData, err := json.Marshal(mqttMsg)
	if err != nil {
		zap.S().Errorf("Error marshalling MQTT message to JSON: %v", err)
	}
	go PushToQueue("pre_handler", jsonData)
}

// Subscribe 订阅一个或多个主题
func (m *MqttInterface) Subscribe(topics string) {
	if token := m.client.Subscribe(topics, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

// Publish 向一个主题发布消息
func (m *MqttInterface) Publish(topic string, payload interface{}) {
	token := m.client.Publish(topic, 0, false, payload)
	token.Wait()
}

// Disconnect 断开与MQTT服务器的连接
func (m *MqttInterface) Disconnect() {
	m.client.Disconnect(250)
}
