package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type DeviceInfoApi struct{}

var deviceInfoBiz = biz.DeviceInfoBiz{}

// CreateDeviceInfo
// @Summary 创建设备详情
// @Description 创建设备详情
// @Tags DeviceInfos
// @Accept json
// @Produce json
// @Param DeviceInfo body models.DeviceInfo true "设备详情"
// @Success 201 {object} servlet.JSONResult{data=models.DeviceInfo} "创建成功的设备详情"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /DeviceInfo/create [post]
func (api *DeviceInfoApi) CreateDeviceInfo(c *gin.Context) {
	var DeviceInfo models.DeviceInfo
	if err := c.ShouldBindJSON(&DeviceInfo); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 DeviceInfo 是否被正确初始化
	if DeviceInfo.SN == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	var Product models.Product
	result := glob.GDb.First(&Product, DeviceInfo.ProductId)
	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	if !DeviceInfo.ManufacturingDate.IsZero() {
		DeviceInfo.WarrantyExpiry = DeviceInfo.ManufacturingDate.AddDate(0, 0, Product.WarrantyPeriod)
	}
	result = glob.GDb.Create(&DeviceInfo)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的设备详情
	servlet.Resp(c, DeviceInfo)
}

// UpdateDeviceInfo
// @Summary 更新一个设备详情
// @Description 更新一个设备详情
// @Tags DeviceInfos
// @Accept json
// @Produce json
// @Param DeviceInfo body models.DeviceInfo true "设备详情"
// @Success 200 {object}  servlet.JSONResult{data=models.DeviceInfo} "设备详情"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "设备详情未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /DeviceInfo/update [post]
func (api *DeviceInfoApi) UpdateDeviceInfo(c *gin.Context) {
	var req models.DeviceInfo
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.DeviceInfo
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "DeviceInfo not found")
		return
	}

	var newV models.DeviceInfo
	newV = old

	var Product models.Product
	result = glob.GDb.First(&Product, newV.ProductId)
	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	if !newV.ManufacturingDate.IsZero() {
		newV.WarrantyExpiry = newV.ManufacturingDate.AddDate(0, 0, Product.WarrantyPeriod)
	}
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageDeviceInfo
// @Summary 分页查询设备详情
// @Description 分页查询设备详情
// @Tags DeviceInfos
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.DeviceInfo}} "设备详情"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /DeviceInfo/page [get]
func (api *DeviceInfoApi) PageDeviceInfo(c *gin.Context) {
	var name = c.Query("sn")
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

	data, err := deviceInfoBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteDeviceInfo
// @Tags      DeviceInfos
// @Summary   删除设备详情
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /DeviceInfo/delete/:id [post]
func (api *DeviceInfoApi) DeleteDeviceInfo(c *gin.Context) {
	var DeviceInfo models.DeviceInfo

	param := c.Param("id")

	result := glob.GDb.First(&DeviceInfo, param)
	if result.Error != nil {
		servlet.Error(c, "DeviceInfo not found")

		return
	}

	if result := glob.GDb.Delete(&DeviceInfo); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdDeviceInfo
// @Tags      DeviceInfos
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /DeviceInfo/:id [get]
func (api *DeviceInfoApi) ByIdDeviceInfo(c *gin.Context) {
	var DeviceInfo models.DeviceInfo

	param := c.Param("id")

	result := glob.GDb.First(&DeviceInfo, param)
	if result.Error != nil {
		servlet.Error(c, "DeviceInfo not found")

		return
	}

	servlet.Resp(c, DeviceInfo)
}
