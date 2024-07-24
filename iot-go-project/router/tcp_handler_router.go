package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type TcpHandlerApi struct{}

var TcpHandlerBiz = biz.TcpHandlerBiz{}

// CreateTcpHandler
// @Summary 创建Tcp数据处理器
// @Description 创建Tcp数据处理器
// @Tags TcpHandlers
// @Accept json
// @Produce json
// @Param TcpHandler body models.TcpHandler true "Tcp数据处理器"
// @Success 201 {object} servlet.JSONResult{data=models.TcpHandler} "创建成功的Tcp数据处理器"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /TcpHandler/create [post]
func (api *TcpHandlerApi) CreateTcpHandler(c *gin.Context) {
	var TcpHandler models.TcpHandler
	if err := c.ShouldBindJSON(&TcpHandler); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 TcpHandler 是否被正确初始化
	if TcpHandler.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&TcpHandler)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	TcpHandlerBiz.SetRedis(TcpHandler)
	// 返回创建成功的Tcp数据处理器
	servlet.Resp(c, TcpHandler)
}

// UpdateTcpHandler
// @Summary 更新一个Tcp数据处理器
// @Description 更新一个Tcp数据处理器
// @Tags TcpHandlers
// @Accept json
// @Produce json
// @Param TcpHandler body models.TcpHandler true "Tcp数据处理器"
// @Success 200 {object}  servlet.JSONResult{data=models.TcpHandler} "Tcp数据处理器"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Tcp数据处理器未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /TcpHandler/update [post]
func (api *TcpHandlerApi) UpdateTcpHandler(c *gin.Context) {
	var req models.TcpHandler
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.TcpHandler
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "TcpHandler not found")
		return
	}

	var newV models.TcpHandler
	newV = old
	newV.Name = req.Name
	newV.Script = req.Script
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	TcpHandlerBiz.SetRedis(newV)
	servlet.Resp(c, old)
}

// PageTcpHandler
// @Summary 分页查询Tcp数据处理器
// @Description 分页查询Tcp数据处理器
// @Tags TcpHandlers
// @Accept json
// @Produce json
// @Param name query string false "Tcp数据处理器名称"
// @Param pid query int false "上级id"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.TcpHandler}} "Tcp数据处理器"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /TcpHandler/page [get]
func (api *TcpHandlerApi) PageTcpHandler(c *gin.Context) {
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

	data, err := TcpHandlerBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteTcpHandler
// @Tags      TcpHandlers
// @Summary   删除Tcp数据处理器
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /TcpHandler/delete/:id [post]
func (api *TcpHandlerApi) DeleteTcpHandler(c *gin.Context) {
	var TcpHandler models.TcpHandler

	param := c.Param("id")

	result := glob.GDb.First(&TcpHandler, param)
	if result.Error != nil {
		servlet.Error(c, "TcpHandler not found")

		return
	}

	if result := glob.GDb.Delete(&TcpHandler); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	TcpHandlerBiz.RemoveRedis(TcpHandler)

	servlet.Resp(c, "删除成功")
}

// ByIdTcpHandler
// @Tags      TcpHandlers
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /TcpHandler/:id [get]
func (api *TcpHandlerApi) ByIdTcpHandler(c *gin.Context) {
	var TcpHandler models.TcpHandler

	param := c.Param("id")

	result := glob.GDb.First(&TcpHandler, param)
	if result.Error != nil {
		servlet.Error(c, "TcpHandler not found")

		return
	}

	servlet.Resp(c, TcpHandler)
}



