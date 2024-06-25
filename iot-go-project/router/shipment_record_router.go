package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
	"time"
)

type ShipmentRecordApi struct{}

var shipmentRecordBiz = biz.ShipmentRecordBiz{}

// CreateShipmentRecord
// @Summary 创建发货记录
// @Description 创建发货记录
// @Tags ShipmentRecords
// @Accept json
// @Produce json
// @Param ShipmentRecord body models.ShipmentRecord true "发货记录"
// @Success 201 {object} servlet.JSONResult{data=models.ShipmentRecord} "创建成功的发货记录"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ShipmentRecord/create [post]
func (api *ShipmentRecordApi) CreateShipmentRecord(c *gin.Context) {
	var param servlet.ShipmentRecordCreateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	tx := glob.GDb.Begin()
	if tx.Error != nil {
		servlet.Error(c, "Failed to begin transaction")
		return
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		servlet.Error(c, "Failed to generate uuid")
		tx.Rollback()
		return
	}
	var shipmentRecord = models.ShipmentRecord{
		ShipmentDate:    param.ShipmentDate,
		Technician:      param.Technician,
		CustomerName:    param.CustomerName,
		CustomerPhone:   param.CustomerPhone,
		CustomerAddress: param.CustomerAddress,
		TrackingNumber:  newUUID.String(),
		Status:          param.Status,
		Description:     param.Description,
	}

	create := tx.Model(models.ShipmentRecord{}).Create(shipmentRecord)
	if create.Error != nil {
		tx.Rollback()
		zap.S().Errorf("创建 ShipmentRecord 异常 %+v", create.Error)
		servlet.Error(c, create.Error.Error())
		return
	}

	var shipmentProductDetail []models.ShipmentProductDetail
	for _, createParam := range param.ProductPlanCreateParams {
		shipmentProductDetail = append(shipmentProductDetail, models.ShipmentProductDetail{
			ShipmentRecordId: shipmentRecord.ID,
			ProductID:        createParam.ProductID,
			Quantity:         createParam.Quantity,
		})
	}

	result := tx.Model(&models.ShipmentProductDetail{}).CreateInBatches(shipmentProductDetail, len(shipmentProductDetail))

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

	servlet.Resp(c, "创建成功")
}

// UpdateShipmentRecord
// @Summary 更新一个发货记录
// @Description 更新一个发货记录
// @Tags ShipmentRecords
// @Accept json
// @Produce json
// @Param ShipmentRecord body models.ShipmentRecord true "发货记录"
// @Success 200 {object}  servlet.JSONResult{data=models.ShipmentRecord} "发货记录"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "发货记录未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ShipmentRecord/update [post]
func (api *ShipmentRecordApi) UpdateShipmentRecord(c *gin.Context) {
	var param servlet.ShipmentRecordCreateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	var old models.ShipmentRecord
	result := glob.GDb.First(&old, param.ID)
	if result.Error != nil {
		servlet.Error(c, "ProductionPlan not found")
		return
	}

	currentTime := time.Now()

	if old.ShipmentDate.Before(currentTime) {
		servlet.Error(c, "发货小于当前时间不允许修改")
		return
	}

	tx := glob.GDb.Begin()
	if tx.Error != nil {
		servlet.Error(c, "Failed to begin transaction")
		return
	}

	var shipmentRecord models.ShipmentRecord
	shipmentRecord.ShipmentDate = param.ShipmentDate
	shipmentRecord.Technician = param.Technician
	shipmentRecord.CustomerName = param.CustomerName
	shipmentRecord.CustomerPhone = param.CustomerPhone
	shipmentRecord.CustomerAddress = param.CustomerAddress
	shipmentRecord.Status = param.Status
	shipmentRecord.Description = param.Description
	create := tx.Model(models.ShipmentRecord{}).Updates(shipmentRecord)
	if create.Error != nil {
		tx.Rollback()
		zap.S().Errorf("更新 ShipmentRecord 异常 %+v", create.Error)
		servlet.Error(c, create.Error.Error())
		return
	}

	db := tx.Model(&models.ShipmentProductDetail{}).Where("shipment_record_id = ?", shipmentRecord.ID).Delete(models.ProductPlan{})
	if db.Error != nil {
		tx.Rollback()
		zap.S().Infoln("Error occurred during deletion:", db.Error)
		servlet.Error(c, "Error occurred during deletion")
		return
	}

	var shipmentProductDetail []models.ShipmentProductDetail
	for _, createParam := range param.ProductPlanCreateParams {
		shipmentProductDetail = append(shipmentProductDetail, models.ShipmentProductDetail{
			ShipmentRecordId: shipmentRecord.ID,
			ProductID:        createParam.ProductID,
			Quantity:         createParam.Quantity,
		})
	}

	result = tx.Model(&models.ShipmentProductDetail{}).CreateInBatches(shipmentProductDetail, len(shipmentProductDetail))

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

	servlet.Resp(c, "更新成功")
}

// PageShipmentRecord
// @Summary 分页查询发货记录
// @Description 分页查询发货记录
// @Tags ShipmentRecords
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.ShipmentRecord}} "发货记录"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /ShipmentRecord/page [get]
func (api *ShipmentRecordApi) PageShipmentRecord(c *gin.Context) {
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

	data, err := shipmentRecordBiz.PageData("", parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteShipmentRecord
// @Tags      ShipmentRecords
// @Summary   删除发货记录
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /ShipmentRecord/delete/:id [post]
func (api *ShipmentRecordApi) DeleteShipmentRecord(c *gin.Context) {
	var ShipmentRecord models.ShipmentRecord

	param := c.Param("id")

	result := glob.GDb.First(&ShipmentRecord, param)
	if result.Error != nil {
		servlet.Error(c, "ShipmentRecord not found")

		return
	}

	if result := glob.GDb.Delete(&ShipmentRecord); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdShipmentRecord
// @Tags      ShipmentRecords
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /ShipmentRecord/:id [get]
func (api *ShipmentRecordApi) ByIdShipmentRecord(c *gin.Context) {
	var ShipmentRecord models.ShipmentRecord

	param := c.Param("id")

	result := glob.GDb.First(&ShipmentRecord, param)
	if result.Error != nil {
		servlet.Error(c, "ShipmentRecord not found")

		return
	}

	servlet.Resp(c, ShipmentRecord)
}
