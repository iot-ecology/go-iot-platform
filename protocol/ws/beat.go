package main

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// registerInfo 注册节点信息
func registerInfo(node *NodeInfo) {
	zap.S().Infof("registerInfo 开始, node = %v", node)

	jsonData, _ := json.Marshal(node)
	globalRedisClient.Set(context.Background(), "pod:info:ws:"+node.Name, jsonData, 1*time.Hour)

}
func BeatTask(f NodeInfo) {
	zap.S().Infof("BeatTask 开始, f = %v", f)
	registerInfo(&f)

	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		registerInfo(&f)
	}
}
