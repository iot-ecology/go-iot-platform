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


func main2() {
	InitLog()
	mqtt.ERROR = log.New(getWriteSync(), "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(getWriteSync(), "[CRIT] ", 0)
	//mqtt.WARN = log.New(getWriteSync(), "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(getWriteSync(), "[DEBUG] ", 0)
	var broker = "192.168.3.101"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client2")
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

	for {

		for _, vc := range fime {
			publish(client, vc.Topic, vc.ID, vc.ID)

		}
		time.Sleep(30*time.Second)
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
	client.Publish(topic, 0, false, marshal)



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
