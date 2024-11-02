package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var globalRedisClient *redis.Client

// InitGlobalRedisClient 初始化全局Redis客户端
//
// 参数:
// config RedisConfig - Redis配置信息
//
// 返回值:
// 无
func InitGlobalRedisClient(config RedisConfig) {

	add := fmt.Sprintf("%s:%d", config.Host, config.Port)
	globalRedisClient = redis.NewClient(&redis.Options{
		Addr:     add,
		Password: config.Password, // 如果没有设置密码，就留空字符串
		DB:       config.Db,       // 使用默认数据库
	})

	// 检查连接是否成功
	if err := globalRedisClient.Ping(context.Background()).Err(); err != nil {
		zap.S().Fatalf("Could not connect to Redis: %v", err)
	}

}

// GetScriptRedis 根据MQTT客户端ID从Redis中获取对应的脚本
// 参数:
//
//	mqttClientId string - MQTT客户端ID
//
// 返回值:
//
//	string - 对应的脚本
func GetScriptRedis(mqttClientId string) string {
	val := globalRedisClient.HGet(context.Background(), "mqtt_script", mqttClientId).Val()
	return val
}
