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
	"errors"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
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
	initLog()

	err = yaml.Unmarshal(yfile, &globalConfig)
	if err != nil {
		zap.S().Fatalf("error: %v", err)
	}

	zap.S().Infof("node name = %v , host = %v , port = %v", globalConfig.NodeInfo.Name, globalConfig.NodeInfo.Host, globalConfig.NodeInfo.Port)
	InitRabbitCon()
	initGlobalRedisClient(globalConfig.RedisConfig)

	globalRedisClient.Del(context.Background(), "coap_uid_f:"+globalConfig.NodeInfo.Name)
	globalRedisClient.Del(context.Background(), "coap_uid:"+globalConfig.NodeInfo.Name)
	go BeatTask(globalConfig.NodeInfo)
	go ListenerCoap()
	go startHttp()

	Create(globalConfig.NodeInfo.Port)
}



func startHttp() {
	http.HandleFunc("/beat", HttpBeat)
	http.Handle("/metrics", promhttp.Handler())

	httpPort := globalConfig.NodeInfo.Port + 10000
	zap.S().Infof("HTTP Server started Port %d", httpPort)

	if err := http.ListenAndServe(":"+strconv.Itoa(httpPort), nil); err != nil {
		zap.S().Fatalf("Failed to start server: %s", err)
	}
}

func HttpBeat(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return
	}
	// 向客户端发送响应消息
	_, _ = fmt.Fprintf(w, "ok")
}

var myTimeEncoder = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 按照 "2006-01-02 15:04:05" 的格式编码时间
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
})

func initLog() {
	encoderConfig := zapcore.EncoderConfig{
		// 使用自定义的时间编码器
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码日志级别
		EncodeTime:     myTimeEncoder,                 // 使用自定义的时间编码器
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码调用者
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), // 使用 Console 编码器
		zapcore.AddSync(os.Stdout),          // 输出到标准输出
		zap.NewAtomicLevelAt(zap.InfoLevel), // 设置日志级别为 Debug
	)

	lg := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // 替换全局 Logger

	// 确保日志被刷新
	defer func(lg *zap.Logger) {
		err := lg.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) {
			zap.S().Errorf("日志同步失败 %+v", err)
		}
	}(lg)

	// 记录一条日志作为示例
	lg.Debug("这是一个调试级别的日志")
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
