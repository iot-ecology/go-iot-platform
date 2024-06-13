package models

import (
	"gorm.io/gorm"
)

type MqttClient struct {
	Host     string `json:"host"`      // 主机
	Port     int    `json:"port"`      // 端口
	ClientId string `json:"client_id"` // 客户端id
	Username string `json:"username"`  // 账号
	Password string `json:"password"`  // 密码
	Subtopic string `json:"subtopic"`  // 订阅的主题
	Start    bool   `json:"start"`     // 是否启动

	Script     string `json:"script" gorm:"type:text"` // 数据处理脚本
	gorm.Model `structs:"-"`
}

type Signal struct {
	MqttClientId   int    `json:"mqtt_client_id"`                  // MQTT客户端表的外键ID
	Name           string `json:"name"`                            // 信号的名称，用于标识不同的信号
	Alias          string `json:"alias" structs:"alias"`           // 信号的别名，用于显示
	Type           string `json:"type"`                            // 信号的数据类型，如整数、字符串等
	MqttClientName string `gorm:"-" json:"mqtt_client_name"`       // MQTT客户端的名称，不存储在数据库中
	Unit           string `json:"unit" structs:"unit"`             // 单位
	CacheSize      int    `json:"cache_size" structs:"cache_size"` // 缓存大小
	gorm.Model     `structs:"-"`
}
type SignalWaringConfig struct {
	SignalId     int     `gorm:"signal_id"  json:"signal_id" structs:"signal_id"` // 信号表的外键ID
	Min          float64 `gorm:"min"  json:"min" structs:"min"`                   // 范围,小值
	Max          float64 `gorm:"max"  json:"max" structs:"max"`                   // 范围,大值
	InOrOut      int     `gorm:"in_or_out"  json:"in_or_out" structs:"in_or_out"` //范围内报警,范围外报警 1 范围内报警 0 范围外报警
	MqttClientId int     `json:"mqtt_client_id" structs:"mqtt_client_id"`         // MQTT客户端表的外键ID
	gorm.Model   `structs:"-"`
	Signal       Signal `gorm:"-" json:"signal" structs:"-"`
}

// SignalDelayWaring : 允许信号数据做一定条数的暂存
type SignalDelayWaring struct {
	Name       string `json:"name" structs:"name"`     // 名称
	Script     string `json:"script" gorm:"type:text"` // 数据处理脚本
	gorm.Model `structs:"-"`
}

type SignalDelayWaringParam struct {
	MqttClientName      string `gorm:"-" json:"mqtt_client_name"`                             // MQTT客户端的名称，不存储在数据库中
	MqttClientId        int    `json:"mqtt_client_id"`                                        // MQTT客户端表的外键ID
	Name                string `json:"name"`                                                  // 参数名称
	SignalName          string `gorm:"signal_name"  json:"signal_name" structs:"signal_name"` // 信号表 name
	SignalId            int    `gorm:"signal_id"  json:"signal_id" structs:"signal_id"`       // 信号表的外键ID
	SignalDelayWaringId int    `json:"signal_delay_waring_id"`                                // SignalDelayWaring 主键
	gorm.Model          `structs:"-"`
}

type NodeInfo struct {
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	Port int    `json:"port,omitempty" yaml:"port,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	Size int64  `json:"size,omitempty" yaml:"size,omitempty"` // 最大处理数量
}

// CalcRule 计算规则
type CalcRule struct {
	Name       string `json:"name" structs:"name"`                              // 任务名称
	Cron       string `json:"cron"`                                             // 定时任务执行周期
	Script     string `json:"script" gorm:"type:text"`                          // 数据处理脚本
	Offset     int64  `json:"offset" structs:"offset"`                          // 执行时间往前推移几秒
	Start      bool   `json:"start" structs:"start"`                            // 是否启动
	MockValue  string `json:"mock_value" structs:"mock_value" gorm:"type:text"` // 模拟数据
	gorm.Model `structs:"-"`
}

// CalcParam 计算参数
type CalcParam struct {
	MqttClientId   int    `json:"mqtt_client_id"`                                        // MQTT客户端表的外键ID
	Name           string `json:"name"`                                                  // 参数名称
	SignalName     string `gorm:"signal_name"  json:"signal_name" structs:"signal_name"` // 信号表 name
	SignalId       int    `gorm:"signal_id"  json:"signal_id" structs:"signal_id"`       // 信号表的外键ID
	Reduce         string `json:"reduce"`                                                // 数据聚合方式 1. mean 2. sum 3. max 4. min
	CalcRuleId     int    `json:"calc_rule_id"`                                          // CalcRule 主键
	MqttClientName string `gorm:"-" json:"mqtt_client_name"`                             // MQTT客户端的名称，不存储在数据库中
	gorm.Model     `structs:"-"`
}

// Dashboard 可视化板
type Dashboard struct {
	Name       string `json:"name"`                    // 名称
	Config     string `json:"config" gorm:"type:text"` // 配置
	gorm.Model `structs:"-"`
}
