package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var ccc = make(map[int]mqtt.Client)

func main1() {
	InitLog()

	zap.S().Infof("启动")
	for i := range 100 {
		mqttlient := createMqttClient(i)
		ccc[i] = mqttlient
	}

	select {
	case msg := <-choke:
		go funcName(msg)
	}
	select {}
}

func funcName(msg Cag) {
	//reader := (*msg.client).OptionsReader()
	//id := reader.ClientID()
	//zap.S().Infof(id)
}

var choke = make(chan Cag,1000)

func createMqttClient(i int) mqtt.Client {
	s := uuid.New().String()
	client := NewMqttClient(s)
	client.Connect("s")
	client.Subscribe("/test_topic/"+strconv.Itoa(i))
	return client.client
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	zap.S().Infof("Connected")
}
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	reader := client.OptionsReader()
	id := reader.ClientID()
	zap.S().Debugf("id %+v", id)
	time.Sleep(1 * time.Second)
	//fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	zap.S().Errorf("Connect lost: %v", err)
}

type Cag struct {
	client *mqtt.Client
	msg string
}