package mongo

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 全局map，用于存储id和MongoDB客户端的映射
var mongoClientMap = make(map[uint]*mongo.Client)

// 用于同步的互斥锁
var mu sync.Mutex

// GetMongoDBClient 根据提供的参数获取或创建MongoDB客户端
func GetMongoDBClient(Host, Username, Password, Db string, Port int, id uint) (*mongo.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	// 使用id作为map的key
	key := id

	client, exists := mongoClientMap[key]
	if exists {
		zap.S().Infof("Reusing existing MongoDB client for id: %d.", id)
		return client, nil
	}

	encodedUsername := url.QueryEscape(Username)
	encodedPassword := url.QueryEscape(Password)

	// 构建MongoDB连接字符串
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", encodedUsername, encodedPassword, Host, Port, Db)

	if Username == "admin" {
		connectionString = connectionString + "?authSource=admin"
	}
	// 创建MongoDB客户端
	clientOptions := options.Client().ApplyURI(connectionString)
	newClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// 检查连接是否成功
	err = newClient.Ping(context.TODO(), nil)
	if err != nil {
		zap.S().Fatalf("Failed to connect to MongoDB: %s", err)
		return nil, err
	}

	// 将新的MongoDB客户端存储到map中
	mongoClientMap[key] = newClient
	zap.S().Infof("New MongoDB client for database '%s' with id %d has been initiated.", Db, id)
	return newClient, nil
}
