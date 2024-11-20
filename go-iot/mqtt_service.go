package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"sync"
	"time"
)

// MqttInterface 定义了MQTT客户端的基本接口
type MqttInterface struct {
	client mqtt.Client
	Id     string
	Chan   chan []byte
	Config MqttConfig
	wg     sync.WaitGroup
}

// NewMqttClient 初始化并返回一个新的MqttInterface实例
func NewMqttClient(id string, config MqttConfig) *MqttInterface {
	return &MqttInterface{
		Id:     id,
		Chan:   make(chan []byte, 1000),
		Config: config,
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
	//opts.SetDefaultPublishHandler(m.messageHandler)
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		zap.S().Errorf("mqtt connection lost id = %s , error = %+v", m.Id, err)
		StopMqttClient(m.Id, m.Config)
	}

	opts.SetOrderMatters(false)
	opts.SetKeepAlive(60 * time.Second)
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
func (m *MqttInterface) Subscribe(topics string) error {
	var token = m.client.Subscribe(topics, 0, func(client mqtt.Client, msg mqtt.Message) {
		go func() {
			mqttMsg := MQTTMessage{
				MQTTClientID: m.Id,
				Message:      string(msg.Payload()),
			}
			jsonData, _ := json.Marshal(mqttMsg)

			m.Chan <- jsonData
		}()

	})

	if token.Wait() && token.Error() != nil {
		zap.S().Errorf(token.Error().Error())
		return token.Error()
	}
	return nil
}

// Publish 向一个主题发布消息
func (m *MqttInterface) Publish(topic string ,qos byte , reatined bool, payload interface{}) {
	token := m.client.Publish(topic, qos, reatined, payload)
	token.Wait()
}

// Disconnect 断开与MQTT服务器的连接
func (m *MqttInterface) Disconnect() {
	m.client.Disconnect(250)
}

func (m *MqttInterface) HandlerMsg() {
	for {
		c := <-m.Chan
		PushToQueue("pre_handler", c)

	}
}