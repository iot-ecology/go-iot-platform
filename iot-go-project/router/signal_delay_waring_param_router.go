package router

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type SignalDelayWaringParamApi struct {
}

var SignalDelayWaringParamBiz = biz.SignalDelayWaringParamBiz{}

// CreateSignalDelayWaring
// @Summary 创建脚本报警参数
// @Description 创建脚本报警参数
// @Tags signal-delay-waring-param
// @Accept json
// @Produce json
// @Param SignalDelayWaringParam body models.SignalDelayWaringParam true "脚本报警参数"
// @Success 201 {object} servlet.JSONResult{data=models.SignalDelayWaringParam} "创建成功的脚本报警参数"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /signal-delay-waring-param/create [post]
func (api *SignalDelayWaringParamApi) CreateSignalDelayWaring(c *gin.Context) {
	var SignalDelayWaringParam models.SignalDelayWaringParam
	if err := c.ShouldBindJSON(&SignalDelayWaringParam); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 SignalDelayWaringParam 是否被正确初始化
	if SignalDelayWaringParam.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&SignalDelayWaringParam)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	configBytes, _ := json.Marshal(SignalDelayWaringParam)

	glob.GRedis.LPush(context.Background(), "delay_param", configBytes)
	// 返回创建成功的脚本报警参数
	servlet.Resp(c, SignalDelayWaringParam)
}

// UpdateSignalDelayWaring
// @Summary 更新一个脚本报警参数
// @Description 更新一个脚本报警参数
// @Tags signal-delay-waring-param
// @Accept json
// @Produce json
// @Param id path int true "脚本报警参数id"
// @Param SignalDelayWaringParam body models.SignalDelayWaringParam true "脚本报警参数"
// @Success 200 {object}  servlet.JSONResult{data=models.SignalDelayWaringParam} "脚本报警参数"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "脚本报警参数未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /signal-delay-waring-param/update [post]
func (api *SignalDelayWaringParamApi) UpdateSignalDelayWaring(c *gin.Context) {
	var req models.SignalDelayWaringParam
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.SignalDelayWaringParam
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "SignalDelayWaringParam not found")
		return
	}
	configBytes, _ := json.Marshal(old)
	count, err := glob.GRedis.LRem(context.Background(), "delay_param", 0, configBytes).Result()
	if err != nil {
		zap.S().Errorf("err %+v", err)
	}
	glob.GLog.Sugar().Infof("删除键 %d 个", count)

	var newV models.SignalDelayWaringParam
	newV = old
	newV.MqttClientId = req.MqttClientId
	newV.Name = req.Name
	newV.SignalName = req.SignalName
	newV.SignalId = req.SignalId

	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}

	configBytes, _ = json.Marshal(newV)

	glob.GRedis.LPush(context.Background(), "delay_param", configBytes)

	servlet.Resp(c, old)
}

// PageSignalDelayWaring
// @Summary 分页查询脚本报警参数
// @Description 分页查询脚本报警参数
// @Tags signal-delay-waring-param
// @Accept json
// @Produce json
// @Param name query string false "名称" maxlength(100)
// @Param signal_delay_waring_id query int false "脚本报警id"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.SignalDelayWaringParam}} "脚本报警参数"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /signal-delay-waring-param/page [get]
func (api *SignalDelayWaringParamApi) PageSignalDelayWaring(c *gin.Context) {
	var name = c.Query("name")
	var signalDelayWaringId = c.Query("signal_delay_waring_id")
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

	data, err := SignalDelayWaringParamBiz.PageData(name, signalDelayWaringId, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteSignalDelayWaring
// @Tags      signal-delay-waring-param
// @Summary   删除脚本报警参数
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /signal-delay-waring-param/delete/:id [post]
func (api *SignalDelayWaringParamApi) DeleteSignalDelayWaring(c *gin.Context) {
	var SignalDelayWaringParam models.SignalDelayWaringParam

	param := c.Param("id")

	result := glob.GDb.First(&SignalDelayWaringParam, param)
	if result.Error != nil {
		servlet.Error(c, "SignalDelayWaringParam not found")

		return
	}

	if result := glob.GDb.Delete(&SignalDelayWaringParam); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}
