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
