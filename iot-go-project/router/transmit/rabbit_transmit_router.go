package transmit

import (
	"github.com/gin-gonic/gin"
	"igp/biz/transmit"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type RabbitmqTransmitApi struct{}

var RabbitmqTransmitBiz = transmit.RabbitTransmitBiz{}

// CreateRabbitmqTransmit
// @Summary 创建Rabbit消息队列管理
// @Description 创建Rabbit消息队列管理
// @Tags RabbitmqTransmits
// @Accept json
// @Produce json
// @Param RabbitmqTransmit body models.RabbitmqTransmit true "Rabbit消息队列管理"
// @Success 200 {object} servlet.JSONResult{data=models.RabbitmqTransmit} "创建成功的Rabbit消息队列管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /RabbitmqTransmit/create [post]
func (api *RabbitmqTransmitApi) CreateRabbitmqTransmit(c *gin.Context) {
	var RabbitmqTransmit models.RabbitmqTransmit
	if err := c.ShouldBindJSON(&RabbitmqTransmit); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 RabbitmqTransmit 是否被正确初始化
	if RabbitmqTransmit.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&RabbitmqTransmit)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	RabbitmqTransmitBiz.SetRedis(RabbitmqTransmit)
	// 返回创建成功的Rabbit消息队列管理
	servlet.Resp(c, RabbitmqTransmit)
}

// UpdateRabbitmqTransmit
// @Summary 更新一个Rabbit消息队列管理
// @Description 更新一个Rabbit消息队列管理
// @Tags RabbitmqTransmits
// @Accept json
// @Produce json
// @Param RabbitmqTransmit body models.RabbitmqTransmit true "Rabbit消息队列管理"
// @Success 200 {object}  servlet.JSONResult{data=models.RabbitmqTransmit} "Rabbit消息队列管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Rabbit消息队列管理未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /RabbitmqTransmit/update [post]
func (api *RabbitmqTransmitApi) UpdateRabbitmqTransmit(c *gin.Context) {
	var req models.RabbitmqTransmit
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.RabbitmqTransmit
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "RabbitmqTransmit not found")
		return
	}

	var newV models.RabbitmqTransmit
	newV = old
	newV.Name = req.Name
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	RabbitmqTransmitBiz.SetRedis(newV)
	servlet.Resp(c, newV)
}

// PageRabbitmqTransmit
// @Summary 分页查询Rabbit消息队列管理
// @Description 分页查询Rabbit消息队列管理
// @Tags RabbitmqTransmits
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.RabbitmqTransmit}} "Rabbit消息队列管理"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /RabbitmqTransmit/page [get]
func (api *RabbitmqTransmitApi) PageRabbitmqTransmit(c *gin.Context) {
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

	data, err := RabbitmqTransmitBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteRabbitmqTransmit
// @Tags      RabbitmqTransmits
// @Summary   删除Rabbit消息队列管理
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /RabbitmqTransmit/delete/:id [post]
// @Success 200 {object}  servlet.JSONResult{data=string} 
func (api *RabbitmqTransmitApi) DeleteRabbitmqTransmit(c *gin.Context) {
	var RabbitmqTransmit models.RabbitmqTransmit

	param := c.Param("id")

	result := glob.GDb.First(&RabbitmqTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "RabbitmqTransmit not found")

		return
	}

	if result := glob.GDb.Delete(&RabbitmqTransmit); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	RabbitmqTransmitBiz.DeleteRedis(RabbitmqTransmit)
	servlet.Resp(c, "删除成功")
}

// ByIdRabbitmqTransmit
// @Tags      RabbitmqTransmits
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /RabbitmqTransmit/:id [get]
// @Success 200 {object}  servlet.JSONResult{data=models.RabbitmqTransmit} 
func (api *RabbitmqTransmitApi) ByIdRabbitmqTransmit(c *gin.Context) {
	var RabbitmqTransmit models.RabbitmqTransmit

	param := c.Param("id")

	result := glob.GDb.First(&RabbitmqTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "RabbitmqTransmit not found")

		return
	}

	servlet.Resp(c, RabbitmqTransmit)
}
