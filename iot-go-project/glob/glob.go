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
