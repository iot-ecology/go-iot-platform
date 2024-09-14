package transmit

import (
	"github.com/gin-gonic/gin"
	"igp/biz/transmit"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type InfluxdbTransmitApi struct{}

var InfluxdbTransmitBiz = transmit.InfluxdbTransmitBiz{}

// CreateInfluxdbTransmit
// @Summary 创建Influxdb数据库管理
// @Description 创建Influxdb数据库管理
// @Tags InfluxdbTransmits
// @Accept json
// @Produce json
// @Param InfluxdbTransmit body models.InfluxdbTransmit true "Influxdb数据库管理"
// @Success 200 {object} servlet.JSONResult{data=models.InfluxdbTransmit} "创建成功的Influxdb数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /InfluxdbTransmit/create [post]
func (api *InfluxdbTransmitApi) CreateInfluxdbTransmit(c *gin.Context) {
	var InfluxdbTransmit models.InfluxdbTransmit
	if err := c.ShouldBindJSON(&InfluxdbTransmit); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := glob.GDb.Create(&InfluxdbTransmit)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	InfluxdbTransmitBiz.SetRedis(InfluxdbTransmit)
	// 返回创建成功的Influxdb数据库管理
	servlet.Resp(c, InfluxdbTransmit)
}

// UpdateInfluxdbTransmit
// @Summary 更新一个Influxdb数据库管理
// @Description 更新一个Influxdb数据库管理
// @Tags InfluxdbTransmits
// @Accept json
// @Produce json
// @Param InfluxdbTransmit body models.InfluxdbTransmit true "Influxdb数据库管理"
// @Success 200 {object}  servlet.JSONResult{data=models.InfluxdbTransmit} "Influxdb数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Influxdb数据库管理未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /InfluxdbTransmit/update [post]
func (api *InfluxdbTransmitApi) UpdateInfluxdbTransmit(c *gin.Context) {
	var req models.InfluxdbTransmit
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.InfluxdbTransmit
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "InfluxdbTransmit not found")
		return
	}

	var newV models.InfluxdbTransmit
	newV = old
	newV.Name =req.Name
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	InfluxdbTransmitBiz.SetRedis(newV)
	servlet.Resp(c, old)
}

// PageInfluxdbTransmit
// @Summary 分页查询Influxdb数据库管理
// @Description 分页查询Influxdb数据库管理
// @Tags InfluxdbTransmits
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.InfluxdbTransmit}} "Influxdb数据库管理"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /InfluxdbTransmit/page [get]
func (api *InfluxdbTransmitApi) PageInfluxdbTransmit(c *gin.Context) {
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

	data, err := InfluxdbTransmitBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteInfluxdbTransmit
// @Tags      InfluxdbTransmits
// @Summary   删除Influxdb数据库管理
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /InfluxdbTransmit/delete/:id [post]
// @Success 200 {object}  servlet.JSONResult{data=string} 
func (api *InfluxdbTransmitApi) DeleteInfluxdbTransmit(c *gin.Context) {
	var InfluxdbTransmit models.InfluxdbTransmit

	param := c.Param("id")

	result := glob.GDb.First(&InfluxdbTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "InfluxdbTransmit not found")

		return
	}

	if result := glob.GDb.Delete(&InfluxdbTransmit); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	InfluxdbTransmitBiz.DeleteRedis(InfluxdbTransmit)
	servlet.Resp(c, "删除成功")
}

// ByIdInfluxdbTransmit
// @Tags      InfluxdbTransmits
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /InfluxdbTransmit/:id [get]
// @Success 200 {object}  servlet.JSONResult{data=models.InfluxdbTransmit} 
func (api *InfluxdbTransmitApi) ByIdInfluxdbTransmit(c *gin.Context) {
	var InfluxdbTransmit models.InfluxdbTransmit

	param := c.Param("id")

	result := glob.GDb.First(&InfluxdbTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "InfluxdbTransmit not found")

		return
	}

	servlet.Resp(c, InfluxdbTransmit)
}
