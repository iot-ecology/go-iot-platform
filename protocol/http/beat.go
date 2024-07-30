package main

import (
	"context"
	"encoding/json"
	"time"
)

// registerInfo 注册节点信息
func registerInfo(node *NodeInfo) {

	jsonData, _ := json.Marshal(node)
	globalRedisClient.Set(context.Background(), "pod:info:http:"+node.Name, jsonData, 1*time.Hour)

}
func BeatTask(f NodeInfo) {
	registerInfo(&f)

	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		registerInfo(&f)
	}
}
