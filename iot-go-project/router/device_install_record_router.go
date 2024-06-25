package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type DeviceInstallRecordApi struct{}

var installRecordBiz = biz.DeviceInstallRecordBiz{}

// CreateDeviceInstallRecord
// @Summary 创建安装记录
// @Description 创建安装记录
// @Tags DeviceInstallRecords
// @Accept json
// @Produce json
// @Param DeviceInstallRecord body models.DeviceInstallRecord true "安装记录"
// @Success 201 {object} servlet.JSONResult{data=models.DeviceInstallRecord} "创建成功的安装记录"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /DeviceInstallRecord/create [post]
func (api *DeviceInstallRecordApi) CreateDeviceInstallRecord(c *gin.Context) {
	var DeviceInstallRecord models.DeviceInstallRecord
	if err := c.ShouldBindJSON(&DeviceInstallRecord); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := glob.GDb.Create(&DeviceInstallRecord)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的安装记录
	servlet.Resp(c, DeviceInstallRecord)
}

// UpdateDeviceInstallRecord
// @Summary 更新一个安装记录
// @Description 更新一个安装记录
// @Tags DeviceInstallRecords
// @Accept json
// @Produce json
// @Param DeviceInstallRecord body models.DeviceInstallRecord true "安装记录"
// @Success 200 {object}  servlet.JSONResult{data=models.DeviceInstallRecord} "安装记录"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "安装记录未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /DeviceInstallRecord/update [post]
func (api *DeviceInstallRecordApi) UpdateDeviceInstallRecord(c *gin.Context) {
	var req models.DeviceInstallRecord
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.DeviceInstallRecord
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "DeviceInstallRecord not found")
		return
	}

	var newV models.DeviceInstallRecord
	newV = old

	newV.InstallDate = req.InstallDate
	newV.Technician = req.Technician
	newV.Description = req.Description
	newV.PhotoURL = req.PhotoURL
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageDeviceInstallRecord
// @Summary 分页查询安装记录
// @Description 分页查询安装记录
// @Tags DeviceInstallRecords
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.DeviceInstallRecord}} "安装记录"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /DeviceInstallRecord/page [get]
func (api *DeviceInstallRecordApi) PageDeviceInstallRecord(c *gin.Context) {
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

	data, err := installRecordBiz.PageData("", parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteDeviceInstallRecord
// @Tags      DeviceInstallRecords
// @Summary   删除安装记录
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /DeviceInstallRecord/delete/:id [post]
func (api *DeviceInstallRecordApi) DeleteDeviceInstallRecord(c *gin.Context) {
	var DeviceInstallRecord models.DeviceInstallRecord

	param := c.Param("id")

	result := glob.GDb.First(&DeviceInstallRecord, param)
	if result.Error != nil {
		servlet.Error(c, "DeviceInstallRecord not found")

		return
	}

	if result := glob.GDb.Delete(&DeviceInstallRecord); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdDeviceInstallRecord
// @Tags      DeviceInstallRecords
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /DeviceInstallRecord/:id [get]
func (api *DeviceInstallRecordApi) ByIdDeviceInstallRecord(c *gin.Context) {
	var DeviceInstallRecord models.DeviceInstallRecord

	param := c.Param("id")

	result := glob.GDb.First(&DeviceInstallRecord, param)
	if result.Error != nil {
		servlet.Error(c, "DeviceInstallRecord not found")

		return
	}

	servlet.Resp(c, DeviceInstallRecord)
}
