package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type ProductionPlanApi struct{}

var ProductionPlanBiz = biz.ProductionPlanBiz{}

// CreateProductionPlan
// @Summary 创建生产计划
// @Description 创建生产计划
// @Tags ProductionPlans
// @Accept json
// @Produce json
// @Param ProductionPlan body models.ProductionPlan true "生产计划"
// @Success 201 {object} servlet.JSONResult{data=models.ProductionPlan} "创建成功的生产计划"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ProductionPlan/create [post]
func (api *ProductionPlanApi) CreateProductionPlan(c *gin.Context) {
	var ProductionPlan models.ProductionPlan
	if err := c.ShouldBindJSON(&ProductionPlan); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 ProductionPlan 是否被正确初始化
	if ProductionPlan.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&ProductionPlan)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的生产计划
	servlet.Resp(c, ProductionPlan)
}

// UpdateProductionPlan
// @Summary 更新一个生产计划
// @Description 更新一个生产计划
// @Tags ProductionPlans
// @Accept json
// @Produce json
// @Param id path int true "生产计划id"
// @Param ProductionPlan body models.ProductionPlan true "生产计划"
// @Success 200 {object}  servlet.JSONResult{data=models.ProductionPlan} "生产计划"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "生产计划未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ProductionPlan/update [post]
func (api *ProductionPlanApi) UpdateProductionPlan(c *gin.Context) {
	var req models.ProductionPlan
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.ProductionPlan
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "ProductionPlan not found")
		return
	}

	var newV models.ProductionPlan
	newV = old
	newV.Name = req.Name
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageProductionPlan
// @Summary 分页查询生产计划
// @Description 分页查询生产计划
// @Tags ProductionPlans
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.ProductionPlan}} "生产计划"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /ProductionPlan/page [get]
func (api *ProductionPlanApi) PageProductionPlan(c *gin.Context) {
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

	data, err := ProductionPlanBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteProductionPlan
// @Tags      ProductionPlans
// @Summary   删除生产计划
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /ProductionPlan/delete/:id [post]
func (api *ProductionPlanApi) DeleteProductionPlan(c *gin.Context) {
	var ProductionPlan models.ProductionPlan

	param := c.Param("id")

	result := glob.GDb.First(&ProductionPlan, param)
	if result.Error != nil {
		servlet.Error(c, "ProductionPlan not found")

		return
	}

	if result := glob.GDb.Delete(&ProductionPlan); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdProductionPlan
// @Tags      ProductionPlans
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /ProductionPlan/:id [get]
func (api *ProductionPlanApi) ByIdProductionPlan(c *gin.Context) {
	var ProductionPlan models.ProductionPlan

	param := c.Param("id")

	result := glob.GDb.First(&ProductionPlan, param)
	if result.Error != nil {
		servlet.Error(c, "ProductionPlan not found")

		return
	}

	servlet.Resp(c, ProductionPlan)
}
