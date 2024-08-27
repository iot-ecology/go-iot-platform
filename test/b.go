package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
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
	reader := (*msg.client).OptionsReader()
	id := reader.ClientID()
	zap.S().Infof(id)
}

var choke = make(chan Cag,1000)

func createMqttClient(i int) mqtt.Client {
	mqtt.ERROR = log.New(getWriteSync(), "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(getWriteSync(), "[CRIT] ", 0)
	//mqtt.WARN = log.New(getWriteSync(), "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(getWriteSync(), "[DEBUG] ", 0)
	var broker = "172.17.0.1"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	s := uuid.New().String()
	opts.SetClientID(s)
	opts.SetUsername("admin")
	opts.SetPassword("public")
	opts.SetOrderMatters(false)
	opts.SetPingTimeout(10*time.Second)
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetAutoReconnect(false)
	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {
		zap.S().Errorf("Reconnecting to %s:%d", broker, port)

	})

	//opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
	//	choke <- Cag{
	//		client: &client,
	//		msg:    msg.Topic(),
	//	}
	//})

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	subscribe := client.Subscribe("/test_topic/"+strconv.Itoa(i), 0, func(client mqtt.Client, message mqtt.Message) {
			choke <- Cag{
				client: &client,
				msg:    message.Topic(),
			}
	})
	if subscribe.Wait() && subscribe.Error() != nil {
		panic(subscribe.Error())

	}

	return client
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