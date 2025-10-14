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

package router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"igp/glob"
	"igp/models"
	"io"
	"net/http"
	"strings"
)


type PodInfo struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	NodeInfo NodeInfo `json:"node_info"`
}

type PodApi struct {
}
type GroupedPodInfo struct {
	Type    string    `json:"type"`
	PodInfo []PodInfo `json:"pod_info"`
}
func (api *PodApi) PodMqtt(c *gin.Context) {
	all, err := glob.GRedis.HGetAll(context.Background(), "register:mqtt").Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve data from Redis"})
		return
	}

	var res []models.NodeInfo
	for _, v := range all {
		var nodeInfo models.NodeInfo
		if err := json.Unmarshal([]byte(v), &nodeInfo); err != nil {
			c.JSON(500, gin.H{"error": "Failed to parse node info"})
			return
		}
		res = append(res, nodeInfo)
	}

	c.JSON(200, res)
}

func (api *PodApi) PodMetrics(c *gin.Context) {
	var pod PodInfo

	// 从请求中获取 PodInfo（假设是 JSON 格式的请求体）
	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var metricsURL string

	// 根据不同的类型构建 metrics URL
	if pod.Type == "tcp" || pod.Type == "coap" {
		metricsURL = fmt.Sprintf("http://%s:%d/metrics", pod.NodeInfo.Host, pod.NodeInfo.Port+10000)
	} else if pod.Type == "http" || pod.Type == "ws" {
		metricsURL = fmt.Sprintf("http://%s:%d/metrics", pod.NodeInfo.Host, pod.NodeInfo.Port)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported pod type"})
		return
	}

	// 发起请求获取 metrics
	resp, err := http.Get(metricsURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching metrics for %s: %s", pod.Name, err.Error())})
		return
	}
	defer resp.Body.Close()

	// 处理响应（假设返回的是文本）
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error reading metrics for %s: %s", pod.Name, err.Error())})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{"metrics": fmt.Sprintf("Metrics for %s: %s", pod.Name, string(bodyBytes))})
}


func (api *PodApi) PodInfo(c *gin.Context) {
	ctx := context.Background() // 确保上下文被创建
	keys, err := glob.GRedis.Keys(ctx, "pod:info:*").Result()
	if err != nil {
		zap.S().Errorf("Could not get keys: %v", err)
		c.JSON(500, gin.H{"error": "Could not get keys"})
		return
	}

	// 存储结果
	groupedResult := make(map[string][]PodInfo)

	for _, key := range keys {
		// 获取每个键的值
		value, err := glob.GRedis.Get(ctx, key).Result()
		if err != nil {
			zap.S().Errorf("Could not get value for key %s: %v", key, err)
			continue
		}

		// 假设 value 是一个 JSON 字符串，你可以根据实际情况解析
		var nodeInfo NodeInfo
		// 这里需要根据 value 的具体格式进行解析
		json.Unmarshal([]byte(value), &nodeInfo)

		// 从键中提取协议和名称（示例处理）
		parts := strings.Split(key, ":")
		if len(parts) < 4 {
			continue // 如果格式不对，跳过
		}
		protocol := parts[2] // 协议
		name := parts[3]     // 名称

		// 将 PodInfo 添加到相应的类型分组
		groupedResult[protocol] = append(groupedResult[protocol], PodInfo{
			Name: name,
			NodeInfo: NodeInfo{
				Host: nodeInfo.Host, // 从实际值中获取
				Port: nodeInfo.Port, // 从实际值中获取
				Name: nodeInfo.Name, // 从实际值中获取
			},
		})
	}

	// 转换为切片形式
	finalResult := make([]GroupedPodInfo, 0)
	for typeKey, podInfos := range groupedResult {
		finalResult = append(finalResult, GroupedPodInfo{
			Type:    typeKey,
			PodInfo: podInfos,
		})
	}


	// 返回结果
	c.JSON(200, finalResult)
}



