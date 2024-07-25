package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type HttpHandlerApi struct{}

var HttpHandlerBiz = biz.HttpHandlerBiz{}

// CreateHttpHandler
// @Summary 创建Http数据处理器
// @Description 创建Http数据处理器
// @Tags HttpHandlers
// @Accept json
// @Produce json
// @Param HttpHandler body models.HttpHandler true "Http数据处理器"
// @Success 201 {object} servlet.JSONResult{data=models.HttpHandler} "创建成功的Http数据处理器"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /HttpHandler/create [post]
func (api *HttpHandlerApi) CreateHttpHandler(c *gin.Context) {
	var HttpHandler models.HttpHandler
	if err := c.ShouldBindJSON(&HttpHandler); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 HttpHandler 是否被正确初始化
	if HttpHandler.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&HttpHandler)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	HttpHandlerBiz.SetRedis(HttpHandler)
	// 返回创建成功的Http数据处理器
	servlet.Resp(c, HttpHandler)
}

// UpdateHttpHandler
// @Summary 更新一个Http数据处理器
// @Description 更新一个Http数据处理器
// @Tags HttpHandlers
// @Accept json
// @Produce json
// @Param HttpHandler body models.HttpHandler true "Http数据处理器"
// @Success 200 {object}  servlet.JSONResult{data=models.HttpHandler} "Http数据处理器"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Http数据处理器未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /HttpHandler/update [post]
func (api *HttpHandlerApi) UpdateHttpHandler(c *gin.Context) {
	var req models.HttpHandler
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.HttpHandler
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "HttpHandler not found")
		return
	}

	var newV models.HttpHandler
	newV = old
	newV.Name = req.Name
	newV.Script = req.Script
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	HttpHandlerBiz.SetRedis(newV)
	servlet.Resp(c, old)
}

// PageHttpHandler
// @Summary 分页查询Http数据处理器
// @Description 分页查询Http数据处理器
// @Tags HttpHandlers
// @Accept json
// @Produce json
// @Param name query string false "Http数据处理器名称"
// @Param pid query int false "上级id"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.HttpHandler}} "Http数据处理器"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /HttpHandler/page [get]
func (api *HttpHandlerApi) PageHttpHandler(c *gin.Context) {
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

	data, err := HttpHandlerBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteHttpHandler
// @Tags      HttpHandlers
// @Summary   删除Http数据处理器
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /HttpHandler/delete/:id [post]
func (api *HttpHandlerApi) DeleteHttpHandler(c *gin.Context) {
	var HttpHandler models.HttpHandler

	param := c.Param("id")

	result := glob.GDb.First(&HttpHandler, param)
	if result.Error != nil {
		servlet.Error(c, "HttpHandler not found")

		return
	}

	if result := glob.GDb.Delete(&HttpHandler); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	HttpHandlerBiz.RemoveRedis(HttpHandler)

	servlet.Resp(c, "删除成功")
}

// ByIdHttpHandler
// @Tags      HttpHandlers
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /HttpHandler/:id [get]
func (api *HttpHandlerApi) ByIdHttpHandler(c *gin.Context) {
	var HttpHandler models.HttpHandler

	param := c.Param("id")

	result := glob.GDb.First(&HttpHandler, param)
	if result.Error != nil {
		servlet.Error(c, "HttpHandler not found")

		return
	}

	servlet.Resp(c, HttpHandler)
}



