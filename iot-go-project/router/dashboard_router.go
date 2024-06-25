package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type DashboardApi struct{}

var dashBiz = biz.DashboardBiz{}

// CreateDashboard
// @Summary 创建面板
// @Description 创建面板
// @Tags dashboards
// @Accept json
// @Produce json
// @Param dashboard body models.Dashboard true "面板"
// @Success 201 {object} servlet.JSONResult{data=models.Dashboard} "创建成功的面板"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /dashboard/create [post]
func (api *DashboardApi) CreateDashboard(c *gin.Context) {
	var dashboard models.Dashboard
	if err := c.ShouldBindJSON(&dashboard); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 dashboard 是否被正确初始化
	if dashboard.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&dashboard)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的面板
	servlet.Resp(c, dashboard)
}

// UpdateDashboard
// @Summary 更新一个面板
// @Description 更新一个面板
// @Tags dashboards
// @Accept json
// @Produce json
// @Param dashboard body models.Dashboard true "面板"
// @Success 200 {object}  servlet.JSONResult{data=models.Dashboard} "面板"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "面板未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /dashboard/update [post]
func (api *DashboardApi) UpdateDashboard(c *gin.Context) {
	var req models.Dashboard
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.Dashboard
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "dashboard not found")
		return
	}

	var newV models.Dashboard
	newV = old
	newV.Name = req.Name
	newV.Config = req.Config
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageDashboard
// @Summary 分页查询面板
// @Description 分页查询面板
// @Tags dashboards
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.Dashboard}} "面板"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /dashboard/page [get]
func (api *DashboardApi) PageDashboard(c *gin.Context) {
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

	data, err := dashBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteDashboard
// @Tags      dashboards
// @Summary   删除面板
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /dashboard/delete/:id [post]
func (api *DashboardApi) DeleteDashboard(c *gin.Context) {
	var dashboard models.Dashboard

	param := c.Param("id")

	result := glob.GDb.First(&dashboard, param)
	if result.Error != nil {
		servlet.Error(c, "dashboard not found")

		return
	}

	if result := glob.GDb.Delete(&dashboard); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdDashboard
// @Tags      dashboards
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /dashboard/:id [get]
func (api *DashboardApi) ByIdDashboard(c *gin.Context) {
	var dashboard models.Dashboard

	param := c.Param("id")

	result := glob.GDb.First(&dashboard, param)
	if result.Error != nil {
		servlet.Error(c, "dashboard not found")

		return
	}

	servlet.Resp(c, dashboard)
}
