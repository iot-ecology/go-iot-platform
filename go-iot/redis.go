package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var globalRedisClient *redis.Client

func initGlobalRedisClient(config RedisConfig) {

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
