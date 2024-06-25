package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"log"
	"strconv"
	"time"
)

type CalcRunBiz struct{}

// Start 根据传入的id启动计算任务
//
// 参数：
// id：计算规则id
//
// 返回值：
// bool：启动计算任务是否成功，成功返回true，否则返回false
func (b CalcRunBiz) Start(id any) bool {

	var calcRule models.CalcRule

	result := glob.GDb.First(&calcRule, id)
	if result.Error != nil {
		glob.GLog.Error("查询异常", zap.Any("err", result.Error))
		panic(result.Error)
	}
	myMap := make(map[string]uint)

	myMap["id"] = calcRule.ID

	nextTime := getNextTime(calcRule.Cron)
	// 放到 zset 队列里面

	jsonStr, _ := json.Marshal(myMap)
	z := redis.Z{Score: float64(nextTime), Member: jsonStr}
	err := glob.GRedis.ZAdd(context.Background(), "calc_queue", z).Err()
	if err != nil {

		glob.GLog.Error("redis zadd 异常", zap.Any("err", err))
		panic(err)
	}

	mup := map[string]interface{}{
		"Start": true, // 更新 Start 字段为 true
	}
	result = glob.GDb.Table("calc_rules").Where("id = ?", calcRule.ID).Updates(mup)

	if result.Error != nil {
		glob.GLog.Error("更新异常", zap.Any("err", result.Error))
		return false
	}
	return true

}

// RefreshRule 根据id刷新计算规则缓存
//
// 参数:
// id: 计算规则id
//
// 返回值:
// 无
func (b CalcRunBiz) RefreshRule(id any) {
	var calcRule models.CalcRule

	result := glob.GDb.First(&calcRule, id)
	if result.Error != nil {
		glob.GLog.Error("查询异常", zap.Any("err", result.Error))
		panic(result.Error)
	}

	var calcParams []models.CalcParam
	result = glob.GDb.Where("calc_rule_id = ?", calcRule.ID).Find(&calcParams)

	// 检查是否有错误
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	var m []servlet.CalcParamCache

	for _, param := range calcParams {
		m = append(m, servlet.CalcParamCache{
			MqttClientId: param.MqttClientId,
			SignalId:     param.SignalId,
			Name:         param.Name,
			SignalName:   param.SignalName,
			Reduce:       param.Reduce,
			CalcRuleId:   param.CalcRuleId,
		})
	}

	cache := servlet.CalcCache{
		ID:     calcRule.ID,
		Param:  m,
		Cron:   calcRule.Cron,
		Script: calcRule.Script,
		Offset: calcRule.Offset,
	}

	jsonData, err := json.Marshal(cache)
	if err != nil {
		glob.GLog.Error("json 序列化异常", zap.Any("err", err))
	}

	glob.GRedis.HSet(context.Background(), "calc_cache", strconv.Itoa(int(calcRule.ID)), jsonData)

}

// Stop 根据传入的id停止计算任务
//
// 参数：
// id：计算规则id
//
// 返回值：
// bool：停止计算任务是否成功，成功返回true，否则返回false
func (b CalcRunBiz) Stop(id any) bool {
	var calcRule models.CalcRule

	result := glob.GDb.First(&calcRule, id)
	if result.Error != nil {
		glob.GLog.Error("查询异常", zap.Any("err", result.Error))
		panic(result.Error)
	}

	ma := make(map[string]interface{})
	ma["id"] = calcRule.ID

	jsonData, err := json.Marshal(ma)
	if err != nil {
		glob.GLog.Sugar().Errorf("Error marshalling struct to JSON: %s", err)
	}

	glob.GRedis.ZRem(context.Background(), "calc_queue", jsonData)
	glob.GRedis.HDel(context.Background(), "calc_cache", strconv.Itoa(int(calcRule.ID)))
	glob.GRedis.Del(context.Background(), "calc_queue_param:"+strconv.Itoa(int(calcRule.ID)))

	mup := map[string]interface{}{
		"Start": false, // 更新 Start 字段为 true
	}
	result = glob.GDb.Table("calc_rules").Where("id = ?", calcRule.ID).Updates(mup)

	if result.Error != nil {
		zap.S().Errorf("update error %+v", result.Error)
		return false
	}
	return true
}

// getNextTime 函数接收一个cron表达式字符串作为参数，返回下一次执行时间的Unix时间戳（秒）。
// 如果解析cron表达式失败，则返回0。
func getNextTime(cronExpr string) int64 {

	// 解析cron表达式
	schedule, err := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).Parse(cronExpr)

	if err != nil {
		fmt.Println("解析cron表达式失败:", err)
		return 0
	}

	// 获取当前时间
	now := time.Now()

	// 计算下一次执行时间
	nextTime := schedule.Next(now)

	// 将下一次执行时间格式化为字符串
	format := nextTime.Format("2006-01-02 15:04:05")
	fmt.Println("下一次执行时间:", format)

	// 将下一次执行时间转换为时间戳（秒）
	nextTimestamp := nextTime.Unix()
	return nextTimestamp
}

// MockCalc 是一个CalcRunBiz类型的方法，用于模拟计算操作
//
// 参数：
// start_time：int64类型，表示查询的起始时间戳
// end_time：int64类型，表示查询的结束时间戳
// id：string类型，表示计算规则的唯一标识
//
// 返回值：
// map[string]interface{}类型，表示计算的结果
func (b CalcRunBiz) MockCalc(startTime, endTime int64, id int) map[string]interface{} {

	var ccc servlet.CalcCache
	result, err := glob.GRedis.HGet(context.Background(), "calc_cache", strconv.Itoa(id)).Result()
	if err != nil {
		zap.S().Errorf("转化异常 %+v", err)
		return nil
	}
	err = json.Unmarshal([]byte(result), &ccc)
	if err != nil {
		zap.S().Infof("Failed to unmarshal message: %s", err)
		return nil
	}
	var m = make(map[string]any)
	for _, cache := range ccc.Param {

		if "原始" == cache.Reduce {
			var fd []string
			fd = append(fd, strconv.Itoa(cache.SignalId))
			config := servlet.InfluxQueryConfig{}
			config.Bucket = glob.GConfig.InfluxConfig.Bucket
			config.Measurement = strconv.Itoa(cache.MqttClientId)
			config.Fields = fd
			config.Aggregation = servlet.AggregationConfig{
				Every:       1,
				Function:    "mean",
				CreateEmpty: false,
			}
			config.StartTime = startTime
			config.EndTime = endTime
			config.Reduce = cache.Reduce
			query := config.GenerateFluxQuery()
			zap.S().Infof("influxdb query line = %s", query)

			result, err := glob.GInfluxdb.QueryAPI(glob.GConfig.InfluxConfig.Org).Query(context.Background(), query)
			if err != nil {
				zap.S().Error(err)
				return nil
			}
			v := make(map[int64]float64)
			for result.Next() {
				if result.TableChanged() {
					fmt.Printf("table: %s\n", result.TableMetadata().String())
				}
				values := result.Record().Values()
				fmt.Printf("value: %v\n", values)
				t := values["_time"].(time.Time)
				println(t.Unix())
				v[t.Unix()] = values["_value"].(float64)

			}
			m[cache.Name] = v
		} else {

			var fd []string
			fd = append(fd, strconv.Itoa(cache.SignalId))

			config := servlet.InfluxQueryConfig{}
			config.Bucket = glob.GConfig.InfluxConfig.Bucket
			config.Measurement = strconv.Itoa(cache.MqttClientId)
			config.Fields = fd
			config.StartTime = startTime
			config.EndTime = endTime
			config.Reduce = cache.Reduce
			query := config.GenerateFluxReduce()

			zap.S().Infof("influxdb query line = %s", query)
			result, err := glob.GInfluxdb.QueryAPI(glob.GConfig.InfluxConfig.Org).Query(context.Background(), query)
			if err != nil {
				zap.S().Error(err)
				return nil
			}

			for result.Next() {
				if result.TableChanged() {
					fmt.Printf("table: %s\n", result.TableMetadata().String())
				}
				values := result.Record().Values()
				fmt.Printf("value: %v\n", values)
				m[cache.Name] = values["_value"].(float64)

			}
		}
	}
	scriot := runCalcScriot(m, ccc.Script)

	var old models.CalcRule
	_ = glob.GDb.First(&old, ccc.ID)

	var newV models.CalcRule
	newV = old
	marshal, _ := json.Marshal(scriot)
	newV.MockValue = string(marshal)
	glob.GDb.Model(&newV).Updates(newV)

	return scriot
}

// runCalcScriot 函数执行传入的 JavaScript 脚本，并将计算结果以 map[string]interface{} 的形式返回
//
// 参数：
// param: 类型为 map[string]float64，表示计算参数
// script: 类型为 string，表示待执行的 JavaScript 脚本
//
// 返回值：
// 类型为 map[string]interface{}，表示执行 JavaScript 脚本后的计算结果
func runCalcScriot(param map[string]any, script string) map[string]interface{} {
	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Errorf("JS代码有问题 %+v", err)

		return nil
	}
	var fn func(string2 map[string]any) map[string]interface{}
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		zap.S().Errorf("Js函数映射到 Go 函数失败！ %+v", err)
		return nil
	}
	a := fn(param)
	return a
}

// QueryRuleExData 函数用于查询指定规则ID在指定时间范围内的扩展数据
//
// 参数：
// rule_id: 规则ID，类型为int64
// start_time: 查询开始时间，类型为int64
// end_time: 查询结束时间，类型为int64
//
// 返回值：
// 返回查询结果，类型为[]bson.M，即bson.M类型的切片
func (b CalcRunBiz) QueryRuleExData(ruleId, startTime, endTime int64) []bson.M {
	database := glob.GMongoClient.Database(glob.GConfig.MongoConfig.Db)
	collection := database.Collection(glob.GConfig.MongoConfig.Collection)

	filter := bson.M{
		"calc_rule_id": ruleId,
		"ex_time": bson.M{
			"$gte": startTime, // 大于等于开始时间
			"$lte": endTime,   // 小于等于结束时间
		},
	}
	// 创建FindOptions实例
	findOptions := options.Find()

	// 指定按照时间字段倒序排序
	// 假设时间字段的名称为"ex_time"
	findOptions.SetSort(bson.M{"ex_time": -1})

	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			zap.S().Errorf("err %+v", err)

		}
	}(cur, context.TODO())
	var c []bson.M

	for cur.Next(context.TODO()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		c = append(c, result)
	}

	return c

}
