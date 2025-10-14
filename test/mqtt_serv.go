/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)




// MqttInterface 定义了MQTT客户端的基本接口
type MqttInterface struct {
	client mqtt.Client
	Id     string
}

// NewMqttClient 初始化并返回一个新的MqttInterface实例
func NewMqttClient( id string) *MqttInterface {
	return &MqttInterface{
		Id: id,
	}
}

// Connect 连接到MQTT服务器
func (m *MqttInterface) Connect(broker string) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://172.17.0.1:1883")
	opts.SetUsername("admin")
	opts.SetAutoReconnect(false)
	opts.SetPassword("public")
	opts.SetClientID(m.Id)
	opts.SetDefaultPublishHandler(m.messageHandler)

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
	//fmt.Printf("Received message on topic: %s\nMessage: %s\n", msg.Topic(), msg.Payload())
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