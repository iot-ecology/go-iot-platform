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

package servlet

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"igp/ut"
	"iot-transmit/common"
	"net/http"
	"strings"
	"time"
)

type GetTodoReq struct {
	Title string `json:"title"` // 主题
}

type JSONResult struct {
	Code    int         `json:"code" `
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationQ struct {
	Size  int         `form:"size" json:"size"`
	Page  int         `form:"page" json:"page"`
	Data  interface{} `json:"data" comment:"muster be a pointer of slice gorm.Model"`
	Total int64       `json:"total"`
}

type MqttScript struct {
	ID uint `json:"id"`

	Script string `json:"script"` // 数据处理脚本

}
type DataRowList struct {
	Time               int64     `json:"time"`                // 秒级时间戳
	DeviceUid          string    `json:"device_uid"`          // 能够产生网络通讯的唯一编码
	IdentificationCode string    `json:"identification_code"` // 设备标识码
	DataRows           []DataRow `json:"data"`
	Nc                 string    `json:"nc"`
}
type DataRow struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CheckScriptReq 用于封装POST请求的数据
type CheckScriptReq struct {
	Param  string `json:"param" binding:"required"`
	Script string `json:"script" binding:"required"`
}

func Resp(c *gin.Context, data interface{}) {
	result := JSONResult{}
	result.Message = "操作成功"
	result.Code = 20000
	result.Data = data
	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, result)
}

func Resp2(c *gin.Context, msg string) {
	result := JSONResult{}
	result.Message = msg
	result.Code = 20000
	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, result)
}
func Error(c *gin.Context, data interface{}) {
	result := JSONResult{}
	result.Message = "操作失败"
	result.Code = 40000
	result.Data = data
	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, result)
}

type ParamStruct struct {
	ClientID string `json:"client_id"`
	Topic    string `json:"topic"`
	QOS      byte   `json:"qos"`
	Retained bool   `json:"retained"`
	Payload  string `json:"payload"`
}

type InfluxResponse struct {
	Series []struct {
		Name    string            `json:"name"`
		Tags    map[string]string `json:"tags,omitempty"`
		Columns []string          `json:"columns"`
		Values  [][]interface{}   `json:"values"`
	} `json:"series"`
	Error string `json:"error,omitempty"`
}

type InfluxQueryConfig struct {
	Bucket      string            `json:"-"`
	DeviceUid   uint              `json:"device_uid,omitempty"`
	Protocol    string            `json:"protocol,omitempty"`
	Measurement string            `json:"measurement,omitempty"`
	Fields      []string          `json:"fields,omitempty"`
	StartTime   int64             `json:"start_time,omitempty"`
	EndTime     int64             `json:"end_time,omitempty"`
	Aggregation AggregationConfig `json:"aggregation,omitempty"`
	Reduce      string            `json:"reduce,omitempty"` // sum min max mean
}

type AggregationConfig struct {
	Every       int    `json:"every"`
	Function    string `json:"function"` // 可能的值："mean", "sum", "min", "max"
	CreateEmpty bool   `json:"create_empty"`
}

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
			|> aggregateWindow(every: %dm, fn: %s, createEmpty: %t)
			|> yield(name: "mean")
	`, iqc.Bucket, timeRange, filterClause, iqc.Measurement, iqc.Aggregation.Every, iqc.Aggregation.Function, iqc.Aggregation.CreateEmpty)
}
func (iqc *InfluxQueryConfig) GenerateFluxQueryString() string {
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
			|> yield(name: "first")
	`, iqc.Bucket, timeRange, filterClause, iqc.Measurement, iqc.Aggregation.Every, iqc.Aggregation.Function, iqc.Aggregation.CreateEmpty)
}
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

type CalcCache struct {
	ID     uint             `json:"id"`
	Param  []CalcParamCache `json:"param"`
	Cron   string           `json:"cron"`
	Script string           `json:"script"`
	Offset int64            `json:"offset"`
}

type CalcParamCache struct {
	Protocol           string `json:"protocol"`
	IdentificationCode string `json:"identification_code"` // 设备标识码
	DeviceUid          int    `json:"device_uid"`          // MQTT客户端表的外键ID

	Name       string `json:"name"`                                                  // 参数名称
	SignalName string `gorm:"signal_name"  json:"signal_name" structs:"signal_name"` // 信号表 name
	Reduce     string `json:"reduce"`                                                // 数据聚合方式 1. 求和 2. 平均值 3. 最大值 4. 最小值 4. 原始
	CalcRuleId int    `json:"calc_rule_id"`                                          // CalcRule 主键
	SignalId   int    `json:"signal_id" structs:"signal_id"`                         // 信号表的外键ID
}
type Event struct {
	StartTime int64 `json:"start_time" bson:"start_time"`
	EndTime   int64 `json:"end_time" bson:"end_time"`
	ID        int   `json:"id" bson:"id"`
}

type WaringRowQuery struct {
	ID          int   `json:"ID"`
	UpTimeStart int64 `json:"up_time_start"`
	UpTimeEnd   int64 `json:"up_time_end"`
}

type DeviceGroupCreateParam struct {
	GroupId  int   `json:"group_id"`
	DeviceId []int `json:"device_id"`
}

type ProductionPlanCreateParam struct {
	ID                      uint                     `json:"id,omitempty" structs:"id"`         // 生产计划名称
	Name                    string                   `json:"name" structs:"name"`               // 生产计划名称
	StartDate               time.Time                `json:"start_date" structs:"start_date"`   // 生产计划开始日期
	EndDate                 time.Time                `json:"end_date" structs:"end_date"`       // 生产计划结束日期
	Status                  string                   `json:"status" structs:"status"`           // 计划状态（准备中,进行中, 已完成）
	Description             string                   `json:"description" structs:"description"` // 生产计划描述
	ProductPlanCreateParams []ProductPlanCreateParam `json:"product_plans" `
}

func (p *ProductionPlanCreateParam) ValidateUniqueIDs() error {
	idMap := make(map[uint]bool)
	for _, param := range p.ProductPlanCreateParams {
		if _, exists := idMap[param.ProductID]; exists {
			return errors.New("error: duplicate ID found")
		}
		idMap[param.ProductID] = true
	}
	return nil
}

type ProductionPlanChangeParam struct {
	ID     uint   `json:"id,omitempty" structs:"id"` // 生产计划名称
	Status string `json:"status" structs:"status"`   // 计划状态（例如：准备中,进行中, 已完成）
}

type ProductPlanCreateParam struct {
	ProductID uint `json:"product_id" structs:"product_id"` // 关联的产品ID
	Quantity  int  `json:"quantity" structs:"quantity"`     // 计划生产数量
}

type ShipmentRecordCreateParam struct {
	ID              uint      `json:"id,omitempty" structs:"id"`                   // 发货记录ID
	ShipmentDate    time.Time `json:"shipment_date" structs:"shipment_date"`       // 发货日期
	Technician      string    `json:"technician" structs:"technician"`             // 发货人员
	CustomerName    string    `json:"customer_name" structs:"customer_name"`       // 客户名称
	CustomerAddress string    `json:"customer_address" structs:"customer_address"` // 客户地址
	TrackingNumber  string    `json:"tracking_number" structs:"tracking_number"`   // 跟踪号码
	Status          string    `json:"status" structs:"status"`                     // 发货状态（例如：pending, shipped, delivered）
	Description     string    `json:"description" structs:"description"`           // 发货描述
	CustomerPhone   string    `json:"customer_phone" structs:"customer_phone"`     // 客户手机

	ProductPlanCreateParams []ProductPlanCreateParam `json:"product_plans"`
}

func (s *ShipmentRecordCreateParam) ValidateUniqueIDs() error {
	idMap := make(map[uint]bool)
	for _, param := range s.ProductPlanCreateParams {
		if _, exists := idMap[param.ProductID]; exists {
			return errors.New("error: duplicate ID found")
		}
		idMap[param.ProductID] = true
	}
	return nil
}

type UserBindRoleParam struct {
	UserId  int   `json:"user_id"`
	RoleIds []int `json:"role_id"`
}

type UserBindDeptParam struct {
	UserId  int   `json:"user_id"`
	DeptIds []int `json:"dept_id"`
}

type UserBindDeviceInfoParam struct {
	UserId        int   `json:"user_id"`
	DeviceInfoIds []int `json:"device_info_id"`
}

type DeviceBindMqttClientParam struct {
	DeviceId     int   `json:"device_id"`
	MqttClientId []int `json:"mqtt_client_id"`
}
type DeviceBindTcpParam struct {
	DeviceId     int   `json:"device_id"`
	TcpHandlerId []int `json:"tcp_handler_id"`
}
type DeviceBindHTTPParam struct {
	DeviceId      int   `json:"device_id"`
	HttpHandlerId []int `json:"http_handler_id"`
}

type DeviceGroupBindMqttClientParam struct {
	DeviceGroupId uint `json:"device_group_id" structs:"device_group_id"` // 设备组ID

	MqttClientId []int `json:"mqtt_client_id"`
}

type LoginParam struct {
	UserName string `json:"user_name" form:"user_name"` // 用户名
	Password string `json:"password" form:"password"`   // 密码
}

type SimUseHistoryResp struct {
	ID           uint   `json:"id" structs:"id"`                         // 历史记录ID
	SimId        uint   `json:"sim_id" structs:"sim_id"`                 // 物联网卡ID
	DeviceInfoId uint   `json:"device_info_id" structs:"device_info_id"` // 设备ID
	Description  string `json:"description" structs:"description"`       // 描述
	SN           string `json:"sn" structs:"sn"`                         // 序列号
}

type TransmitScriptParam struct {
	DataRowList []common.DataRowList `json:"data_row_list"`
	Script      string               `json:"script"`
}

type DeviceInfoRes struct {
	ProductId         uint         `json:"product_id" structs:"product_id"`                                                               // 产品ID
	SN                string       `json:"sn" structs:"sn"`                                                                               // 设备编号
	ManufacturingDate ut.LocalTime `json:"manufacturing_date,omitempty" gorm:"type:DATETIME; default:NULL;" structs:"manufacturing_date"` // 制造日期
	ProcurementDate   ut.LocalTime `json:"procurement_date,omitempty" gorm:"type:DATETIME; default:NULL;" structs:"procurement_date"`     // 采购日期
	Source            int          `json:"source" structs:"source"`                                                                       // 设备来源,1: 内部,2: 外源
	WarrantyExpiry    ut.LocalTime `json:"warranty_expiry,omitempty" gorm:"type:DATETIME; default:NULL;" structs:"warranty_expiry"`       // 保修截止日期
	PushInterval      int          `json:"push_interval,omitempty" structs:"push_interval"`                                               // 推送间隔（秒）
	ErrorRate         float64      `json:"error_rate,omitempty" structs:"error_rate"`                                                     // 推送时间误差（秒）
	gorm.Model        `structs:"-"`
	ProductName       string `gorm:"-" json:"product_name,omitempty" ` // 产品名称
	Protocol          string `json:"protocol,omitempty" structs:"protocol,omitempty"`                                                                         // 协议


}
