package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"igp/glob"
	"igp/servlet"
	"reflect"
)

type InfluxDbApi struct{}

// QueryInfluxdb
// @Tags      DATA
// @Summary   数据查询
// @accept    application/json
// @Produce   application/json
// @Param     data  body      servlet.InfluxQueryConfig true "查询参数"
// @Success   200  {object}  servlet.JSONResult
// @Router    /query/influxdb [post]
func (s *InfluxDbApi) QueryInfluxdb(c *gin.Context) {
	json := servlet.InfluxQueryConfig{}
	err := c.ShouldBind(&json)
	if err != nil {
		glob.GLog.Sugar().Error("操作异常", err)
		panic(err)

		return
	}
	json.Bucket = glob.GConfig.InfluxConfig.Bucket
	query := json.GenerateFluxQuery()
	glob.GLog.Sugar().Info(query)
	result, err := glob.GInfluxdb.QueryAPI(glob.GConfig.InfluxConfig.Org).Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var v []map[string]interface{}

	for result.Next() {
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		values := result.Record().Values()
		fmt.Printf("value: %v\n", values)
		v = append(v, values)
	}
	field := groupByField(v)

	servlet.Resp(c, field)

	return
}

// QueryInfluxdbString
// @Tags      DATA
// @Summary   数据查询字符串
// @accept    application/json
// @Produce   application/json
// @Param     data  body      servlet.InfluxQueryConfig true "查询参数"
// @Success   200  {object}  servlet.JSONResult
// @Router    /query/str-influxdb [post]
func (s *InfluxDbApi) QueryInfluxdbString(c *gin.Context) {
	json := servlet.InfluxQueryConfig{}
	err := c.ShouldBind(&json)
	if err != nil {
		glob.GLog.Sugar().Error("操作异常", err)
		panic(err)

		return
	}
	json.Bucket = glob.GConfig.InfluxConfig.Bucket
	query := json.GenerateFluxQueryString()
	glob.GLog.Sugar().Info(query)
	result, err := glob.GInfluxdb.QueryAPI(glob.GConfig.InfluxConfig.Org).Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var v []map[string]interface{}

	for result.Next() {
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		values := result.Record().Values()
		fmt.Printf("value: %v\n", values)
		v = append(v, values)
	}
	field := groupByField(v)

	servlet.Resp(c, field)

	return
}

func groupByField(records []map[string]interface{}) map[string][]map[string]interface{} {
	grouped := make(map[string][]map[string]interface{})

	for _, record := range records {
		// 检查_field是否存在于记录中
		if fieldVal, ok := record["_field"]; ok {
			if fieldValStr, ok := fieldVal.(string); ok {
				// 使用_field的值作为分组的键
				grouped[fieldValStr] = append(grouped[fieldValStr], record)
			} else {
				fmt.Printf("Expected _field to be a string, got %s\n", reflect.TypeOf(fieldVal))
			}
		} else {
			fmt.Println("_field not found in record")
		}
	}

	return grouped
}
