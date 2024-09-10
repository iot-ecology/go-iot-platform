package router

import (
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WebsocketHandlerApi struct{}

var WebsocketHandlerBiz = biz.WebsocketHandlerBiz{}

// CreateWebsocketHandler
// @Summary 创建Websocket数据处理器
// @Description 创建Websocket数据处理器
// @Tags WebsocketHandlers
// @Accept json
// @Produce json
// @Param WebsocketHandler body models.WebsocketHandler true "Websocket数据处理器"
// @Success 200 {object} servlet.JSONResult{data=models.WebsocketHandler} "创建成功的Websocket数据处理器"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /WebsocketHandler/create [post]
func (api *WebsocketHandlerApi) CreateWebsocketHandler(c *gin.Context) {
	var WebsocketHandler models.WebsocketHandler
	if err := c.ShouldBindJSON(&WebsocketHandler); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 WebsocketHandler 是否被正确初始化
	if WebsocketHandler.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&WebsocketHandler)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	WebsocketHandlerBiz.SetRedis(WebsocketHandler)
	// 返回创建成功的Websocket数据处理器
	servlet.Resp(c, WebsocketHandler)
}

// UpdateWebsocketHandler
// @Summary 更新一个Websocket数据处理器
// @Description 更新一个Websocket数据处理器
// @Tags WebsocketHandlers
// @Accept json
// @Produce json
// @Param WebsocketHandler body models.WebsocketHandler true "Websocket数据处理器"
// @Success 200 {object}  servlet.JSONResult{data=models.WebsocketHandler} "Websocket数据处理器"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Websocket数据处理器未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /WebsocketHandler/update [post]
func (api *WebsocketHandlerApi) UpdateWebsocketHandler(c *gin.Context) {
	var req models.WebsocketHandler
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.WebsocketHandler
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "WebsocketHandler not found")
		return
	}

	var newV models.WebsocketHandler
	newV = old
	newV.Name = req.Name
	newV.Script = req.Script
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	WebsocketHandlerBiz.SetRedis(newV)
	servlet.Resp(c, old)
}

// PageWebsocketHandler
// @Summary 分页查询Websocket数据处理器
// @Description 分页查询Websocket数据处理器
// @Tags WebsocketHandlers
// @Accept json
// @Produce json
// @Param name query string false "Websocket数据处理器名称"
// @Param pid query int false "上级id"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.WebsocketHandler}} "Websocket数据处理器"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /WebsocketHandler/page [get]
func (api *WebsocketHandlerApi) PageWebsocketHandler(c *gin.Context) {
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

	data, err := WebsocketHandlerBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteWebsocketHandler
// @Tags      WebsocketHandlers
// @Summary   删除Websocket数据处理器
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /WebsocketHandler/delete/:id [post]
// @Success 200 {object}  servlet.JSONResult{data=string}
func (api *WebsocketHandlerApi) DeleteWebsocketHandler(c *gin.Context) {
	var WebsocketHandler models.WebsocketHandler

	param := c.Param("id")

	result := glob.GDb.First(&WebsocketHandler, param)
	if result.Error != nil {
		servlet.Error(c, "WebsocketHandler not found")

		return
	}

	if result := glob.GDb.Delete(&WebsocketHandler); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	WebsocketHandlerBiz.RemoveRedis(WebsocketHandler)

	servlet.Resp(c, "删除成功")
}

// ByIdWebsocketHandler
// @Tags      WebsocketHandlers
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /WebsocketHandler/:id [get]
// @Success 200 {object}  servlet.JSONResult{data=models.WebsocketHandler}
func (api *WebsocketHandlerApi) ByIdWebsocketHandler(c *gin.Context) {
	var WebsocketHandler models.WebsocketHandler

	param := c.Param("id")

	result := glob.GDb.First(&WebsocketHandler, param)
	if result.Error != nil {
		servlet.Error(c, "WebsocketHandler not found")

		return
	}

	servlet.Resp(c, WebsocketHandler)
}
