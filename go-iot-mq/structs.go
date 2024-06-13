package main

import (
	"fmt"
	"strings"
)

type MQTTMessage struct {
	MQTTClientID string `json:"mqtt_client_id"`
	Message      string `json:"message"`
}

type DataRowList struct {
	Time      int64     `json:"time"` // 秒级时间戳
	DeviceUid string    `json:"device_uid"`
	DataRows  []DataRow `json:"data"`
	Nc        string    `json:"nc"`
}
type DataRow struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
}

type ServerConfig struct {
	RedisConfig  RedisConfig  `yaml:"redis_config" json:"redis_config"`
	MQConfig     MQConfig     `yaml:"mq_config" json:"mq_config"`
	InfluxConfig InfluxConfig `yaml:"influx_config" json:"influx_config"`
	MongoConfig  MongoConfig  `yaml:"mongo_config" json:"mongo_config"`
	NodeInfo     NodeInfo     `yaml:"node_info" json:"node_info"`
}

type RedisConfig struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	Db       int    `json:"db,omitempty" yaml:"db,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

type InfluxConfig struct {
	Host   string `json:"host,omitempty" yaml:"host,omitempty"`
	Port   int    `json:"port,omitempty" yaml:"port,omitempty"`
	Token  string `json:"token,omitempty" yaml:"token,omitempty"`
	Org    string `json:"org,omitempty" yaml:"org,omitempty"`
	Bucket string `json:"bucket,omitempty" yaml:"bucket,omitempty"`
}
type MQConfig struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}
type MongoConfig struct {
	Host                   string `json:"host,omitempty" yaml:"host,omitempty"`
	Port                   int    `json:"port,omitempty" yaml:"port,omitempty"`
	Username               string `json:"username,omitempty" yaml:"username,omitempty"`
	Password               string `json:"password,omitempty" yaml:"password,omitempty"`
	Db                     string `json:"db,omitempty" yaml:"db,omitempty"`
	Collection             string `json:"collection,omitempty" yaml:"collection,omitempty"`
	WaringCollection       string `json:"waring_collection,omitempty" yaml:"waring_collection,omitempty"`
	ScriptWaringCollection string `json:"script_waring_collection,omitempty" yaml:"script_waring_collection,omitempty"`
}

type Signal struct {
	MqttClientId int    `json:"mqtt_client_id"` // MQTT客户端表的外键ID
	Name         string `json:"name"`           // 信号的名称，用于标识不同的信号
	Type         string `json:"type"`           // 信号的数据类型，如整数、字符串等
	ID           int    `json:"ID"`
	CacheSize    int64  `json:"cache_size"` // 缓存大小
}

type SignalWaringConfig struct {
	SignalId int     `json:"signal_id"` // 信号表的外键ID
	Min      float64 `json:"min"`       // 范围,小值
	Max      float64 `json:"max"`       // 范围,大值
	InOrOut  int     `json:"in_or_out"` //  1 范围内报警 0 范围外报警
	ID       int     `json:"ID"`
}

type SignalDelayWaringParam struct {
	MqttClientName      string `gorm:"-" json:"mqtt_client_name"`                             // MQTT客户端的名称，不存储在数据库中
	MqttClientId        int    `json:"mqtt_client_id"`                                        // MQTT客户端表的外键ID
	Name                string `json:"name"`                                                  // 参数名称
	SignalName          string `gorm:"signal_name"  json:"signal_name" structs:"signal_name"` // 信号表 name
	SignalId            int    `gorm:"signal_id"  json:"signal_id" structs:"signal_id"`       // 信号表的外键ID
	SignalDelayWaringId int    `json:"signal_delay_waring_id"`                                // SignalDelayWaring 主键
	ID                  int    `json:"ID"`
}
type SignalDelayWaring struct {
	Name   string `json:"name" structs:"name"`     // 名称
	Script string `json:"script" gorm:"type:text"` // 数据处理脚本
	ID     int    `json:"ID"`
}
type Tv struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

type CalcCache struct {
	ID     uint             `json:"id"`
	Param  []CalcParamCache `json:"param"`
	Cron   string           `json:"cron"`
	Script string           `json:"script"`
	Offset int64            `json:"offset"`
}

type CalcParamCache struct {
	MqttClientId int    `json:"mqtt_client_id"`                                        // MQTT客户端表的外键ID
	Name         string `json:"name"`                                                  // 参数名称
	SignalName   string `gorm:"signal_name"  json:"signal_name" structs:"signal_name"` // 信号表 name
	SignalId     int    `json:"signal_id" structs:"signal_id"`                         // 信号表的外键ID
	Reduce       string `json:"reduce"`                                                // 数据聚合方式 1. 求和 2. 平均值 3. 最大值 4. 最小值
	CalcRuleId   int    `json:"calc_rule_id"`                                          // CalcRule 主键
}

type InfluxQueryConfig struct {
	Bucket      string            `json:"-"`
	Measurement string            `json:"measurement"`
	Fields      []string          `json:"fields"`
	StartTime   int64             `json:"start_time"`
	EndTime     int64             `json:"end_time"`
	Aggregation AggregationConfig `json:"aggregation"`
	Reduce      string            `json:"reduce"` // sum min max mean
}

type AggregationConfig struct {
	Every       int    `json:"every"`
	Function    string `json:"function"` // 可能的值："mean", "sum", "min", "max"
	CreateEmpty bool   `json:"create_empty"`
}

// GenerateFluxQuery 根据InfluxQueryConfig生成Flux语言的查询语句
//
// 参数：
//
//	iqc *InfluxQueryConfig - InfluxQueryConfig类型指针，表示查询配置信息
//
// 返回值：
//
//	string - 生成的Flux语言的查询语句
func (iqc *InfluxQueryConfig) GenerateFluxQuery() string {
	var filters []string
	for _, field := range iqc.Fields {
		filters = append(filters, fmt.Sprintf(`r["_field"] == "%s"`, field))
	}

	timeRange := ""
	if iqc.StartTime != 0 && iqc.EndTime != 0 {
		timeRange = fmt.Sprintf(`start: %d, stop: %d`, iqc.StartTime, iqc.EndTime)
	}

	filterClause := ""
	if len(filters) > 0 {
		filterClause = fmt.Sprintf("|> filter(fn: (r) => %s)", strings.Join(filters, " or "))
	}

	return fmt.Sprintf(`
		from(bucket: "%s")
			|> range(%s)
			%s
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> aggregateWindow(every: %ds, fn: %s, createEmpty: %t)
			|> yield(name: "mean")
	`, iqc.Bucket, timeRange, filterClause, iqc.Measurement, iqc.Aggregation.Every, iqc.Aggregation.Function, iqc.Aggregation.CreateEmpty)
}

// GenerateFluxReduce 根据InfluxQueryConfig生成Flux语言的reduce查询语句
//
// 参数：
//
//	iqc *InfluxQueryConfig - InfluxQueryConfig类型指针，表示查询配置信息
//
// 返回值：
//
//	string - 生成的Flux语言的reduce查询语句
func (iqc *InfluxQueryConfig) GenerateFluxReduce() string {
	var filters []string
	for _, field := range iqc.Fields {
		filters = append(filters, fmt.Sprintf(`r["_field"] == "%s"`, field))
	}

	timeRange := ""
	if iqc.StartTime != 0 && iqc.EndTime != 0 {
		timeRange = fmt.Sprintf(`start: %d, stop: %d`, iqc.StartTime, iqc.EndTime)
	}

	filterClause := ""
	if len(filters) > 0 {
		filterClause = fmt.Sprintf("|> filter(fn: (r) => %s)", strings.Join(filters, " or "))
	}

	return fmt.Sprintf(`
		from(bucket: "%s")
			|> range(%s)
			%s
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> %s()
	`, iqc.Bucket, timeRange, filterClause, iqc.Measurement, iqc.Reduce)
}
