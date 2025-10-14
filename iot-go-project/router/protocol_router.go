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
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"igp/glob"
	"strings"
)


type ProtocolService struct{}


type NodeInfo struct {
	// Host 表示节点的主机地址。
	Host string `json:"host,omitempty" yaml:"host,omitempty"`

	// Port 表示节点监听的端口号。
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// Name 表示节点的名称。
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Type 表示节点的类型。
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// Size 表示节点可以处理的最大数量。
	Size int64 `json:"size,omitempty" yaml:"size,omitempty"`

	ClientInfo map[string]interface{} `json:"client_info,omitempty" yaml:"client_info,omitempty"`
}

// WsServerInfo
// @Summary WebSocket服务端信息
// @Description 协议服务信息
// @Tags protocol
// @Accept json
// @Produce json
// @Success 200 {object} servlet.JSONResult{data=router.NodeInfo[]} "节点信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /protocol/ws_info [get]
func (api *ProtocolService) WsServerInfo(c *gin.Context){
	val := glob.GRedis.Keys(context.Background(), "pod:info:ws:*").Val()

	var  res []NodeInfo
	for _, s := range val {


		s2 := glob.GRedis.Get(context.Background(), s).Val()
		v := NodeInfo{}
		err := json.Unmarshal([]byte(s2), &v)


		if err != nil {
			zap.S().Error(err)
			continue
		}
		strings := glob.GRedis.LRange(context.Background(), "ws_uid:" + v.Name, 0, -1).Val()
		zap.S().Info(strings)

		clientInfo := make(map[string]interface{})
		for _, s3 := range strings {
			result, err := glob.GRedis.Get(context.Background(), "ws:last:"+s3).Result()
			if err !=nil {
				zap.S().Error(err)
				continue
			}
			clientInfo[s3] = result
		}
		v.ClientInfo = clientInfo
		res = append(res, v)

	}
	c.JSON(200, res)
}




// TcpServerInfo
// @Summary Tcp服务端信息
// @Description 协议服务信息
// @Tags protocol
// @Accept json
// @Produce json
// @Success 200 {object} servlet.JSONResult{data=router.NodeInfo[]} "节点信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /protocol/tcp_info [get]
func (api *ProtocolService) TcpServerInfo(c *gin.Context){
	val := glob.GRedis.Keys(context.Background(), "pod:info:tcp:*").Val()

	var  res []NodeInfo
	for _, s := range val {




		s2 := glob.GRedis.Get(context.Background(), s).Val()
		v := NodeInfo{}
		err := json.Unmarshal([]byte(s2), &v)


		if err != nil {
			zap.S().Error(err)
			continue
		}

		parts := strings.Split(s, ":")
		lastElement := parts[len(parts)-1]

		m := glob.GRedis.HRandField(context.Background(), "tcp_uid_f:"+lastElement,-1).Val()
		zap.S().Info(m)

		clientInfo := make(map[string]interface{})
		for _, s3 := range m {
			replace := strings.Replace(s3, ":", "@", -1)
			result, err := glob.GRedis.Get(context.Background(), "tcp:last:"+replace).Result()
			if err !=nil {
				zap.S().Error(err)
				continue
			}
			clientInfo[s3] = result
		}

		v.ClientInfo = clientInfo
		res = append(res, v)

	}
	c.JSON(200, res)
}




// CoapServerInfo
// @Summary Coap服务端信息
// @Description 协议服务信息
// @Tags protocol
// @Accept json
// @Produce json
// @Success 200 {object} servlet.JSONResult{data=router.NodeInfo[]} "节点信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /protocol/coap_info [get]
func (api *ProtocolService) CoapServerInfo(c *gin.Context){
	val := glob.GRedis.Keys(context.Background(), "pod:info:coap:*").Val()

	var  res []NodeInfo
	for _, s := range val {




		s2 := glob.GRedis.Get(context.Background(), s).Val()
		v := NodeInfo{}
		err := json.Unmarshal([]byte(s2), &v)


		if err != nil {
			zap.S().Error(err)
			continue
		}

		parts := strings.Split(s, ":")
		lastElement := parts[len(parts)-1]

		m := glob.GRedis.HRandField(context.Background(), "coap_uid_f:"+lastElement,-1).Val()
		zap.S().Info(m)

		clientInfo := make(map[string]interface{})
		for _, s3 := range m {
			replace := strings.Replace(s3, ":", "@", -1)
			result, err := glob.GRedis.Get(context.Background(), "coap:last:"+replace).Result()
			if err !=nil {
				zap.S().Error(err)
				continue
			}
			clientInfo[s3] = result
		}

		v.ClientInfo = clientInfo
		res = append(res, v)

	}
	c.JSON(200, res)
}