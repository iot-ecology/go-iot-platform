package biz

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"igp/glob"
	"igp/models"
	"strconv"
	"strings"
)

func InitRedisExpireHandler(client *redis.Client) {
	client.ConfigSet(context.Background(), "notify-keyspace-events", "Ex")

	// 订阅Redis的过期事件。
	pubsub := client.Subscribe(context.Background(), "__keyevent@"+strconv.Itoa(glob.GConfig.RedisConfig.Db)+"__:expired")
	// 确保订阅在函数结束时关闭。
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
	}(pubsub)

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			// 如果接收消息时出现错误，记录错误并退出。
			zap.S().Fatalf("Error %s", err)
			return
		}
		zap.S().Infof("Redis Expire key: %s", msg.Payload)

		parts := strings.Split(msg.Payload, ":")
		if len(parts) < 2 {
			continue
		}
		refId := parts[len(parts)-1]
		model := models.MessageList{
			MessageTypeId: int(glob.StartNotification),
			RefId:         refId,
		}
		if strings.HasPrefix(msg.Payload, glob.StartNotification.String()) {
			model.Content = "您有新的生产计划需要开始请注意!"
			model.EnContent = "You have a new production plan that needs to start, please take note!"
			glob.GDb.Model(&models.MessageList{}).Create(model)
		}
		if strings.HasPrefix(msg.Payload, glob.DueSoonNotification.String()) {
			model.Content = "您的生产计划马上到达截止时间请注意!"
			model.EnContent = "Your production plan is about to reach its deadline, please take note!"
			glob.GDb.Model(&models.MessageList{}).Create(model)
		}
		if strings.HasPrefix(msg.Payload, glob.DueNotification.String()) {
			model.Content = "您的生产计划按要求应当完成，请确认是否完成!"
			model.EnContent = "Your production plan should be completed as required, please confirm if it is done!"
			glob.GDb.Model(&models.MessageList{}).Create(model)
		}

		if strings.HasPrefix(msg.Payload, glob.MaintenanceNotification.String()) {
			model.Content = "您有新的维修任务，请查收!"
			model.EnContent = "You have a new maintenance task, please check it out!"
			glob.GDb.Model(&models.MessageList{}).Create(model)
		}
		if strings.HasPrefix(msg.Payload, glob.MaintenanceStartNotification.String()) {
			model.Content = "维修任务已开始."
			model.EnContent = "The maintenance task has started."
			glob.GDb.Model(&models.MessageList{}).Create(model)
		}
		if strings.HasPrefix(msg.Payload, glob.MaintenanceEndNotification.String()) {
			model.Content = "维修任务已完成."
			model.EnContent = "The maintenance task has been completed."
			glob.GDb.Model(&models.MessageList{}).Create(model)
		}

	}
}
