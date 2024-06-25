package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type RepairRecordApi struct{}

var RepairRecordBiz = biz.RepairRecordBiz{}

// CreateRepairRecord
// @Summary 创建维修日志
// @Description 创建维修日志
// @Tags RepairRecords
// @Accept json
// @Produce json
// @Param RepairRecord body models.RepairRecord true "维修日志"
// @Success 201 {object} servlet.JSONResult{data=models.RepairRecord} "创建成功的维修日志"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /RepairRecord/create [post]
func (api *RepairRecordApi) CreateRepairRecord(c *gin.Context) {
	var RepairRecord models.RepairRecord
	if err := c.ShouldBindJSON(&RepairRecord); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := glob.GDb.Create(&RepairRecord)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的维修日志
	servlet.Resp(c, RepairRecord)
}

// UpdateRepairRecord
// @Summary 更新一个维修日志
// @Description 更新一个维修日志
// @Tags RepairRecords
// @Accept json
// @Produce json
// @Param RepairRecord body models.RepairRecord true "维修日志"
// @Success 200 {object}  servlet.JSONResult{data=models.RepairRecord} "维修日志"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "维修日志未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /RepairRecord/update [post]
func (api *RepairRecordApi) UpdateRepairRecord(c *gin.Context) {
	var req models.RepairRecord
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.RepairRecord
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "RepairRecord not found")
		return
	}

	var newV models.RepairRecord
	newV = old
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageRepairRecord
// @Summary 分页查询维修日志
// @Description 分页查询维修日志
// @Tags RepairRecords
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.RepairRecord}} "维修日志"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /RepairRecord/page [get]
func (api *RepairRecordApi) PageRepairRecord(c *gin.Context) {
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

	data, err := RepairRecordBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteRepairRecord
// @Tags      RepairRecords
// @Summary   删除维修日志
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /RepairRecord/delete/:id [post]
func (api *RepairRecordApi) DeleteRepairRecord(c *gin.Context) {
	var RepairRecord models.RepairRecord

	param := c.Param("id")

	result := glob.GDb.First(&RepairRecord, param)
	if result.Error != nil {
		servlet.Error(c, "RepairRecord not found")

		return
	}

	if result := glob.GDb.Delete(&RepairRecord); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdRepairRecord
// @Tags      RepairRecords
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /RepairRecord/:id [get]
func (api *RepairRecordApi) ByIdRepairRecord(c *gin.Context) {
	var RepairRecord models.RepairRecord

	param := c.Param("id")

	result := glob.GDb.First(&RepairRecord, param)
	if result.Error != nil {
		servlet.Error(c, "RepairRecord not found")

		return
	}

	servlet.Resp(c, RepairRecord)
}
