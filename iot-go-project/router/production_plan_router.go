package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
	"time"
)

type ProductionPlanApi struct{}

var ProductionPlanBiz = biz.ProductionPlanBiz{}

// CreateProductionPlan
// @Summary 创建生产计划
// @Description 创建生产计划
// @Tags ProductionPlans
// @Accept json
// @Produce json
// @Param ProductionPlan body servlet.ProductionPlanCreateParam true "生产计划"
// @Success 201 {object} servlet.JSONResult{data=models.ProductionPlan} "创建成功的生产计划"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ProductionPlan/create [post]
func (api *ProductionPlanApi) CreateProductionPlan(c *gin.Context) {
	var param servlet.ProductionPlanCreateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	e := param.ValidateUniqueIDs()
	if e != nil {
		servlet.Error(c, e.Error())
		return
	}
	tx := glob.GDb.Begin()
	if tx.Error != nil {
		servlet.Error(c, "Failed to begin transaction")
		return
	}

	var productionPlan models.ProductionPlan

	productionPlan.Name = param.Name
	productionPlan.Description = param.Description
	productionPlan.StartDate = param.StartDate
	productionPlan.EndDate = param.EndDate

	create := tx.Model(models.ProductionPlan{}).Create(productionPlan)
	if create.Error != nil {
		tx.Rollback()
		zap.S().Errorf("创建 ProductionPlan 异常 %+v", create.Error)
		servlet.Error(c, create.Error.Error())
		return
	}

	var productPlans []models.ProductPlan
	for _, createParam := range param.ProductPlanCreateParams {
		productPlans = append(productPlans, models.ProductPlan{
			ProductionPlanID: productionPlan.ID,
			ProductID:        createParam.ProductID,
			Quantity:         createParam.Quantity,
		})
	}

	result := tx.Model(&models.ProductPlan{}).CreateInBatches(productPlans, len(productPlans))

	if result.Error != nil {
		tx.Rollback()
		zap.S().Infoln("Error occurred during creation:", result.Error)
		servlet.Error(c, "Error occurred during creation")
		return
	}
	if err := tx.Commit().Error; err != nil {
		servlet.Error(c, "Failed to commit transaction")
		return
	}

	ProductionPlanBiz.BeforeCreate(productionPlan)
	servlet.Resp(c, productionPlan)
}

// UpdateProductionPlan
// @Summary 更新一个生产计划
// @Description 更新一个生产计划
// @Tags ProductionPlans
// @Accept json
// @Produce json
// @Param ProductionPlan body models.ProductionPlan true "生产计划"
// @Success 200 {object}  servlet.JSONResult{data=models.ProductionPlan} "生产计划"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "生产计划未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ProductionPlan/update [post]
func (api *ProductionPlanApi) UpdateProductionPlan(c *gin.Context) {

	var param servlet.ProductionPlanCreateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	e := param.ValidateUniqueIDs()
	if e != nil {
		servlet.Error(c, e.Error())
		return
	}

	var old models.ProductionPlan
	result := glob.GDb.First(&old, param.ID)
	if result.Error != nil {
		servlet.Error(c, "ProductionPlan not found")
		return
	}

	currentTime := time.Now()

	if old.StartDate.Before(currentTime) {
		servlet.Error(c, "开始时间小于当前时间不允许修改")
		return
	}

	tx := glob.GDb.Begin()
	if tx.Error != nil {
		servlet.Error(c, "Failed to begin transaction")
		return
	}

	var productionPlan models.ProductionPlan
	productionPlan = old

	productionPlan.Name = param.Name
	productionPlan.Description = param.Description
	productionPlan.StartDate = param.StartDate
	productionPlan.EndDate = param.EndDate

	create := tx.Model(models.ProductionPlan{}).Updates(productionPlan)
	if create.Error != nil {
		tx.Rollback()
		zap.S().Errorf("更新 ProductionPlan 异常 %+v", create.Error)
		servlet.Error(c, create.Error.Error())
		return
	}

	db := tx.Model(&models.ProductPlan{}).Where("production_plan_id = ?", productionPlan.ID).Delete(models.ProductPlan{})
	if db.Error != nil {
		tx.Rollback()
		zap.S().Infoln("Error occurred during deletion:", db.Error)
		servlet.Error(c, "Error occurred during deletion")
		return
	}

	var productPlans []models.ProductPlan
	for _, createParam := range param.ProductPlanCreateParams {
		productPlans = append(productPlans, models.ProductPlan{
			ProductionPlanID: productionPlan.ID,
			ProductID:        createParam.ProductID,
			Quantity:         createParam.Quantity,
		})
	}

	result = tx.Model(&models.ProductPlan{}).CreateInBatches(productPlans, len(productPlans))

	if result.Error != nil {
		tx.Rollback()
		zap.S().Infoln("Error occurred during creation:", result.Error)
		servlet.Error(c, "Error occurred during creation")
		return
	}
	if err := tx.Commit().Error; err != nil {
		servlet.Error(c, "Failed to commit transaction")
		return
	}

	servlet.Resp(c, "绑定成功")

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
// @Success 200 {object} servlet.JSONResult{data=servlet.ProductionPlanCreateParam} "生产计划"
func (api *ProductionPlanApi) ByIdProductionPlan(c *gin.Context) {
	var ProductionPlan models.ProductionPlan

	param := c.Param("id")

	result := glob.GDb.First(&ProductionPlan, param)
	if result.Error != nil {
		servlet.Error(c, "ProductionPlan not found")

		return
	}

	var res servlet.ProductionPlanCreateParam
	res.ID = ProductionPlan.ID
	res.Name = ProductionPlan.Name
	res.StartDate = ProductionPlan.StartDate
	res.EndDate = ProductionPlan.EndDate
	res.Description = ProductionPlan.Description
	var dt []models.ProductPlan

	tx := glob.GDb.Model(models.ProductPlan{}).Where("production_plan_id = ?", ProductionPlan.ID).Find(&dt)
	if tx.Error != nil {
		servlet.Error(c, tx.Error.Error())
		return
	}

	res.ProductPlanCreateParams = make([]servlet.ProductPlanCreateParam, len(dt))
	for i, v := range dt {
		res.ProductPlanCreateParams[i].ProductID = v.ProductID
		res.ProductPlanCreateParams[i].Quantity = v.Quantity
	}

	servlet.Resp(c, res)
}

// ChangeProductionPlanState
// @Tags      ProductionPlans
// @Summary   修改生产计划状态
// @Param ProductionPlan body servlet.ProductionPlanChangeParam true "修改参数"
// @Produce   application/json
// @Router    /ProductionPlan/change_state [post]
func (api *ProductionPlanApi) ChangeProductionPlanState(c *gin.Context) {
	var param servlet.ProductionPlanChangeParam
	if err := c.ShouldBindJSON(&param); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	state := ProductionPlanBiz.ChangeProductionPlanState(param)
	if state {

		servlet.Resp(c, "修改成功")
	} else {
		servlet.Error(c, "修改失败")
	}
}
