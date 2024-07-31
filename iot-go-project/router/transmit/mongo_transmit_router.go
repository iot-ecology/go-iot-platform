package transmit

import (
	"github.com/gin-gonic/gin"
	"igp/biz/transmit"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type MongoTransmitApi struct{}

var MongoTransmitBiz = transmit.MongoTransmitBiz{}

// CreateMongoTransmit
// @Summary 创建Mongo数据库管理
// @Description 创建Mongo数据库管理
// @Tags MongoTransmits
// @Accept json
// @Produce json
// @Param MongoTransmit body models.MongoTransmit true "Mongo数据库管理"
// @Success 201 {object} servlet.JSONResult{data=models.MongoTransmit} "创建成功的Mongo数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /MongoTransmit/create [post]
func (api *MongoTransmitApi) CreateMongoTransmit(c *gin.Context) {
	var MongoTransmit models.MongoTransmit
	if err := c.ShouldBindJSON(&MongoTransmit); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := glob.GDb.Create(&MongoTransmit)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	MongoTransmitBiz.SetRedis(MongoTransmit)
	// 返回创建成功的Mongo数据库管理
	servlet.Resp(c, MongoTransmit)
}

// UpdateMongoTransmit
// @Summary 更新一个Mongo数据库管理
// @Description 更新一个Mongo数据库管理
// @Tags MongoTransmits
// @Accept json
// @Produce json
// @Param MongoTransmit body models.MongoTransmit true "Mongo数据库管理"
// @Success 200 {object}  servlet.JSONResult{data=models.MongoTransmit} "Mongo数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Mongo数据库管理未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /MongoTransmit/update [post]
func (api *MongoTransmitApi) UpdateMongoTransmit(c *gin.Context) {
	var req models.MongoTransmit
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.MongoTransmit
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "MongoTransmit not found")
		return
	}

	var newV models.MongoTransmit
	newV = old

	newV.Name = req.Name
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	MongoTransmitBiz.SetRedis(newV)

	servlet.Resp(c, old)
}

// PageMongoTransmit
// @Summary 分页查询Mongo数据库管理
// @Description 分页查询Mongo数据库管理
// @Tags MongoTransmits
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.MongoTransmit}} "Mongo数据库管理"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /MongoTransmit/page [get]
func (api *MongoTransmitApi) PageMongoTransmit(c *gin.Context) {
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

	data, err := MongoTransmitBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteMongoTransmit
// @Tags      MongoTransmits
// @Summary   删除Mongo数据库管理
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /MongoTransmit/delete/:id [post]
func (api *MongoTransmitApi) DeleteMongoTransmit(c *gin.Context) {
	var MongoTransmit models.MongoTransmit

	param := c.Param("id")

	result := glob.GDb.First(&MongoTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "MongoTransmit not found")

		return
	}

	if result := glob.GDb.Delete(&MongoTransmit); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	MongoTransmitBiz.DeleteRedis(MongoTransmit)
	servlet.Resp(c, "删除成功")
}

// ByIdMongoTransmit
// @Tags      MongoTransmits
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /MongoTransmit/:id [get]
func (api *MongoTransmitApi) ByIdMongoTransmit(c *gin.Context) {
	var MongoTransmit models.MongoTransmit

	param := c.Param("id")

	result := glob.GDb.First(&MongoTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "MongoTransmit not found")

		return
	}

	servlet.Resp(c, MongoTransmit)
}