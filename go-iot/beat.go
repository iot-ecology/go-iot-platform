package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// ServerConfig 定义了服务器配置的结构体，包含了节点信息、Redis配置和消息队列配置。
type ServerConfig struct {
	// NodeInfo 定义了节点的信息，包括主机地址、端口、节点名称、节点类型和最大处理数量。
	NodeInfo NodeInfo `yaml:"node_info" json:"node_info"`

	// RedisConfig 定义了Redis服务器的配置，包括主机地址、端口、数据库索引和密码。
	RedisConfig RedisConfig `yaml:"redis_config" json:"redis_config"`

	// MQConfig 定义了消息队列服务器的配置，包括主机地址、端口、用户名和密码。
	MQConfig MQConfig `yaml:"mq_config" json:"mq_config"`
}

// NodeInfo 定义了节点的基本信息。
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
}

// RedisConfig 定义了Redis服务器的配置信息。
type RedisConfig struct {
	// Host 表示Redis服务器的主机地址。
	Host string `json:"host,omitempty" yaml:"host,omitempty"`

	// Port 表示Redis服务器监听的端口号。
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// Db 表示Redis服务器的数据库索引。
	Db int `json:"db,omitempty" yaml:"db,omitempty"`

	// Password 表示Redis服务器的访问密码。
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

// MQConfig 定义了消息队列服务器的配置信息。
type MQConfig struct {
	// Host 表示消息队列服务器的主机地址。
	Host string `json:"host,omitempty" yaml:"host,omitempty"`

	// Port 表示消息队列服务器监听的端口号。
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// Username 表示用于访问消息队列服务器的用户名。
	Username string `json:"username,omitempty" yaml:"username,omitempty"`

	// Password 表示用于访问消息队列服务器的密码。
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

// ListenerBeat 监听Redis中以 "beat" 开头的键的过期事件，并处理这些事件。
func ListenerBeat() {
	client := globalRedisClient

	// 配置Redis以启用过期事件的通知。
	client.ConfigSet(context.Background(), "notify-keyspace-events", "Ex")

	// 订阅Redis的过期事件。
	pubsub := client.Subscribe(context.Background(), "__keyevent@"+strconv.Itoa(globalConfig.RedisConfig.Db)+"__:expired")
	// 确保订阅在函数结束时关闭。
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
		}
	}(pubsub)

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			// 如果接收消息时出现错误，记录错误并退出。
			zap.S().Fatalf("Error %s", err)
			return
		}

		// 如果消息负载以 "beat" 开头，则进一步处理。
		if strings.HasPrefix(msg.Payload, "beat") {
			// 分割负载字符串以获取最后一个元素。
			parts := strings.Split(msg.Payload, ":")
			lastElement := parts[len(parts)-1]

			// 调用 HandlerOffNode 处理过期的 "beat" 键。
			HandlerOffNode(lastElement)
		}
	}
}
