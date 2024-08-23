package main

import (
	"encoding/json"
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
			topic := "/test_topic/" + strconv.Itoa(i)
			go publish(client, topic,i,i)
		}
		time.Sleep(1 * time.Second) // 暂停1秒
	}

}
func publish(client mqtt.Client, topic string, i int, i2 int) {
	// 初始化随机数生成器的种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机数

	var dataRows []DataRow
	for i3 := range 200 {
		randomNum := rand.Intn(21) // Intn返回一个[0, n)范围内的随机数
		dataRows  = append(dataRows, DataRow{
			Name:  "信号-" + strconv.Itoa(i3) ,
			Value: strconv.Itoa(randomNum),
		})
	}


	//
	DataRowList := DataRowList{
		Time:               time.Now().Unix(),
		DeviceUid:          strconv.Itoa(i2+1),
		IdentificationCode: strconv.Itoa(i2+1),
		DataRows:           dataRows,
		Nc:                 strconv.Itoa(i2),
	}

	marshal, _ := json.Marshal(DataRowList)

	client.Publish(topic, 0, false, marshal)
}



type DataRowList struct {
	Time      int64     `json:"Time"`       // 秒级时间戳
	DeviceUid string    `json:"DeviceUid"` // 能够产生网络通讯的唯一编码
	IdentificationCode string `json:"IdentificationCode"` // 设备标识码
	DataRows  []DataRow `json:"DataRows"`
	Nc        string    `json:"Nc"`
}
type DataRow struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}
