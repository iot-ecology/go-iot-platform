package main

import (
	"context"
	"encoding/json"
	"github.com/dop251/goja"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

// HandlerDataStorage 函数处理从AMQP通道接收到的MQTT消息数据
// 参数：
//
//	messages <-chan amqp.Delivery：接收AMQP消息的通道
//
// 返回值：
//
//	无
func HandlerDataStorage(messages <-chan amqp.Delivery) {

	go func() {

		for d := range messages {
			HandlerDataStorageString(d)
			err := d.Ack(false)
			if err != nil {
				zap.S().Errorf("消息确认异常：%+v", err)

			}
		}
	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}

func HandlerDataStorageString(d amqp.Delivery) {
	var msg MQTTMessage
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		zap.S().Infof("Failed to unmarshal message: %s", err)
		return
	}

	// 打印反序列化后的消息
	zap.S().Infof("处理 pre_handler 数据 : %+v", msg)

	script := GetScriptRedis(msg.MQTTClientID)
	if script != "" {
		data := runScript(msg.Message, script)
		// 这个原始数据
		// 1. 需要持久化到influxdb

		for i := 0; i < len(*data); i++ {
			// 解引用指针并访问切片中的元素
			row := (*data)[i]

			StorageDataRowList(row)
		}
		// 2. 发送到报警消息队列
		zap.S().Debugf("DataRowList: %+v", data)

		jsonData, err := json.Marshal(data)
		if err != nil {
			zap.S().Errorf("推送报警原始数据异常 %s", err)
			return
		}
		zap.S().Infof("推送报警原始数据: %s", jsonData)
		writeAPI.Flush()

		PushToQueue("waring_handler", jsonData)
		PushToQueue("waring_delay_handler", jsonData)
	} else {
		zap.S().Infof("执行脚本为空")
	}
	//err = d.Ack(false)
	//if err != nil {
	//	zap.S().Errorf("Failed to ack message: %s", err)
	//}

}

// StorageDataRowList 函数将DataRowList类型指针dt中的数据写入InfluxDB数据库
// 参数：
//
//	dt *DataRowList - DataRowList类型指针，存储待写入InfluxDB的数据
//
// 返回值：
//
//	无
func StorageDataRowList(dt DataRowList) {
	signal2 := GetMqttClientSignal2(dt.DeviceUid)
	timeFromUnix := time.Unix(dt.Time, 0)
	p := influxdb2.NewPointWithMeasurement(dt.DeviceUid).
		AddField("storage_time", time.Now().Unix()).
		AddField("push_time", dt.Time).
		SetTime(timeFromUnix)

	for _, row := range dt.DataRows {
		b := signal2[row.Name].Numb
		if b {
			float, _ := strconv.ParseFloat(row.Value, 64)
			p.AddField(strconv.Itoa(signal2[row.Name].ID), float)

		} else {
			p.AddField(strconv.Itoa(signal2[row.Name].ID), row.Value)

		}

		if signal2[row.Name].CacheSize > 0 {
			// 获取当前 ZSet 的大小
			currentSize := globalRedisClient.ZCard(context.Background(), "signal_delay_warning:"+dt.DeviceUid+":"+strconv.Itoa(signal2[row.Name].ID)).Val()

			// 如果 ZSet 的大小已经达到或超过配置的缓存大小，则移除第一个元素
			if currentSize >= signal2[row.Name].CacheSize {
				// 移除 ZSet 中分数最低的元素，即最早的元素
				i := signal2[row.Name].CacheSize + 1 - currentSize
				if i == 1 {

				} else {
					err := globalRedisClient.ZRemRangeByRank(context.Background(), "signal_delay_warning:"+dt.DeviceUid+":"+strconv.Itoa(signal2[row.Name].ID), 0, i).Err()
					if err != nil {
						// 处理错误
						zap.S().Errorf("移除 ZSet 元素异常：%+v", err)
					}
				}
			}

			// 写入缓存
			err := globalRedisClient.ZAdd(context.Background(), "signal_delay_warning:"+dt.DeviceUid+":"+strconv.Itoa(signal2[row.Name].ID), redis.Z{Score: float64(dt.Time), Member: row.Value}).Err()
			if err != nil {
				// 处理错误
				zap.S().Errorf("写入 ZSet 元素异常：%+v", err)
			}
		}

	}

	writeAPI.WritePoint(p)
	// fixme: 持久化完成后需要进一步推送给delay的报警队列，找出关联的signal_delay报警配置表的id

}

// runScript 函数接收两个字符串参数，param 和 script，返回一个指向 DataRowList 类型的指针
//
// 参数：
// param：string 类型，传递给 JS 脚本的参数
// script：string 类型，待执行的 JS 脚本
//
// 返回值：
// *DataRowList 类型指针，JS 脚本执行后的结果，如果执行失败则返回 nil
func runScript(param string, script string) *[]DataRowList {

	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Errorf("JS代码有问题！")
		return nil
	}
	var fn func(string2 string) *[]DataRowList
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		zap.S().Errorf("Js函数映射到 Go 函数失败！")
		return nil
	}
	a := fn(param)
	return a

}

// GetMqttClientSignal 函数根据MQTT客户端ID获取对应的信号映射表
// 参数：
//
//	mqtt_client_id string - MQTT客户端ID
//
// 返回值：
//
//	map[string]bool - 信号映射表，其中key为信号名称，value表示信号类型是否为数字类型（忽略大小写）
func GetMqttClientSignal(mqttClientId string) (map[int]bool, map[int]int64, map[string]int) {
	background := context.Background()
	result, err := globalRedisClient.LRange(background, "signal:"+mqttClientId, 0, -1).Result()
	if err != nil {
		// 处理错误，例如记录日志或返回错误
		zap.S().Errorf("获取信号映射表失败：%+v", err)
	}
	mapping := make(map[int]bool)
	mappingName := make(map[string]int)
	cacheSizeMapping := make(map[int]int64)
	for _, strSignal := range result {
		var signal Signal
		err := json.Unmarshal([]byte(strSignal), &signal)
		if err != nil {
			continue // 如果反序列化失败，跳过当前信号
		}

		mapping[signal.ID] = strings.EqualFold(signal.Type, "数字")
		cacheSizeMapping[signal.ID] = signal.CacheSize
		mappingName[signal.Name] = signal.ID
	}
	return mapping, cacheSizeMapping, mappingName

}

func GetMqttClientSignal2(mqttClientId string) map[string]signalMapping {
	background := context.Background()
	result, err := globalRedisClient.LRange(background, "signal:"+mqttClientId, 0, -1).Result()
	if err != nil {
		// 处理错误，例如记录日志或返回错误
		zap.S().Errorf("获取信号映射表失败：%+v", err)
	}
	mapping := make(map[string]signalMapping)
	for _, strSignal := range result {
		var signal Signal
		err := json.Unmarshal([]byte(strSignal), &signal)
		if err != nil {
			continue // 如果反序列化失败，跳过当前信号
		}

		mapping[signal.Name] = signalMapping{
			CacheSize: signal.CacheSize,
			ID:        signal.ID,
			Numb:      strings.EqualFold(signal.Type, "数字"),
		}
	}
	return mapping

}

type signalMapping struct {
	CacheSize int64
	ID        int
	Numb      bool
}
