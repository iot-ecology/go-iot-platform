package router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type MqttApi struct{}

var bizMqtt = biz.MqttClientBiz{}
var nodeBiz = biz.NodeBiz{}
var scriptBiz = biz.ScirptBiz{}

// CreateMqtt
// @Tags      MQTT
// @Summary   创建MQTT客户端
// @accept    application/json
// @Produce   application/json
// @Param     data  body      models.MqttClient true "创建参数"
// @Success   200  {object}  servlet.JSONResult{data=models.MqttClient}
// @Router    /mqtt/create [post]
func (s *MqttApi) CreateMqtt(c *gin.Context) {
	mqttClient := models.MqttClient{}
	err := c.ShouldBind(&mqttClient)
	if err != nil {
		glob.GLog.Sugar().Error("操作异常", err)
		panic(err)
		return
	}

	servlet.Resp(c, bizMqtt.CreateMqtt(mqttClient))
	return
}

// UpdateMqtt
// @Tags      MQTT
// @Summary   修改MQTT客户端
// @accept    application/json
// @Produce   application/json
// @Param     data  body      models.MqttClient true "创建参数"
// @Success   200  {object}  servlet.JSONResult{data=models.MqttClient}
// @Router    /mqtt/update [post]
func (s *MqttApi) UpdateMqtt(c *gin.Context) {
	var req models.MqttClient

	if err := c.ShouldBindJSON(&req); err != nil {
		servlet.Error(c, err.Error())

		return
	}
	var old models.MqttClient

	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {
		servlet.Error(c, "MQTT_CLIENT not found")
		return
	}
	var newV models.MqttClient
	newV = old
	newV.Subtopic = req.Subtopic
	newV.Username = req.Username
	newV.Password = req.Password
	newV.Host = req.Host
	newV.Port = req.Port
	// 更新记录
	result = glob.GDb.Model(&newV).Updates(newV)
	if result.Error != nil {
		servlet.Error(c, result.Error.Error())

		return
	}
	servlet.Resp(c, old)
}

// StartMqtt
// @Tags      MQTT
// @Summary   启动MQTT客户端
// @Produce   application/json
// @Param id query string false "mqtt_client表id"
// @Router    /mqtt/start [get]
func (s *MqttApi) StartMqtt(c *gin.Context) {
	var id = c.Query("id")

	start := bizMqtt.Start(id)
	param := nodeBiz.SendCreateParam(start)
	glob.GLog.Sugar().Info(param)

	var m map[string]interface{}
	err := json.Unmarshal([]byte(param), &m)
	if err != nil {
		zap.S().Error("Error unmarshalling JSON", zap.Error(err))
		// 这里可以返回错误或者处理错误
		return
	}
	msg := m["message"]

	servlet.Resp2(c, fmt.Sprintf("%v", msg))

	return
}

// StopMqtt
// @Tags      MQTT
// @Summary   停止MQTT客户端
// @Produce   application/json
// @Param id query string false "mqtt_client表id"
// @Router    /mqtt/stop [get]
func (s *MqttApi) StopMqtt(c *gin.Context) {
	var id = c.Query("id")

	start := bizMqtt.Start(id)
	param := nodeBiz.SendStopParam(start)
	glob.GLog.Sugar().Info(param)

	var m map[string]interface{}
	err := json.Unmarshal([]byte(param), &m)
	if err != nil {
		zap.S().Error("Error unmarshalling JSON", zap.Error(err))
		// 这里可以返回错误或者处理错误
		return
	}
	glob.GLog.Sugar().Info(param)

	msg := m["message"]

	servlet.Resp2(c, fmt.Sprintf("%v", msg))

	return
}

// SendMqttMessage
// @Tags      MQTT
// @Summary   发送MQTT消息
// @Produce   application/json
// @Param id query string false "客户端ID"
// @Param     data  body      servlet.ParamStruct true "消息"
// @Router    /mqtt/send [post]
func (s *MqttApi) SendMqttMessage(c *gin.Context) {
	var id = c.Query("id")
	requestBody, err := c.GetRawData()
	if err != nil {
		servlet.Error(c, err.Error())
		return
	}

	param := nodeBiz.SendPushData(id, requestBody)

	var m map[string]interface{}

	err = json.Unmarshal([]byte(param), &m)
	if err != nil {
		zap.S().Error("Error unmarshalling JSON", zap.Error(err))
		// 这里可以返回错误或者处理错误
		return
	}
	glob.GLog.Sugar().Info(param)

	msg := m["message"]

	servlet.Resp2(c, fmt.Sprintf("%v", msg))

	return
}

// PageMqtt
// @Tags      MQTT
// @Summary   分页查询MQTT客户端
// @Produce   application/json
// @Param client_id query string false "客户端名称" maxlength(100)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success   200  {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.MqttClient}}
// @Router    /mqtt/page [get]
func (s *MqttApi) PageMqtt(c *gin.Context) {

	var name = c.Query("client_id")
	var page = c.DefaultQuery("page", "0")
	var pageSize = c.DefaultQuery("page_size", "10")
	parseUint, err := strconv.Atoi(page)
	if err != nil {
		servlet.Error(c, "无效的页码")
		return
	}
	u, err := strconv.Atoi(pageSize)

	if err != nil {
		servlet.Error(c, "无效的页长")

		return
	}

	data, err := bizMqtt.PageMqttData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")

		return
	}
	servlet.Resp(c, data)
	return
}

// NodeUsingStatus
// @Tags      MQTT
// @Summary   查询节点使用情况
// @Produce   application/json
// @Router    /mqtt/node-using-status [get]
func (s *MqttApi) NodeUsingStatus(c *gin.Context) {

	status := nodeBiz.SendNodeUsingStatus()

	glob.GLog.Sugar().Info(status)

	servlet.Resp(c, status)

}

// DeleteMqtt
// @Tags      MQTT
// @Summary   删除MQTT客户端
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /mqtt/delete/:id [post]
func (s *MqttApi) DeleteMqtt(c *gin.Context) {
	var mqttClient models.MqttClient

	param := c.Param("id")

	result := glob.GDb.First(&mqttClient, param)
	if result.Error != nil {
		servlet.Error(c, "MQTT_CLIENT not found")
		return
	}

	if result := glob.GDb.Delete(&mqttClient); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	glob.GRedis.HDel(context.Background(), "mqtt_script", param)

	servlet.Resp(c, "删除成功")
}

// SetScript
// @Tags      MQTT
// @Summary   设置解析脚本
// @accept    application/json
// @Param     data  body      servlet.MqttScript true "创建参数"
// @Produce   application/json
// @Router    /mqtt/set-script [post]
func (s *MqttApi) SetScript(c *gin.Context) {

	var scriptData servlet.MqttScript

	if err := c.BindJSON(&scriptData); err != nil {
		servlet.Error(c, "Invalid request body")
		return
	}

	var mqttClient models.MqttClient
	result := glob.GDb.First(&mqttClient, scriptData.ID)

	if result.Error != nil {
		servlet.Error(c, "MQTT_CLIENT not found")
		return
	}

	mqttClient.Script = scriptData.Script

	result = glob.GDb.Model(&mqttClient).Updates(mqttClient)

	if result.Error != nil {
		servlet.Error(c, "Failed to update script")

		return
	}
	bizMqtt.SetScriptRedis(mqttClient.ClientId, mqttClient.Script)
	servlet.Resp(c, "ok")
	return

}

// CheckScript
// @Tags      MQTT
// @Summary   设置解析脚本
// @accept    application/json
// @Param     data  body      servlet.CheckScriptReq true "创建参数"
// @Produce   application/json
// @Router    /mqtt/check-script [post]
func (s *MqttApi) CheckScript(c *gin.Context) {
	var req servlet.CheckScriptReq
	if err := c.BindJSON(&req); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := scriptBiz.CheckScript(req.Param, req.Script)
	if result == nil {
		servlet.Error(c, "执行脚本失败")
		return
	} else {
		servlet.Resp(c, result)
		return
	}

}
