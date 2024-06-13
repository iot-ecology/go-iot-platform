package main

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// HandlerWaring 是一个处理预警消息的函数
//
// 参数：
// messages <-chan amqp.Delivery: 接收AMQP消息的通道
//
// 返回值：
// 无
func HandlerWaring(messages <-chan amqp.Delivery) {

	go func() {
		for d := range messages {

			HandlerWaringString(d)
			err := d.Ack(false)
			if err != nil {
				zap.S().Errorf("消息确认异常：%+v", err)

			}

		}
	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}

func HandlerWaringString(d amqp.Delivery) bool {
	var data []DataRowList
	err := json.Unmarshal(d.Body, &data)
	if err != nil {
		zap.S().Infof("Failed to unmarshal message: %s body = %s", err, string(d.Body))
		return false
	}
	for i := 0; i < len(data); i++ {
		// 解引用指针并访问切片中的元素
		row := (data)[i]

		handlerWaringOnce(row)

	}
	return true
}

// handlerWaringOnce 处理警告处理器数据的函数
//
// 参数：
// msg DataRowList - 包含反序列化后的消息的数据行列表
//
// 返回值：
// bool - 表示是否处理成功
func handlerWaringOnce(msg DataRowList) {
	// 打印反序列化后的消息
	zap.S().Debugf("处理 waring_handler 数据: %+v", msg)

	uid := msg.DeviceUid
	// 1. 根据设备UID（mqtt客户端ID）获取所有信号

	mapping := getMqttClientMappingSignalWarningConfig(uid)
	db := GMongoClient.Database(globalConfig.MongoConfig.Db)
	collection := db.Collection(globalConfig.MongoConfig.WaringCollection)

	var toInsert []interface{}
	for _, row := range msg.DataRows {
		configs := mapping[row.Name]

		floatValue, err := strconv.ParseFloat(row.Value, 64)
		if err != nil {
			zap.S().Errorf("字符串转换出错: %s", err)
			continue
		}
		for _, config := range configs {

			// fixme: 消息推送
			zap.S().Infof("配置ID: %d 信号名称: %s, 范围小值: %f, 范围大值: %f, 1=范围内报警 0=范围外报警: %d\n", config.ID, row.Name, config.Min, config.Max, config.InOrOut)

			if config.InOrOut == 1 {
				if config.Min <= floatValue && floatValue <= config.Max {
					// 在范围内，根据需求执行操作
					zap.S().Infof("当前信号 %s  值在范围内: %+v 命中规则ID %d", row.Name, floatValue, config.ID)
					toInsert = append(toInsert, bson.M{
						"device_uid":  uid,
						"signal_name": row.Name,
						"signal_id":   config.SignalId,
						"value":       floatValue,
						"rule_id":     config.ID,
						"insert_time": time.Now().Unix(),
						"up_time":     msg.Time,
					})
				}
			} else {
				if floatValue < config.Min || floatValue > config.Max {
					// 范围外报警
					zap.S().Infof("当前信号 %s 范围外报警: %+v 命中规则ID %d", row.Name, floatValue, config.ID)
					toInsert = append(toInsert, bson.M{
						"device_uid":  uid,
						"signal_name": row.Name,
						"signal_id":   config.SignalId,
						"value":       floatValue,
						"rule_id":     config.ID,
						"insert_time": time.Now().Unix(),
						"up_time":     msg.Time,
					})
				}
			}
		}

	}
	_, err := collection.InsertMany(context.Background(), toInsert)
	if err != nil {
		zap.S().Errorf("消息确认异常：%+v", err)
	}
}

// getMqttClientMappingSignalWarningConfig 根据 MQTT 客户端 ID 获取信号警告配置的映射
// 参数:
//
//	mqtt_client_id string - MQTT 客户端 ID
//
// 返回值:
//
//	map[string][]SignalWaringConfig - 信号名称到信号警告配置切片的映射
func getMqttClientMappingSignalWarningConfig(mqttClientId string) map[string][]SignalWaringConfig {
	background := context.Background()
	result, err := globalRedisClient.LRange(background, "signal:"+mqttClientId, 0, -1).Result()
	if err != nil {
		// 处理错误，例如记录日志或返回错误
		zap.S().Errorf("获取信号列表失败: %+v", err)
	}

	// 创建一个映射，用于存放 signal.Name 到 swcs 的映射
	mapping := make(map[string][]SignalWaringConfig)

	for _, strSignal := range result {
		var signal Signal
		err := json.Unmarshal([]byte(strSignal), &signal)
		if err != nil {
			continue // 如果反序列化失败，跳过当前信号
		}

		result2, err := globalRedisClient.LRange(background, "waring:"+strconv.Itoa(signal.ID), 0, -1).Result()
		if err != nil {
			continue // 如果获取 warning 列表失败，跳过当前信号
		}

		var swcs []SignalWaringConfig
		for _, sw := range result2 {
			var swc SignalWaringConfig
			err := json.Unmarshal([]byte(sw), &swc)
			if err != nil {
				continue // 如果反序列化失败，跳过当前警告配置
			}
			swcs = append(swcs, swc)
		}

		// 将解析后的 swcs 切片与 signal.Name 关联，并存储到映射中
		mapping[signal.Name] = swcs
	}

	return mapping
}
