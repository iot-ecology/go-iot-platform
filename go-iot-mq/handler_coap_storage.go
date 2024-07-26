package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// CoapMessage 用于处理tcp转发后的数据
type CoapMessage struct {
	Uid     string `json:"uid"`
	Message string `json:"message"`
}

// HandlerCoapDataStorage 函数处理从AMQP通道接收到的coap消息数据
// 参数：
//
//	messages <-chan amqp.Delivery：接收AMQP消息的通道
//
// 返回值：
//
//	无
func HandlerCoapDataStorage(messages <-chan amqp.Delivery) {

	go func() {

		for d := range messages {
			HandlerDataCoapStorageString(d)
			err := d.Ack(false)
			if err != nil {
				zap.S().Errorf("消息确认异常：%+v", err)

			}
		}
	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}

func HandlerDataCoapStorageString(d amqp.Delivery) {
	var msg CoapMessage
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		zap.S().Infof("Failed to unmarshal message: %s", err)
		return
	}
	zap.S().Infof("处理 pre_handler 数据 : %+v", msg)

	script := GetScriptRedisForCoap(msg.Uid)
	if script != "" {
		data := runScript(msg.Message, script)
		for i := 0; i < len(*data); i++ {
			row := (*data)[i]
			StorageDataRowList(row,"coap")
		}
		zap.S().Debugf("DataRowList: %+v", data)

		jsonData, err := json.Marshal(data)
		if err != nil {
			zap.S().Errorf("推送报警原始数据异常 %s", err)
			return
		}
		zap.S().Infof("推送报警原始数据: %s", jsonData)
		writeAPI.Flush()
		HandlerCoapLastTime(*data)
		PushToQueue("waring_handler", jsonData)
		PushToQueue("waring_delay_handler", jsonData)
		PushToQueue("transmit_handler", jsonData)
	} else {
		zap.S().Infof("执行脚本为空")
	}

}

// HandlerCoapLastTime 和上一次推送事件进行对比，判断是否超过阈值，如果超过则发送额外的消息通知
func HandlerCoapLastTime(data []DataRowList) {
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

		val := globalRedisClient.LRange(context.Background(), "coap_bind_device_info:"+deviceUid, 0, -1).Val()

		for _, s := range val {
			handlerCoapOne(s)
		}

	}

}

func handlerCoapOne(deviceUid string) bool {
	val := globalRedisClient.Get(context.Background(), "coap_bind_device_info:"+deviceUid).Val()
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

// GetScriptRedisForCoap 根据 http 的设备ID从Redis中获取对应的脚本
// 参数:
//
//	tcp id string - tcp id
//
// 返回值:
//
//	string - 对应的脚本
func GetScriptRedisForCoap(tcpId string) string {
	val := globalRedisClient.HGet(context.Background(), "struct:Coap", tcpId).Val()
	return val
}
