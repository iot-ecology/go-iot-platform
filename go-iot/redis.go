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
