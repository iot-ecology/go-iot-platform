package servlet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
	Time      int64  `json:"time"` // 秒级时间戳
	DeviceUid string `json:"device_uid"`

	DataRows []DataRow `json:"data"`
	Nc       string    `json:"nc"`
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
	return
}

func Resp2(c *gin.Context, msg string) {
	result := JSONResult{}
	result.Message = msg
	result.Code = 20000
	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, result)
	return
}
func Error(c *gin.Context, data interface{}) {
	result := JSONResult{}
	result.Message = "操作失败"
	result.Code = 40000
	result.Data = data
	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, result)
	return
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
	MqttClientId int    `json:"mqtt_client_id"`                                        // MQTT客户端表的外键ID
	Name         string `json:"name"`                                                  // 参数名称
	SignalName   string `gorm:"signal_name"  json:"signal_name" structs:"signal_name"` // 信号表 name
	Reduce       string `json:"reduce"`                                                // 数据聚合方式 1. 求和 2. 平均值 3. 最大值 4. 最小值 4. 原始
	CalcRuleId   int    `json:"calc_rule_id"`                                          // CalcRule 主键
	SignalId     int    `json:"signal_id" structs:"signal_id"`                         // 信号表的外键ID
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
