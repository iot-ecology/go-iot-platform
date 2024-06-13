package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type SignalApi struct{}

var bizSignal = biz.SignalBiz{}

// CreateSignal
// SignalApi 结构体用于封装与信号相关的API
// @Summary 创建一个新的信号
// @Description 创建一个新的信号记录
// @Tags signals
// @Accept json
// @Produce json
// @Param signal body models.Signal true "信号信息"
// @Success 201 {object} servlet.JSONResult{data=models.Signal} "创建成功的信号"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /signal/create [post]
func (api *SignalApi) CreateSignal(c *gin.Context) {
	var signal models.Signal
	if err := c.ShouldBindJSON(&signal); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 signal 是否被正确初始化
	if signal.Name == "" {
		servlet.Error(c, "信号名称不能为空")
		return
	}

	result := glob.GDb.Create(&signal)
	bizSignal.SetSignalCache(&signal)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的信号
	servlet.Resp(c, signal)
}

// UpdateSignal
// @Summary 更新一个已有的信号
// @Description 根据ID更新一个已有的信号记录
// @Tags signals
// @Accept json
// @Produce json
// @Param id path int true "信号的ID"
// @Param signal body models.Signal true "信号信息"
// @Success 200 {object}  servlet.JSONResult{data=models.Signal} "更新成功的信号"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "信号未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /signal/update [post]
func (api *SignalApi) UpdateSignal(c *gin.Context) {
	var req models.Signal
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.Signal
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "Signal not found")
		return
	}

	bizSignal.RemoveSignalCache(&old)
	var newV models.Signal
	newV = old
	newV.Name = req.Name
	newV.Type = req.Type
	newV.Alias = req.Alias
	newV.Unit = req.Unit
	newV.CacheSize = req.CacheSize

	result = glob.GDb.Model(&newV).Updates(newV)

	removeOldCache(newV.CacheSize, newV.ID, newV.MqttClientId)
	bizSignal.SetSignalCache(&newV)
	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// removeOldCache 从Redis有序集合中移除旧缓存
//
// 参数：
// threshold：int类型，表示保留的元素数量阈值
// signalId：uint类型，表示信号ID
// mqttId：int类型，表示MQTT连接ID
//
// 返回值：无
func removeOldCache(threshold int, signalId uint, mqttId int) {
	ctx := context.Background()

	redisKey := "signal_delay_warning:" + strconv.Itoa(mqttId) + ":" + strconv.Itoa(int(signalId))
	count, err := glob.GRedis.ZCard(ctx, redisKey).Result()
	if err != nil {
		panic(err)
	}
	if count <= int64(threshold) {
		fmt.Println("No need to remove elements.")
		return
	}
	extraCount := count - int64(threshold)
	deletedCount, err := glob.GRedis.ZRemRangeByRank(ctx, redisKey, 0, int64(extraCount)-1).Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Removed %d extra elements from sorted set '%s'\n", deletedCount, redisKey)
}

// PageSignal
// @Summary 分页查询信号
// @Description 根据条件分页查询信号列表
// @Tags signals
// @Accept json
// @Produce json
// @Param mqtt_client_id query string false "MQTT客户端ID"
// @Param type query string false "数据类型,数字、中文"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.Signal}} "信号列表"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /signal/page [get]
func (api *SignalApi) PageSignal(c *gin.Context) {
	var mqqtClientId = c.Query("mqtt_client_id")
	var ty = c.Query("type")
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

	data, err := bizSignal.PageSignal(mqqtClientId, ty, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteSignal
// @Tags      signals
// @Param id path int true "主键"
// @Summary   删除信号
// @Produce   application/json
// @Router    /signal/delete/:id [post]
func (api *SignalApi) DeleteSignal(c *gin.Context) {
	var signal models.Signal

	param := c.Param("id")

	result := glob.GDb.First(&signal, param)
	bizSignal.RemoveSignalCache(&signal)
	if result.Error != nil {
		servlet.Error(c, "Signal not found")

		return
	}

	if result := glob.GDb.Delete(&signal); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}
