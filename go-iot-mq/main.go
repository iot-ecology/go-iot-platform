package main

import (
	"context"
	"flag"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

var globalConfig ServerConfig
var writeAPI api.WriteAPI

func main() {

	var configPath string
	flag.StringVar(&configPath, "config", "app-node1.yml", "Path to the config file")
	flag.Parse()

	yfile, err := os.ReadFile(configPath)
	if err != nil {
		zap.S().Fatalf("error: %+v", err)
	}

	err = yaml.Unmarshal(yfile, &globalConfig)
	if err != nil {
		zap.S().Fatalf("error: %+v", err)
	}
	InitLog()

	InitGlobalRedisClient(globalConfig.RedisConfig)
	InitInfluxDbClient(globalConfig.InfluxConfig)
	writeAPI = GlobalInfluxDbClient.WriteAPI(globalConfig.InfluxConfig.Org, globalConfig.InfluxConfig.Bucket)
	//InitRabbitCon(globalConfig.MQConfig)
	err = ConnectToRMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	zap.S().Infof("消息队列类型 %s", globalConfig.NodeInfo.Type)

	CreateRabbitQueue("waring_handler")
	CreateRabbitQueue("waring_delay_handler")
	initMongo()
	failOnError(err, "Failed to open a channel")
	go startHttp()
	cus := NewConsumer("", "amqp://guest:guest@localhost:5672", "", "", "")
	err = cus.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	if globalConfig.NodeInfo.Type == "pre_handler" {
		deliveries, err := cus.AnnounceQueue("pre_handler", "")
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		}
		cus.Handle(deliveries, HandlerDataStorage, 1, "pre_handler", "")
	}
	if globalConfig.NodeInfo.Type == "waring_handler" {
		waring_handler, err := cus.AnnounceQueue("waring_handler", "")
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		}
		cus.Handle(waring_handler, HandlerWaring, 1, "waring_handler", "")
	}
	if globalConfig.NodeInfo.Type == "calc_queue" {
		calc_queue, err := cus.AnnounceQueue("calc_queue", "")
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		}
		cus.Handle(calc_queue, HandlerCalc, 1, "calc_queue", "")
	}
	if globalConfig.NodeInfo.Type == "waring_delay_handler" {
		waring_delay_handler, err := cus.AnnounceQueue("waring_delay_handler", "")
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		}
		cus.Handle(waring_delay_handler, HandlerWaringDelay, 1, "waring_delay_handler", "")
	}

}

func startHttp() {
	http.HandleFunc("/beat", HttpBeat)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":"+strconv.Itoa(globalConfig.NodeInfo.Port), nil); err != nil {
		zap.S().Fatalf("Failed to start server: %s", err)
	}
	zap.S().Infof("Server started Port %d", globalConfig.NodeInfo.Port)
}

func HttpBeat(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return
	}
	// 向客户端发送响应消息
	_, _ = fmt.Fprintf(w, "ok")
}

func failOnError(err error, msg string) {
	if err != nil {
		zap.S().Panicf("%s: %s", msg, err)
	}
}

var GlobalInfluxDbClient influxdb2.Client

// InitInfluxDbClient 函数用于初始化InfluxDB客户端
//
// 参数：
// config InfluxConfig - InfluxDB配置信息
//
// 返回值：
// 无
func InitInfluxDbClient(config InfluxConfig) {
	GlobalInfluxDbClient = influxdb2.NewClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port), config.Token)

}

// PushToQueue 将消息推送到RabbitMQ队列中
//
// 参数：
// queueName string - 目标队列名称
// body []byte - 消息体
//
// 返回值：
// 无
func PushToQueue(queueName string, body []byte) {

	zap.S().Infof("开始推送消息到队列 %s msg %s", queueName, body)

	err := chann.PublishWithContext(context.Background(), "", queueName, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to open a channel")

}

var GRabbitMq *amqp.Connection

// InitRabbitCon 初始化RabbitMQ连接
func InitRabbitCon(config MQConfig) {
	zap.S().Infof("开始处理rabbitmq")
	conn, err := amqp.Dial(genUrl(config))
	if err != nil {
		zap.S().Fatalf("Failed to connect to RabbitMQ  %+v", err)
	}

	GRabbitMq = conn
	zap.S().Infof("Connected to RabbitMQ")

}

// genUrl 生成RabbitMQ连接字符串
//
// 返回值为一个字符串，格式为amqp://用户名:密码@主机名:端口号/
//
// 参数说明：
//
//	无
//
// 返回值说明：
//
//	返回RabbitMQ的连接字符串
func genUrl(config MQConfig) string {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.Username, config.Password, config.Host, config.Port)
	return connStr
}

// CreateRabbitQueue 创建一个RabbitMQ队列
//
// 参数：
// queueName string - 队列名称
//
// 返回值：
// 无
func CreateRabbitQueue(queueName string) {

	_, err := chann.QueueDeclare(queueName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		zap.S().Fatalf("创建queue异常 %s", queueName)
	}
}

var GMongoClient *mongo.Client

// initMongo 函数用于初始化 MongoDB 连接
func initMongo() {
	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%d", globalConfig.MongoConfig.Username, globalConfig.MongoConfig.Password, globalConfig.MongoConfig.Host, globalConfig.MongoConfig.Port)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	GMongoClient = client

	zap.S().Infof("Connected to MongoDB!")

}
