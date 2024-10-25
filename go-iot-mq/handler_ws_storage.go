package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// WsMessage 用于处理ws转发后的数据
type WsMessage struct {
	Uid     string `json:"uid"`
	Message string `json:"message"`
}

// HandlerWsDataStorage 函数处理从AMQP通道接收到的websocket消息数据
//
// 参数：
//
//	messages <-chan amqp.Delivery：接收AMQP消息的通道
//
// 返回值：
//
//	无
func HandlerWsDataStorage(messages <-chan amqp.Delivery) {

	go func() {

		for d := range messages {
			HandlerDataWsStorageString(d)
			err := d.Ack(false)
			if err != nil {
				zap.S().Errorf("消息确认异常：%+v", err)

			}
		}
	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}

func HandlerDataWsStorageString(d amqp.Delivery) {
	var msg WsMessage
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		zap.S().Infof("Failed to unmarshal message: %s", err)
		return
	}
	zap.S().Infof("处理 pre_ws_handler 数据 : %+v", msg)

	script := GetScriptRedisForWs(msg.Uid)
	if script != "" {
		data := runScript(msg.Message, script)
		if data == nil {
			zap.S().Infof("执行脚本为空")
			return
		}
		for i := 0; i < len(*data); i++ {
			row := (*data)[i]
			(*data)[i].Protocol = "WebSocket"
			StorageDataRowList(row, "WebSocket")
		}
		zap.S().Debugf("DataRowList: %+v", data)

		jsonData, err := json.Marshal(data)
		if err != nil {
			zap.S().Errorf("推送报警原始数据异常 %s", err)
			return
		}
		zap.S().Infof("推送报警原始数据: %s", jsonData)
		HandlerWebsocketLastTime(*data)
		PushToQueue("waring_handler", jsonData)
		PushToQueue("waring_delay_handler", jsonData)
		PushToQueue("transmit_handler", jsonData)
	} else {
		zap.S().Infof("执行脚本为空")
	}

}

// HandlerWebsocketLastTime 和上一次推送事件进行对比，判断是否超过阈值，如果超过则发送额外的消息通知
func HandlerWebsocketLastTime(data []DataRowList) {
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

		val := globalRedisClient.LRange(context.Background(), "ws_bind_device_info:"+deviceUid, 0, -1).Val()

		for _, s := range val {
			handlerWebsocketOne(s)
		}

	}

}

func handlerWebsocketOne(deviceUid string) bool {
	val := globalRedisClient.Get(context.Background(), "ws_bind_device_info:"+deviceUid).Val()
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

// GetScriptRedisForWs 根据 http 的设备ID从Redis中获取对应的脚本
//
// 参数:
//
//	- tcpid string tcp id
//
// 返回值:
//
//	- string 对应的脚本
func GetScriptRedisForWs(tcpId string) string {
	val := globalRedisClient.HGet(context.Background(), "struct:Websocket", tcpId).Val()
	return val
}
