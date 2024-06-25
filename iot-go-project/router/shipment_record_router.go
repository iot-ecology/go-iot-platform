package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
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
	var ShipmentRecord models.ShipmentRecord
	if err := c.ShouldBindJSON(&ShipmentRecord); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := glob.GDb.Create(&ShipmentRecord)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的发货记录
	servlet.Resp(c, ShipmentRecord)
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
	var req models.ShipmentRecord
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.ShipmentRecord
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "ShipmentRecord not found")
		return
	}

	var newV models.ShipmentRecord
	newV = old
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
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
