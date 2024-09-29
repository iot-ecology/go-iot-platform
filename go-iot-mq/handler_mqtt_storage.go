// 用于处理MQTT转发过来的数据

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
//
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
			go funcName(d)
		}
	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}

func funcName(d amqp.Delivery) {
	HandlerDataStorageString(d)
	err := d.Ack(false)
	if err != nil {
		zap.S().Errorf("消息确认异常：%+v", err)

	}
}

// HandlerDataStorageString 是一个处理来自AMQP的消息的函数
// 它从传入的AMQP交付对象中反序列化MQTT消息
// 如果反序列化失败，则记录错误并返回
// 如果成功，则打印反序列化后的MQTT消息
// 它通过调用GetScriptRedis获取MQTT客户端的脚本
// 如果脚本不为空，则执行脚本并获取数据
// 遍历数据，并将每一行存储到InfluxDB
// 最后，将数据存储为JSON格式，并推送到两个消息队列中
// 如果脚本为空，则记录消息
func HandlerDataStorageString(d amqp.Delivery) {
	var msg MQTTMessage
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		zap.S().Errorf("Failed to unmarshal message: %s", err)
		return
	}
	zap.S().Debugf("处理 pre_handler 数据 : %+v", msg)

	script := GetScriptRedis(msg.MQTTClientID)
	if script != "" {
		data := runScript(msg.Message, script)

		if data == nil {
			zap.S().Errorf("执行脚本为空")
			return
		}
		for i := 0; i < len(*data); i++ {
			row := (*data)[i]
			(*data)[i].Protocol = "MQTT"
			StorageDataRowList(row, "MQTT")
		}
		zap.S().Debugf("DataRowList: %+v", data)

		jsonData, err := json.Marshal(data)
		if err != nil {
			zap.S().Errorf("推送报警原始数据异常 %s", err)
			return
		}
		zap.S().Debugf("推送报警原始数据: %s", jsonData)
		HandlerMqttLastTime(*data)
		PushToQueue("waring_handler", jsonData)
		PushToQueue("waring_delay_handler", jsonData)
		PushToQueue("transmit_handler", jsonData)
	} else {
		zap.S().Errorf("执行脚本为空")
	}

}

// HandlerMqttLastTime 和上一次推送事件进行对比，判断是否超过阈值，如果超过则发送额外的消息通知
func HandlerMqttLastTime(data []DataRowList) {
	if len(data) == 0 {
		return
	}

	var deviceUid = data[0].DeviceUid
	key := "last_push_time:" + deviceUid
	// 1. 从redis中获取这个设备上次推送的时间
	lastTime, err := globalRedisClient.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		zap.S().Errorf("获取设备上次推送时间异常：%+v", err)
		return
	}
	now := time.Now().Unix()

	// 如果没有这个时间则设置时间(当前时间）
	if errors.Is(err, redis.Nil) {
		err := globalRedisClient.Set(context.Background(), key, now, 0).Err()
		if err != nil {
			zap.S().Errorf("设置设备上次推送时间异常：%+v", err)
			return
		}
		lastTime = fmt.Sprintf("%d", now)
	}

	if lastTime != fmt.Sprintf("%d", now) {

		val := globalRedisClient.LRange(context.Background(), "mqtt_client_id_bind_device_info:"+deviceUid, 0, -1).Val()

		for _, s := range val {
			handlerOne(s)
		}

	}

}

func handlerOne(deviceUid string) bool {
	val := globalRedisClient.Get(context.Background(), "mqtt_client_id_bind_device_info:"+deviceUid).Val()
	if val == "" {
		return true
	}
	parseUint, _ := strconv.ParseUint(val, 10, 64)
	withRedis := FindByIdWithRedis(parseUint)
	if withRedis == nil {
		return true
	}
	globalRedisClient.Expire(context.Background(), "Device_Off_Message:"+deviceUid, time.Duration(withRedis.PushInterval)*time.Second)
	return false
}
func FindByIdWithRedis(id uint64) *DeviceInfo {
	val := globalRedisClient.HGet(context.Background(), "struct:device_info", strconv.Itoa(int(id))).Val()

	var res DeviceInfo
	err := json.Unmarshal([]byte(val), &res)
	if err != nil {
		return nil
	}
	return &res
}

func genMeasurement(dt DataRowList, protocol string) string {
	return protocol + "_" + dt.DeviceUid + "_" + dt.IdentificationCode
}

// CalcBucketName 函数根据前缀、协议和id计算桶名
// prefix: 桶名前缀
// protocol: 使用的协议
// id: 桶的ID
// 返回值: 计算得到的桶名
func CalcBucketName(prefix, protocol string, id uint) string {
	return prefix + "_" + protocol + "_" + strconv.Itoa(int(id%100))
}

// StorageDataRowList 函数将DataRowList类型指针dt中的数据写入InfluxDB数据库
// 参数：
//
//	dt *DataRowList - DataRowList类型指针，存储待写入InfluxDB的数据
//
// 返回值：
//
//	无
func StorageDataRowList(dt DataRowList, protocol string) {
	signal2 := GetMqttClientSignal2(dt.DeviceUid, dt.IdentificationCode)
	zap.S().Debugf("获取的mqtt信号数据signal2: %+v", signal2)
	loc, _ := time.LoadLocation("Asia/Shanghai")

	timeFromUnix := time.Unix(dt.Time, 0).In(loc)

	if CheckPushTime(dt.Protocol, dt.IdentificationCode, dt.DeviceUid, timeFromUnix.Unix()) {
		zap.S().Errorf("推送时间间隔异常")
	}
	i, err := strconv.Atoi(dt.DeviceUid)
	if err != nil {
		zap.S().Debugf("转换错误: %+v", err)
	} else {
		zap.S().Debugf("转换后的整数: %+v", i)
	}

	bucketName := CalcBucketName(globalConfig.InfluxConfig.Bucket, protocol, uint(i))
	writeAPI := GlobalInfluxDbClient.WriteAPI(globalConfig.InfluxConfig.Org,
		bucketName)

	measurement := genMeasurement(dt, protocol)
	zap.S().Infof("bucketName = %s , measurement = %s",bucketName,measurement)
	unix := time.Now().Unix()
	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("storage_time", unix).
		AddField("push_time", dt.Time).
		AddField("time-sub",unix-dt.Time).
		SetTime(timeFromUnix)

	for _, row := range dt.DataRows {

		b := signal2[row.Name].Numb
		if b {
			float, _ := strconv.ParseFloat(row.Value, 64)
			p.AddField(strconv.Itoa(signal2[row.Name].ID), float)

		} else {
			p.AddField(strconv.Itoa(signal2[row.Name].ID), row.Value)

		}
		zap.S().Debugf("当前信号的的CacheSize:%+v=============rowName:%+v", signal2[row.Name].CacheSize, row.Name)
		if signal2[row.Name].CacheSize > 0 {
			// 获取当前 ZSet 的大小
			currentSize := globalRedisClient.ZCard(context.Background(),
				"signal_delay_warning:"+dt.DeviceUid+":"+dt.IdentificationCode+":"+strconv.Itoa(signal2[row.Name].ID)).Val()
			zap.S().Debugf("当前signal_delay_warning的大小: %+v", currentSize)
			// 如果 ZSet 的大小已经达到或超过配置的缓存大小，则移除第一个元素
			if currentSize >= signal2[row.Name].CacheSize {
				zap.S().Debugf("当前signal_delay_warning的currentSize大于等于配置大小:%+v", signal2[row.Name].CacheSize)
				// 移除 ZSet 中分数最低的元素，即最早的元素
				i := signal2[row.Name].CacheSize + 1 - currentSize
				zap.S().Debugf("计算后的i: %+v", i)
				if i == 1 {
					zap.S().Debugf("计算后的i的值为1")
				} else {
					zap.S().Debugf("开始移除之前的元素")
					err := globalRedisClient.ZRemRangeByRank(context.Background(),
						"signal_delay_warning:"+dt.DeviceUid+":"+dt.IdentificationCode+":"+strconv.Itoa(
							signal2[row.Name].ID), 0,
						i-1).Err()
					if err != nil {
						// 处理错误
						zap.S().Errorf("移除 ZSet 元素异常：%+v", err)
					}
				}
			} else {
				zap.S().Errorf("当前大小未超过配置大小,写入缓存")
				// 写入缓存
				// 根据zset的特效,如果value一致的话,则会修改score,此处体现为修改了该值的时间,也就是说最新的值和之前的值相同的话只会保留最新时间的这一份
				err := globalRedisClient.ZAdd(context.Background(),
					"signal_delay_warning:"+dt.DeviceUid+":"+dt.IdentificationCode+":"+strconv.Itoa(signal2[row.
						Name].
						ID),
					redis.Z{Score: float64(dt.Time), Member: row.Value}).Err()
				if err != nil {
					// 处理错误
					zap.S().Errorf("写入 ZSet 元素异常：%+v", err)
				}
			}

		}

	}

	SetPushTime(dt.Protocol,dt.IdentificationCode,dt.DeviceUid,dt.Time)
	writeAPI.WritePoint(p)
	writeAPI.Flush()

}

// runScript 函数接收两个字符串参数，param 和 script，返回一个指向 DataRowList 类型的指针
//
// 参数：
// 	- param：string 类型，传递给 JS 脚本的参数
// 	- script：string 类型，待执行的 JS 脚本
//
// 返回值：
// 	- *DataRowList 类型指针，JS 脚本执行后的结果，如果执行失败则返回 nil
func runScript(param string, script string) *[]DataRowList {
	vm := goja.New()

	// 执行 JavaScript 脚本
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Errorf("JS代码有问题: %v", err)
		return nil
	}

	// 将 JavaScript 中的 main 函数映射到 Go 的 fn 函数
	var fn func(string) *[]DataRowList
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		zap.S().Errorf("Js函数映射到 Go 函数失败: %v", err)
		return nil
	}

	// 使用 defer 和 recover 来捕获 fn 函数中的 panic
	var result *[]DataRowList
	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorf("在执行 JavaScript 函数时发生 panic: %v", r)
			// 这里可以进行一些清理工作或者返回一个特定的错误结果
			result = nil // 或者设置为一个特定的错误结果
		}
	}()

	// 调用映射的函数
	result = fn(param)


	return result
}

func GetMqttClientSignal2(mqttClientId, IdentificationCode string) map[string]signalMapping {
	background := context.Background()
	result, err := globalRedisClient.LRange(background, "signal:"+mqttClientId+":"+IdentificationCode, 0, -1).Result()
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
