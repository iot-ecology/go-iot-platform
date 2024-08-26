package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	InitLog()
	mqtt.ERROR = log.New(getWriteSync(), "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(getWriteSync(), "[CRIT] ", 0)
	mqtt.WARN = log.New(getWriteSync(), "[WARN]  ", 0)
	mqtt.DEBUG = log.New(getWriteSync(), "[DEBUG] ", 0)
	var broker = "172.17.0.1"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("admin")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetAutoReconnect(true)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fime := readFime("1.txt")

	go c(client,fime[0])

}
func c(client mqtt.Client , vc Vc) {

		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			publish(client, vc.Topic, vc.ID, vc.ID)
		}

}

func readFime(path string) []Vc {
	// 读取 path 每一行用空格分割分隔的数据 第一个是 Topic 第二个是ID

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var vcs []Vc
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) != 2 {
			log.Printf("Invalid line format: %s", line)
			continue
		}
		topic := fields[0]
		id, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Printf("Invalid ID: %s", fields[1])
			continue
		}
		vcs = append(vcs, Vc{Topic: topic, ID: id})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return vcs
}

type Vc struct {
	Topic string
	ID    int
}

func publish(client mqtt.Client, topic string, i int, i2 int) {
	// 初始化随机数生成器的种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机数

	var dataRows []DataRow
	for i3 := range 200 {
		randomNum := rand.Intn(21) // Intn返回一个[0, n)范围内的随机数
		dataRows = append(dataRows, DataRow{
			Name:  "信号-" + strconv.Itoa(i3),
			Value: strconv.Itoa(randomNum),
		})
	}

	//
	DataRowList := DataRowList{
		Time:               time.Now().Unix(),
		DeviceUid:          strconv.Itoa(i2),
		IdentificationCode: strconv.Itoa(i2),
		DataRows:           dataRows,
		Nc:                 strconv.Itoa(i2),
	}

	marshal, _ := json.Marshal(DataRowList)

	zap.S().Infof("发送消息: %s  消息主题: %s\n", DataRowList.Time, topic)
	token := client.Publish(topic, 0, false, marshal)

	if token.Wait() && token.Error() != nil {
		zap.S().Error(token.Error())
	}


}

type DataRowList struct {
	Time               int64     `json:"Time"`               // 秒级时间戳
	DeviceUid          string    `json:"DeviceUid"`          // 能够产生网络通讯的唯一编码
	IdentificationCode string    `json:"IdentificationCode"` // 设备标识码
	DataRows           []DataRow `json:"DataRows"`
	Nc                 string    `json:"Nc"`
}
type DataRow struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}
