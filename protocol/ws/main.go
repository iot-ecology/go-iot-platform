package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"os"
	"syscall"
	"time"
)

var globalConfig ServerConfig

func main() {
	initLog()
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
	go BeatTask(globalConfig.NodeInfo)
	e := gin.Default()
	//html页面位置
	e.LoadHTMLGlob("template/*")
	//静态文件位置
	e.Static("/static", "./static")

	e.GET("/", Index)
	e.GET("/auth", AuthCtr)
	e.GET("/ws", InitWebSocket)

	go ListenerWs()
	e.Run(fmt.Sprintf(":%d", globalConfig.NodeInfo.Port))
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
