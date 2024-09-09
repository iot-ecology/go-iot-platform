package transmit_mqtt

import (
	"github.com/gin-gonic/gin"
	transmit "igp/biz/transmit/mqtt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type RabbitmqTransmitBindApi struct{}

var RabbitmqTransmitBindBiz = transmit.RabbitmqTransmitBindBiz{}

// CreateRabbitmqTransmitBind
// @Summary 创建RabbitmqTransmitBind
// @Description 创建RabbitmqTransmitBind
// @Tags RabbitmqTransmitBinds
// @Accept json
// @Produce json
// @Param RabbitmqTransmitBind body models.RabbitmqTransmitBind true "RabbitmqTransmitBind"
// @Success 200 {object} servlet.JSONResult{data=models.RabbitmqTransmitBind} "创建成功的RabbitmqTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /RabbitmqTransmitBind/create [post]
func (api *RabbitmqTransmitBindApi) CreateRabbitmqTransmitBind(c *gin.Context) {
	var RabbitmqTransmitBind models.RabbitmqTransmitBind
	if err := c.ShouldBindJSON(&RabbitmqTransmitBind); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 RabbitmqTransmitBind 是否被正确初始化
	if RabbitmqTransmitBind.Exchange == "" {
		servlet.Error(c, "Exchange不能为空")
		return
	}

	result := glob.GDb.Create(&RabbitmqTransmitBind)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	RabbitmqTransmitBindBiz.HandlerRedis(RabbitmqTransmitBind)
	// 返回创建成功的RabbitmqTransmitBind
	servlet.Resp(c, RabbitmqTransmitBind)
}

// UpdateRabbitmqTransmitBind
// @Summary 更新一个RabbitmqTransmitBind
// @Description 更新一个RabbitmqTransmitBind
// @Tags RabbitmqTransmitBinds
// @Accept json
// @Produce json
// @Param RabbitmqTransmitBind body models.RabbitmqTransmitBind true "RabbitmqTransmitBind"
// @Success 200 {object}  servlet.JSONResult{data=models.RabbitmqTransmitBind} "RabbitmqTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "RabbitmqTransmitBind未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /RabbitmqTransmitBind/update [post]
func (api *RabbitmqTransmitBindApi) UpdateRabbitmqTransmitBind(c *gin.Context) {
	var req models.RabbitmqTransmitBind
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.RabbitmqTransmitBind
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "RabbitmqTransmitBind not found")
		return
	}

	var newV models.RabbitmqTransmitBind
	newV = old
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	RabbitmqTransmitBindBiz.HandlerRedis(newV)
	servlet.Resp(c, old)
}

// PageRabbitmqTransmitBind
// @Summary 分页查询RabbitmqTransmitBind
// @Description 分页查询RabbitmqTransmitBind
// @Tags RabbitmqTransmitBinds
// @Accept json
// @Produce json
// @Param exchange query string false "表名" default(0)
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.RabbitmqTransmitBind}} "RabbitmqTransmitBind"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /RabbitmqTransmitBind/page [get]
func (api *RabbitmqTransmitBindApi) PageRabbitmqTransmitBind(c *gin.Context) {
	var name = c.Query("exchange")
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

	data, err := RabbitmqTransmitBindBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteRabbitmqTransmitBind
// @Tags      RabbitmqTransmitBinds
// @Summary   删除RabbitmqTransmitBind
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /RabbitmqTransmitBind/delete/:id [post]
func (api *RabbitmqTransmitBindApi) DeleteRabbitmqTransmitBind(c *gin.Context) {
	var RabbitmqTransmitBind models.RabbitmqTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&RabbitmqTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "RabbitmqTransmitBind not found")

		return
	}

	if result := glob.GDb.Delete(&RabbitmqTransmitBind); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	RabbitmqTransmitBind.Enable = false
	RabbitmqTransmitBindBiz.HandlerRedis(RabbitmqTransmitBind)

	servlet.Resp(c, "删除成功")
}

// ByIdRabbitmqTransmitBind
// @Tags      RabbitmqTransmitBinds
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /RabbitmqTransmitBind/:id [get]
func (api *RabbitmqTransmitBindApi) ByIdRabbitmqTransmitBind(c *gin.Context) {
	var RabbitmqTransmitBind models.RabbitmqTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&RabbitmqTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "RabbitmqTransmitBind not found")

		return
	}

	servlet.Resp(c, RabbitmqTransmitBind)
}
