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

package task

import (
	"context"
	"encoding/json"
	"igp/biz"
	"igp/glob"
	"igp/initialize"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// handlerCalcQueue 函数处理计算任务队列
//
// 此函数创建一个带有秒级精度的新cron调度器，
// 每隔10秒执行一次handler_calc_task函数，
// 并将任务添加到调度器中。
// 如果添加任务失败，则记录错误并退出程序。
//
// 最后，启动cron调度器，开始执行任务。
func handlerCalcQueue() {
	cronIns := cron.New(cron.WithSeconds())

	_, err := cronIns.AddFunc("*/3 * * * * *", handlerCalcTask)
	if err != nil {
		log.Fatal(err)
	}

	// 启动 cron 调度器
	cronIns.Start()

}

var calcRunBiz = biz.CalcRunBiz{}

// handlerCalcTask 函数处理计算任务
func handlerCalcTask() {
	glob.GLog.Debug("handler_calc_task 开始执行")
	lock := initialize.NewRedisDistLock(glob.GRedis, "handlerCalcQueue")
	if lock.TryLock() {
		members, err := glob.GRedis.ZRevRangeWithScores(context.Background(), "calc_queue", 0, 99).Result()
		if err != nil {
			glob.GLog.Sugar().Error("Error getting ZSet members: %v\n", err)
			return
		}

		// 当前时间
		nowTime := time.Now().Unix()
		for _, member := range members {
			// 判断是否过期 ，过期了则表示要执行计算
			if member.Score < float64(nowTime) {

				var myMap map[string]uint
				bytes := []byte(member.Member.(string))
				err = json.Unmarshal(bytes, &myMap)
				if err != nil {
					log.Fatal(err)
				}
				glob.GRedis.Set(context.Background(), "calc_queue_param:"+strconv.Itoa(int(myMap["id"])), member.Score, 0)

				// 发送给消息队列
				pushToQueue("calc_queue", bytes)
				// 启动下一轮的计算任务
				calcRunBiz.Start(myMap["id"])
			}

		}
		lock.Unlock()

		return
	} else {
		glob.GLog.Sugar().Errorf("没有获取到handlerCalcQueue锁")
		return
	}
}

func pushToQueue(queueName string, body []byte) {

	ch, _ := glob.GRabbitMq.Channel()
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			zap.S().Errorf("channel close error: %+v", err)
		}
	}(ch)

	// 这里需要对每个消息做延迟处理，而不是对队列做延迟
	_ = ch.PublishWithContext(context.Background(), "", queueName, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			//Expiration: "过期时间",
			ContentType: "text/plain",
			Body:        body,
		})
	glob.GLog.Sugar().Infof(" [x] 发送到 %s 消息体 %s", queueName, body)

}

func InitTask() {
	zap.S().Info("初始化任务")
	handlerCalcQueue()
}
