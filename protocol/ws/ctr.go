package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
func AuthCtr(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "请求头 authorization 丢失"})
		return
	}
	user, pass, ok := parseBasicAuth(authHeader)
	zap.S().Infof("解析的账号密码 username: %s, password: %s", user, pass)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "请求头 authorization 无法正常解析"})
		return
	}
	deviceId := c.Request.Header.Get("device_id")

	username, password := FindDeviceMappingUP(deviceId)
	zap.S().Infof("device_id: %s", deviceId)
	zap.S().Infof("有效账号密码 username: %s, password: %s", username, password)
	if username != user || password != pass {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "账号密码不匹配"})
		return
	}

	// fixme: 确认客户端数量

	curSize := globalRedisClient.LLen(context.Background(), "ws_uid:"+globalConfig.NodeInfo.Name).Val()
	if curSize <= globalConfig.NodeInfo.Size {

		newUUID, _ := uuid.NewUUID()

		uidString := newUUID.String()
		storageId := deviceId + "@" + uidString
		globalRedisClient.LPush(context.Background(), "ws_uid:"+globalConfig.NodeInfo.Name, storageId)

		c.JSON(http.StatusOK, gin.H{
			"message": "认证通过",
			"uid":     storageId,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "当前节点已满",
		})

	}
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func parseBasicAuth(authHeader string) (username, password string, ok bool) {
	zap.S().Infof("parseBasicAuth 开始, authHeader = %v", authHeader)
	// 基本认证格式："Basic <base64-encoded-string>"
	const prefix = "Basic "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return
	}

	// 解码base64
	enc := authHeader[len(prefix):]
	decoded, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return
	}

	// 解析用户名和密码
	cs := string(decoded)
	if !strings.Contains(cs, ":") {
		return
	}
	split := strings.Split(cs, ":")
	return split[0], split[1], true
}
func FindDeviceMappingUP(deviceId string) (string, string) {
	zap.S().Infof("FindDeviceMappingUP 开始, deviceId = %v", deviceId)
	// todo: 从redis中根据deviceId获取用户名和密码
	val := globalRedisClient.HGet(context.Background(), "auth:ws", deviceId).Val()
	var auth Auth
	err := json.Unmarshal([]byte(val), &auth)
	if err != nil {
		return "", ""
	}
	return auth.Username, auth.Password
}
