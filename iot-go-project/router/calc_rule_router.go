package router

import (
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type CalcRuleApi struct{}

var calcRuleBiz = biz.CalcRuleBiz{}
var calcRunBiz = biz.CalcRunBiz{}

// CreateCalcRule
// @Summary 创建计算规则
// @Description 创建计算规则
// @Tags calc-rule
// @Accept json
// @Produce json
// @Param CalcRule body models.CalcRule true "计算规则"
// @Success 201 {object} servlet.JSONResult{data=models.CalcRule} "创建成功的计算规则"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /calc-rule/create [post]
func (api *CalcRuleApi) CreateCalcRule(c *gin.Context) {
	var CalcRule models.CalcRule
	if err := c.ShouldBindJSON(&CalcRule); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 CalcRule 是否被正确初始化
	if CalcRule.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&CalcRule)
	calcRunBiz.RefreshRule(CalcRule.ID)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的计算规则
	servlet.Resp(c, CalcRule)
	return
}

// UpdateCalcRule
// @Summary 更新一个计算规则
// @Description 更新一个计算规则
// @Tags calc-rule
// @Accept json
// @Produce json
// @Param id path int true "计算规则id"
// @Param CalcRule body models.CalcRule true "计算规则"
// @Success 200 {object}  servlet.JSONResult{data=models.CalcRule} "计算规则"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "计算规则未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /calc-rule/update [post]
func (api *CalcRuleApi) UpdateCalcRule(c *gin.Context) {
	var req models.CalcRule
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.CalcRule
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "CalcRule not found")
		return
	}
	if old.Start {
		servlet.Error(c, "已启动不能修改.")
		return
	}

	var newV models.CalcRule
	newV = old
	newV.Name = req.Name
	newV.Script = req.Script
	newV.Cron = req.Cron
	newV.Offset = req.Offset
	newV.MockValue = ""
	m := structs.Map(newV)
	result = glob.GDb.Table("calc_rules").Where("id = ?", newV.ID).Updates(m)
	calcRunBiz.RefreshRule(old.ID)
	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageCalcRule
// @Summary 分页查询计算规则
// @Description 分页查询计算规则
// @Tags calc-rule
// @Accept json
// @Produce json
// @Param name query string false "名称"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.CalcRule}} "计算规则"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /calc-rule/page [get]
func (api *CalcRuleApi) PageCalcRule(c *gin.Context) {
	var name = c.Query("name")
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

	data, err := calcRuleBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
	return
}

// DeleteCalcRule
// @Tags      calc-rule
// @Summary   删除计算规则
// @Produce   application/json
// @Router    /calc-rule/delete/:id [post]
func (api *CalcRuleApi) DeleteCalcRule(c *gin.Context) {
	var CalcRule models.CalcRule

	param := c.Param("id")

	result := glob.GDb.First(&CalcRule, param)
	if result.Error != nil {
		servlet.Error(c, "CalcRule not found")

		return
	}

	if result := glob.GDb.Delete(&CalcRule); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// StartCalcRule
// @Tags      calc-rule
// @Summary   启动任务
// @Produce   application/json
// @Router    /calc-rule/start/:id [post]
func (api *CalcRuleApi) StartCalcRule(c *gin.Context) {

	param := c.Param("id")

	start := calcRunBiz.Start(param)
	if start {

		servlet.Resp(c, "启动成功")
		return
	} else {
		servlet.Resp(c, "启动失败")
		return
	}

}

// MockCalcRule
// @Tags      calc-rule
// @Summary   模拟执行
// @Param CalcRule body servlet.Event true "模拟参数"
// @Produce   application/json
// @Router    /calc-rule/mock [post]
func (api *CalcRuleApi) MockCalcRule(c *gin.Context) {

	var req servlet.Event
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	start := calcRunBiz.MockCalc(req.StartTime, req.EndTime, req.ID)

	servlet.Resp(c, start)
	return

}

// Refresh
// @Tags      calc-rule
// @Summary   刷新规则
// @Produce   application/json
// @Router    /calc-rule/refresh/:id [post]
func (api *CalcRuleApi) Refresh(c *gin.Context) {
	param := c.Param("id")

	calcRunBiz.RefreshRule(param)
	servlet.Resp(c, "刷新成功")
	return

}

// StopCalcRule
// @Tags      calc-rule
// @Summary   停止任务
// @Produce   application/json
// @Router    /calc-rule/stop/:id [post]
func (api *CalcRuleApi) StopCalcRule(c *gin.Context) {

	param := c.Param("id")

	start := calcRunBiz.Stop(param)
	if start {

		servlet.Resp(c, "停止成功")
		return
	} else {
		servlet.Resp(c, "停止失败")
		return
	}

}

// CalcRuleResult
// @Tags      calc-rule
// @Summary   获取计算结果
// @Param rule_id path int true "规则id"
// @Param start_time path int true "开始时间"
// @Param end_time path int true "结束时间"
// @Produce   application/json
// @Router    /calc-rule/rd [get]
func (api *CalcRuleApi) CalcRuleResult(c *gin.Context) {
	ruleIdstr := c.Query("rule_id")
	startTimestr := c.Query("start_time")
	endTimestr := c.Query("end_time")

	ruleId, err := strconv.ParseInt(ruleIdstr, 10, 64)
	if err != nil {
		servlet.Error(c, "invalid rule_id format")
		return
	}

	startTime, err1 := strconv.ParseInt(startTimestr, 10, 64)
	if err1 != nil {
		servlet.Error(c, "invalid start_time format")
		return
	}
	endTime, err2 := strconv.ParseInt(endTimestr, 10, 64)
	if err2 != nil {
		servlet.Error(c, "invalid end_time format")
		return
	}

	start := calcRunBiz.QueryRuleExData(ruleId, startTime, endTime)

	servlet.Resp(c, start)
	return
}
