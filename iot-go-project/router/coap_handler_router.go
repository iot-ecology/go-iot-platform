package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type CoapHandlerApi struct{}

var CoapHandlerBiz = biz.CoapHandlerBiz{}

// CreateCoapHandler
// @Summary 创建Coap数据处理器
// @Description 创建Coap数据处理器
// @Tags CoapHandlers
// @Accept json
// @Produce json
// @Param CoapHandler body models.CoapHandler true "Coap数据处理器"
// @Success 201 {object} servlet.JSONResult{data=models.CoapHandler} "创建成功的Coap数据处理器"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /CoapHandler/create [post]
func (api *CoapHandlerApi) CreateCoapHandler(c *gin.Context) {
	var CoapHandler models.CoapHandler
	if err := c.ShouldBindJSON(&CoapHandler); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 CoapHandler 是否被正确初始化
	if CoapHandler.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&CoapHandler)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	CoapHandlerBiz.SetRedis(CoapHandler)
	// 返回创建成功的Coap数据处理器
	servlet.Resp(c, CoapHandler)
}

// UpdateCoapHandler
// @Summary 更新一个Coap数据处理器
// @Description 更新一个Coap数据处理器
// @Tags CoapHandlers
// @Accept json
// @Produce json
// @Param CoapHandler body models.CoapHandler true "Coap数据处理器"
// @Success 200 {object}  servlet.JSONResult{data=models.CoapHandler} "Coap数据处理器"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Coap数据处理器未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /CoapHandler/update [post]
func (api *CoapHandlerApi) UpdateCoapHandler(c *gin.Context) {
	var req models.CoapHandler
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.CoapHandler
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "CoapHandler not found")
		return
	}

	var newV models.CoapHandler
	newV = old
	newV.Name = req.Name
	newV.Script = req.Script
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	CoapHandlerBiz.SetRedis(newV)
	servlet.Resp(c, old)
}

// PageCoapHandler
// @Summary 分页查询Coap数据处理器
// @Description 分页查询Coap数据处理器
// @Tags CoapHandlers
// @Accept json
// @Produce json
// @Param name query string false "Coap数据处理器名称"
// @Param pid query int false "上级id"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.CoapHandler}} "Coap数据处理器"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /CoapHandler/page [get]
func (api *CoapHandlerApi) PageCoapHandler(c *gin.Context) {
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

	data, err := CoapHandlerBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteCoapHandler
// @Tags      CoapHandlers
// @Summary   删除Coap数据处理器
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /CoapHandler/delete/:id [post]
func (api *CoapHandlerApi) DeleteCoapHandler(c *gin.Context) {
	var CoapHandler models.CoapHandler

	param := c.Param("id")

	result := glob.GDb.First(&CoapHandler, param)
	if result.Error != nil {
		servlet.Error(c, "CoapHandler not found")

		return
	}

	if result := glob.GDb.Delete(&CoapHandler); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	CoapHandlerBiz.RemoveRedis(CoapHandler)

	servlet.Resp(c, "删除成功")
}

// ByIdCoapHandler
// @Tags      CoapHandlers
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /CoapHandler/:id [get]
func (api *CoapHandlerApi) ByIdCoapHandler(c *gin.Context) {
	var CoapHandler models.CoapHandler

	param := c.Param("id")

	result := glob.GDb.First(&CoapHandler, param)
	if result.Error != nil {
		servlet.Error(c, "CoapHandler not found")

		return
	}

	servlet.Resp(c, CoapHandler)
}



