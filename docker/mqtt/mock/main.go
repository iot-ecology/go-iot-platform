package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"strconv"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	var broker = "localhost"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("admin")
	opts.SetPassword("admin")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		for i := 0; i < 100; i++ {
			//topic := "/test_topic/" + strconv.Itoa(i)
			go publish(client, "123")
		}
		time.Sleep(1 * time.Second) // 暂停1秒
	}

}
func publish(client mqtt.Client, topic string) {
	// 初始化随机数生成器的种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机数
	randomNum := rand.Intn(21) // Intn返回一个[0, n)范围内的随机数

	client.Publish(topic, 0, false, strconv.Itoa(randomNum))
}
