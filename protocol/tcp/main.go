package main

import (
	"flag"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)
var globalConfig ServerConfig

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "app-local.yml", "Path to the config file")
	flag.Parse()

	yfile, err := os.ReadFile(configPath)
	if err != nil {
		zap.S().Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(yfile, &globalConfig)
	if err != nil {
		zap.S().Fatalf("error: %v", err)
	}

	zap.S().Infof("node name = %v , host = %v , port = %v", globalConfig.NodeInfo.Name, globalConfig.NodeInfo.Host, globalConfig.NodeInfo.Port)
	InitRabbitCon()

	initGlobalRedisClient(globalConfig.RedisConfig)
	server := New(&Config{
		Host: "localhost",
		Port: strconv.Itoa(globalConfig.NodeInfo.Port),
	})
	server.Run()
}




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