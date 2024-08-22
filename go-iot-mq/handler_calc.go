package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dop251/goja"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// HandlerCalc 处理计算消息队列中的消息
//
// 参数：
//   - messages <-chan amqp.Delivery：消息队列通道，用于接收消息
//
// 返回值：
//   - 无
func HandlerCalc(messages <-chan amqp.Delivery) {

	go func() {
		for d := range messages {

			HandlerCalcStr(d)
			err := d.Ack(false)
			if err != nil {
				zap.S().Errorf("消息确认异常：%+v", err)
			}

		}

	}()

	zap.S().Infof(" [*] Waiting for messages. To exit press CTRL+C")
}
func calcMeasurement(deviceUid int, IdentificationCode, protocol string) string {
	zap.S().Infof("calcMeasurement 开始, deviceUid = %v, IdentificationCode = %v, protocol = %v", deviceUid, IdentificationCode, protocol)
	return protocol + "_" + strconv.Itoa(deviceUid) + "_" + IdentificationCode
}

func HandlerCalcStr(d amqp.Delivery) bool {
	var myMap map[string]int64

	// 使用Unmarshal函数将JSON数据解码到map中
	err := json.Unmarshal(d.Body, &myMap)
	if err != nil {
		log.Fatal(err)
	}

	result, err := globalRedisClient.HGet(context.Background(), "calc_cache", strconv.Itoa(int(myMap["id"]))).Result()
	if err != nil {
		zap.S().Errorf("转化异常 %+v", err)
		return true

	}
	var ccc CalcCache

	err = json.Unmarshal([]byte(result), &ccc)
	if err != nil {
		zap.S().Infof("Failed to unmarshal message: %s", err)
		return true
	}
	nextTime := getNextTime(ccc.Cron)
	zap.S().Infof("执行时间 = %d", nextTime)

	val := globalRedisClient.Get(context.Background(), "calc_queue_param:"+strconv.Itoa(int(myMap["id"]))).Val()
	preTime, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		zap.S().Errorf("转换错误：%s", err)
		return true

	}
	var m = make(map[string]any)
	for _, cache := range ccc.Param {
		// todo: 查询逻辑调整。  Measurement 不能直接和 MQTT_CLIENT_ID 对应了
		if "原始" == cache.Reduce {
			var fd []string
			fd = append(fd, strconv.Itoa(cache.SignalId))
			config := InfluxQueryConfig{}
			config.Bucket = globalConfig.InfluxConfig.Bucket
			config.Measurement = calcMeasurement(cache.DeviceUid, cache.IdentificationCode, cache.Protocol)
			config.Fields = fd
			config.Aggregation = AggregationConfig{
				Every:       1,
				Function:    "mean",
				CreateEmpty: false,
			}
			config.StartTime = preTime - ccc.Offset
			config.EndTime = preTime
			config.Reduce = cache.Reduce
			query := config.GenerateFluxQuery()
			zap.S().Infof("influxdb query line = %s", query)

			result, err := GlobalInfluxDbClient.QueryAPI(globalConfig.InfluxConfig.Org).Query(context.Background(), query)
			if err != nil {
				zap.S().Error(err)
			}
			v := make(map[int64]float64)
			for result.Next() {
				if result.TableChanged() {
					fmt.Printf("table: %s\n", result.TableMetadata().String())
				}
				values := result.Record().Values()
				fmt.Printf("value: %+v\n", values)
				t := values["_time"].(time.Time)
				println(t.Unix())
				v[t.Unix()] = values["_value"].(float64)

			}
			m[cache.Name] = v
		} else {
			var fd []string
			fd = append(fd, strconv.Itoa(cache.SignalId))

			config := InfluxQueryConfig{}
			config.Bucket = globalConfig.InfluxConfig.Bucket
			config.Measurement = calcMeasurement(cache.DeviceUid, cache.IdentificationCode, cache.Protocol)
			config.Fields = fd

			config.StartTime = preTime - ccc.Offset
			config.EndTime = preTime
			//config.StartTime = 1716134400
			//config.EndTime = 1716449111
			config.Reduce = cache.Reduce
			query := config.GenerateFluxReduce()

			zap.S().Infof("influxdb query line = %s", query)
			result, err := GlobalInfluxDbClient.QueryAPI(globalConfig.InfluxConfig.Org).Query(context.Background(), query)
			if err != nil {
				zap.S().Error(err)
				return true

			}

			for result.Next() {
				if result.TableChanged() {
					zap.S().Errorf("table: %s\n", result.TableMetadata().String())
				}
				values := result.Record().Values()
				zap.S().Infof("value: %+v", values)
				m[cache.Name] = values["_value"].(float64)

			}
		}
	}
	scriot := runCalcScript(m, ccc.Script)
	zap.S().Infof("结算结果 = %+v", scriot)

	// 获取数据库和集合
	db := GMongoClient.Database(globalConfig.MongoConfig.Db)
	// fixme: 暂时使用固定集合名，后续需要改成根据规则ID动态获取
	name := CalcCollectionName(globalConfig.MongoConfig.Collection, ccc.ID)
	CheckCollectionAndCreate(globalConfig.MongoConfig.Collection, name)
	collection := db.Collection(name)

	// 插入数据
	insertResult, err := collection.InsertOne(context.Background(), bson.M{

		"calc_rule_id": myMap["id"],
		"ex_time":      nextTime,
		"start_time":   nextTime - ccc.Offset,
		"end_time":     nextTime,
		"param":        m,
		"script":       ccc.Script,
		"result":       scriot})
	if err != nil {
		log.Fatal(err)
		return true
	}
	zap.S().Infof("写入完成 %+v", insertResult)

	subSec := nextTime - time.Now().Unix()
	zap.S().Infof("时间差 %d", subSec)

	sendData := make(map[string]uint)

	sendData["id"] = ccc.ID

	jsonData, err := json.Marshal(sendData)
	zap.S().Infof("jsonData  = %s", jsonData)
	if err != nil {
		zap.S().Errorf("数据异常 %s", err)
		return true
	}
	// fixme: 修订redis的时间
	//  1. 删除老的时间数据

	jsonStr, _ := json.Marshal(myMap)
	z := redis.Z{Score: float64(nextTime), Member: jsonStr}
	globalRedisClient.Set(context.Background(), "calc_queue_param:"+strconv.Itoa(int(myMap["id"])), nextTime, 0)
	err = globalRedisClient.ZAdd(context.Background(), "calc_queue", z).Err()
	if err != nil {

		zap.S().Errorf("redis zadd 异常 %s", err)
	}

	return false
}

// runCalcScript 函数用于执行JavaScript脚本，并返回计算结果
//
// 参数：
//
//   - param: map[string]float64类型，表示传递给JavaScript脚本的参数，其中键为参数名，值为参数值
//
//   - script: string类型，表示待执行的JavaScript脚本
//
// 返回值：
//
//	- map[string]interface{}类型，表示JavaScript脚本执行后的结果，其中键为结果名，值为结果值
func runCalcScript(param map[string]any, script string) map[string]interface{} {
	vm := goja.New()

	// 执行 JavaScript 脚本
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Error("JS代码有问题！", zap.Error(err))
		return nil
	}

	// 将 JavaScript 中的 main 函数映射到 Go 的 fn 函数
	var fn func(map[string]any) map[string]interface{}
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		zap.S().Error("Js函数映射到 Go 函数失败！", zap.Error(err))
		return nil
	}

	// 使用 defer 和 recover 来捕获 fn 函数中的 panic
	result := make(map[string]interface{})
	defer func() {
		if r := recover(); r != nil {
			zap.S().Error("在执行 JavaScript 函数时发生 panic:", zap.Any("panic value", r))
			// 可以选择返回空的 map 或者包含错误信息的 map
			result["error"] = fmt.Sprintf("panic occurred: %v", r)
		}
	}()

	// 调用映射的函数
	result = fn(param)
	return result
}

// getNextTime 函数根据传入的cron表达式，返回下一次执行时间的时间戳（秒）
//
// 参数：
//
//	- cronExpr string cron表达式
//
// 返回值：
//
//	- int64  下一次执行时间的时间戳（秒）
func getNextTime(cronExpr string) int64 {

	// 解析cron表达式
	schedule, err := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).Parse(cronExpr)
	if err != nil {
		zap.S().Error("解析cron表达式失败:", err)
		return 0
	}

	// 获取当前时间
	now := time.Now()

	// 计算下一次执行时间
	nextTime := schedule.Next(now)

	// 将下一次执行时间格式化为字符串
	format := nextTime.Format("2006-01-02 15:04:05")
	zap.S().Infof("下一次执行时间: %s", format)

	// 将下一次执行时间转换为时间戳（秒）
	nextTimestamp := nextTime.Unix()
	return nextTimestamp
}
