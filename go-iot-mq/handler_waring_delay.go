package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func HandlerWaringDelay(messages <-chan amqp.Delivery) {

	go func() {
		for d := range messages {

			HandlerWaringDelayStr(d)
			d.Ack(false)
		}
	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}

func HandlerWaringDelayStr(d amqp.Delivery) bool {
	var data []DataRowList
	err := json.Unmarshal(d.Body, &data)
	if err != nil {
		zap.S().Infof("Failed to unmarshal message: %s body = %s", err, string(d.Body))
		return false
	}

	for i := 0; i < len(data); i++ {
		// 解引用指针并访问切片中的元素
		row := (data)[i]

		handlerWaringDelayOnce(row)

	}
	return true
}

// handlerWaringDelayOnce 函数用于处理handlerWaringDelayOnce数据
// 参数：
//   - msg DataRowList: 包含反序列化后的消息的数据行列表
//
// 返回值：
//
//	无
func handlerWaringDelayOnce(msg DataRowList) {
	zap.S().Infof("处理 handlerWaringDelayOnce 数据: %+v", msg)
	uid := msg.DeviceUid
	mapping := getDelayParam(uid, msg.DataRows)
	background := context.Background()
	var scriptParam = make(map[string][]Tv)
	for _, param := range mapping {
		key := "signal_delay_warning:" + strconv.Itoa(param.MqttClientId) + ":" + strconv.Itoa(param.SignalId)
		zap.S().Infof("key = %s", key)
		members, _ := globalRedisClient.ZRevRangeWithScores(background, key, 0, -1).Result()
		var vs []Tv

		for _, member := range members {
			float, _ := strconv.ParseFloat(member.Member.(string), 64)
			vs = append(vs, Tv{Time: int64(member.Score), Value: float})
		}
		scriptParam[param.Name] = vs

	}
	script := getDelayScript(mapping)
	zap.S().Infof("脚本报警参数 = %+v", scriptParam)
	db := GMongoClient.Database(globalConfig.MongoConfig.Db)
	collection := db.Collection(globalConfig.MongoConfig.ScriptWaringCollection)
	var toInsert []interface{}
	for _, waring := range script {
		zap.S().Infof("key = %+v", waring)
		delayScript := runWaringDelayScript(waring.Script, scriptParam)
		toInsert = append(toInsert, bson.M{
			"device_uid":  uid,
			"param":       scriptParam,
			"script":      waring.Script,
			"value":       delayScript,
			"rule_id":     waring.ID,
			"insert_time": time.Now().Unix(),
			"up_time":     msg.Time,
		})
	}
	if toInsert != nil {

		one, err := collection.InsertMany(context.Background(), toInsert)
		if err != nil {
			zap.S().Errorf("插入数据失败 %+v", err)
		} else {
			zap.S().Infof("插入数据成功 %+v", one)
		}
		return
	}
}

// runWaringDelayScript 函数执行传入的JavaScript脚本，并将传入的参数map[string][]Tv传递给该脚本执行
// 参数：
//
//	script string - 要执行的JavaScript脚本
//	param map[string][]Tv - 传递给JavaScript脚本的参数
//
// 返回值：
//
//	bool - JavaScript脚本执行后返回的结果
func runWaringDelayScript(script string, param map[string][]Tv) bool {
	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		fmt.Println("JS代码有问题！")
	}
	var fn func(string2 map[string][]Tv) bool
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		fmt.Println("Js函数映射到 Go 函数失败！")
		panic(err)
	}
	a := fn(param)
	return a
}

// getDelayScript 从Redis中获取SignalDelayWaring信息列表
// 参数：
//
//	mapping []SignalDelayWaringParam - SignalDelayWaringParam类型的切片，用于从Redis中查询SignalDelayWaring信息
//
// 返回值：
//
//	[]SignalDelayWaring - SignalDelayWaring类型的切片，包含了从Redis中查询到的SignalDelayWaring信息
func getDelayScript(mapping []SignalDelayWaringParam) []SignalDelayWaring {
	var res []SignalDelayWaring

	for _, param := range mapping {
		id := param.SignalDelayWaringId
		val := globalRedisClient.HGet(context.Background(), "signal_delay_config", strconv.Itoa(id)).Val()
		var singw SignalDelayWaring
		err := json.Unmarshal([]byte(val), &singw)
		if err != nil {
			zap.S().Errorf("解析 json 异常 %+v", err)

		}
		res = append(res, singw)
	}
	return res
}

// getDelayParam 函数根据用户UID和DataRow切片从Redis中获取延迟报警参数
//
// 参数：
// uid string - 用户UID
// rows []DataRow - DataRow切片
//
// 返回值：
// []SignalDelayWaringParam - SignalDelayWaringParam切片，包含符合要求的延迟报警参数
func getDelayParam(uid string, rows []DataRow) []SignalDelayWaringParam {
	val := globalRedisClient.LRange(context.Background(), "delay_param", 0, -1).Val()
	var mapping []SignalDelayWaringParam
	for _, s := range val {
		var param SignalDelayWaringParam
		err := json.Unmarshal([]byte(s), &param)
		if err != nil {
			continue // 如果反序列化失败，跳过当前信号
		}
		if strconv.Itoa(param.MqttClientId) == uid && nameInDataRow(param.SignalName, rows) {

			mapping = append(mapping, param)
		}
	}
	return mapping
}

// nameInDataRow 函数用于判断给定的名称是否存在于DataRow切片中
//
// 参数：
//
//	name string - 需要查找的名称
//	rows []DataRow - DataRow切片
//
// 返回值：
//
//	bool - 如果找到名称，则返回true；否则返回false
func nameInDataRow(name string, rows []DataRow) bool {
	for _, row := range rows {
		if row.Name == name {
			return true
		}
	}

	return false

}
