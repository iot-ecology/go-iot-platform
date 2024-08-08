package transmit_mqtt

import (
	"github.com/gin-gonic/gin"
	transmit "igp/biz/transmit/mqtt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type KafkaTransmitBindApi struct{}

var KafkaTransmitBindBiz = transmit.KafkaTransmitBindBiz{}

// CreateKafkaTransmitBind
// @Summary 创建KafkaTransmitBind
// @Description 创建KafkaTransmitBind
// @Tags KafkaTransmitBinds
// @Accept json
// @Produce json
// @Param KafkaTransmitBind body models.KafkaTransmitBind true "KafkaTransmitBind"
// @Success 201 {object} servlet.JSONResult{data=models.KafkaTransmitBind} "创建成功的KafkaTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /KafkaTransmitBind/create [post]
func (api *KafkaTransmitBindApi) CreateKafkaTransmitBind(c *gin.Context) {
	var KafkaTransmitBind models.KafkaTransmitBind
	if err := c.ShouldBindJSON(&KafkaTransmitBind); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 KafkaTransmitBind 是否被正确初始化
	if KafkaTransmitBind.Topic == "" {
		servlet.Error(c, "topic不能为空")
		return
	}

	result := glob.GDb.Create(&KafkaTransmitBind)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	KafkaTransmitBindBiz.HandlerRedis(KafkaTransmitBind)
	// 返回创建成功的KafkaTransmitBind
	servlet.Resp(c, KafkaTransmitBind)
}

// UpdateKafkaTransmitBind
// @Summary 更新一个KafkaTransmitBind
// @Description 更新一个KafkaTransmitBind
// @Tags KafkaTransmitBinds
// @Accept json
// @Produce json
// @Param KafkaTransmitBind body models.KafkaTransmitBind true "KafkaTransmitBind"
// @Success 200 {object}  servlet.JSONResult{data=models.KafkaTransmitBind} "KafkaTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "KafkaTransmitBind未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /KafkaTransmitBind/update [post]
func (api *KafkaTransmitBindApi) UpdateKafkaTransmitBind(c *gin.Context) {
	var req models.KafkaTransmitBind
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.KafkaTransmitBind
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "KafkaTransmitBind not found")
		return
	}

	var newV models.KafkaTransmitBind
	newV = old
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	KafkaTransmitBindBiz.HandlerRedis(newV)
	servlet.Resp(c, old)
}

// PageKafkaTransmitBind
// @Summary 分页查询KafkaTransmitBind
// @Description 分页查询KafkaTransmitBind
// @Tags KafkaTransmitBinds
// @Accept json
// @Produce json
// @Param topic query string false "主题" default(0)
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.KafkaTransmitBind}} "KafkaTransmitBind"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /KafkaTransmitBind/page [get]
func (api *KafkaTransmitBindApi) PageKafkaTransmitBind(c *gin.Context) {
	var name = c.Query("topic")
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

	data, err := KafkaTransmitBindBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteKafkaTransmitBind
// @Tags      KafkaTransmitBinds
// @Summary   删除KafkaTransmitBind
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /KafkaTransmitBind/delete/:id [post]
func (api *KafkaTransmitBindApi) DeleteKafkaTransmitBind(c *gin.Context) {
	var KafkaTransmitBind models.KafkaTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&KafkaTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "KafkaTransmitBind not found")

		return
	}

	if result := glob.GDb.Delete(&KafkaTransmitBind); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	KafkaTransmitBind.Enable = false
	KafkaTransmitBindBiz.HandlerRedis(KafkaTransmitBind)
	servlet.Resp(c, "删除成功")
}

// ByIdKafkaTransmitBind
// @Tags      KafkaTransmitBinds
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /KafkaTransmitBind/:id [get]
func (api *KafkaTransmitBindApi) ByIdKafkaTransmitBind(c *gin.Context) {
	var KafkaTransmitBind models.KafkaTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&KafkaTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "KafkaTransmitBind not found")

		return
	}

	servlet.Resp(c, KafkaTransmitBind)
}
