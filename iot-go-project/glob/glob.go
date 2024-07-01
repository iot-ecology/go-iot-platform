package glob

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"igp/config"
)

var (
	GDb          *gorm.DB
	GLog         *zap.Logger
	GRedis       *redis.Client
	GInfluxdb    influxdb2.Client
	GRabbitMq    *amqp.Connection
	GConfig      *config.ServerConfig
	GMongoClient *mongo.Client
)

type MessageType int

var (
	// StartNotification 计划开始通知
	StartNotification MessageType = 1
	// DueSoonNotification 计划临期通知
	DueSoonNotification MessageType = 2
	// DueNotification 计划到期通知
	DueNotification MessageType = 3
	// ProductionStartNotification 生产开始通知
	ProductionStartNotification MessageType = 4
	// ProductionCompleteNotification 生产完成通知
	ProductionCompleteNotification MessageType = 5
	// MaintenanceNotification 维修通知
	MaintenanceNotification MessageType = 6
	// MaintenanceStartNotification 维修开始通知
	MaintenanceStartNotification MessageType = 7
	// MaintenanceEndNotification 维修结束通知
	MaintenanceEndNotification MessageType = 8
)

func (mt *MessageType) String() string {
	switch *mt {
	case StartNotification:
		return "Plan_Start_Notification"
	case DueSoonNotification:
		return "Plan_Due_Soon_Notification"
	case DueNotification:
		return "Plan_Due_Notification"
	case ProductionStartNotification:
		return "Production_Start_Notification"
	case ProductionCompleteNotification:
		return "Production_Complete_Notification"
	case MaintenanceNotification:
		return "Maintenance_Notification"
	case MaintenanceStartNotification:
		return "Maintenance_Start_Notification"
	case MaintenanceEndNotification:
		return "Maintenance_End_Notification"

	default:
		return "未知消息类型"
	}
}
