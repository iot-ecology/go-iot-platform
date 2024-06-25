package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type CalcParamApi struct{}

var calcParamBiz = biz.CalcParamBiz{}

// CreateCalcParam
// @Summary 创建计算参数
// @Description 创建计算参数
// @Tags calc-param
// @Accept json
// @Produce json
// @Param CalcParam body models.CalcParam true "计算参数"
// @Success 201 {object} servlet.JSONResult{data=models.CalcParam} "创建成功的计算参数"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /calc-param/create [post]
func (api *CalcParamApi) CreateCalcParam(c *gin.Context) {
	var CalcParam models.CalcParam
	if err := c.ShouldBindJSON(&CalcParam); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 CalcParam 是否被正确初始化
	if CalcParam.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&CalcParam)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	calcRunBiz.RefreshRule(CalcParam.CalcRuleId)

	// 返回创建成功的计算参数
	servlet.Resp(c, CalcParam)
}

// UpdateCalcParam
// @Summary 更新一个计算参数
// @Description 更新一个计算参数
// @Tags calc-param
// @Accept json
// @Produce json
// @Param id path int true "计算参数id"
// @Param CalcParam body models.CalcParam true "计算参数"
// @Success 200 {object}  servlet.JSONResult{data=models.CalcParam} "计算参数"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "计算参数未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /calc-param/update [post]
func (api *CalcParamApi) UpdateCalcParam(c *gin.Context) {
	var req models.CalcParam
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.CalcParam
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "CalcParam not found")
		return
	}

	var newV models.CalcParam
	newV = old
	newV.Name = req.Name
	newV.Reduce = req.Reduce
	newV.SignalName = req.SignalName
	newV.MqttClientId = req.MqttClientId
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	calcRunBiz.RefreshRule(old.CalcRuleId)

	servlet.Resp(c, old)
}

// PageCalcParam
// @Summary 分页查询计算参数
// @Description 分页查询计算参数
// @Tags calc-param
// @Accept json
// @Produce json
// @Param name query string false "客户端名称" maxlength(100)
// @Param  mqtt_client_id query int false "mqtt client 表id"
// @Param  rule_id query int false "计算规则 表id"
// @Param  signal_name query string false "信号名称"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.CalcParam}} "计算参数"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /calc-param/page [get]
func (api *CalcParamApi) PageCalcParam(c *gin.Context) {
	var name = c.Query("name")
	var mqttClientId = c.Query("mqtt_client_id")
	var ruleId = c.Query("rule_id")
	var signalName = c.Query("signal_name")
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

	data, err := calcParamBiz.PageData(name, mqttClientId, signalName, ruleId, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteCalcParam
// @Tags      calc-param
// @Summary   删除计算参数
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /calc-param/delete/:id [post]
func (api *CalcParamApi) DeleteCalcParam(c *gin.Context) {
	var CalcParam models.CalcParam

	param := c.Param("id")

	result := glob.GDb.First(&CalcParam, param)
	if result.Error != nil {
		servlet.Error(c, "CalcParam not found")

		return
	}

	if result := glob.GDb.Delete(&CalcParam); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}
