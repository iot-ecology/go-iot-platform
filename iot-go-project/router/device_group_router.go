package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type DeviceGroupApi struct{}

var deviceGroupBiz = biz.DeviceGroupBiz{}

// CreateDeviceGroup
// @Summary 创建设备组
// @Description 创建设备组
// @Tags DeviceGroups
// @Accept json
// @Produce json
// @Param DeviceGroup body models.DeviceGroup true "设备组"
// @Success 201 {object} servlet.JSONResult{data=models.DeviceGroup} "创建成功的设备组"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /device_group/create [post]
func (api *DeviceGroupApi) CreateDeviceGroup(c *gin.Context) {
	var DeviceGroup models.DeviceGroup
	if err := c.ShouldBindJSON(&DeviceGroup); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 DeviceGroup 是否被正确初始化
	if DeviceGroup.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&DeviceGroup)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的设备组
	servlet.Resp(c, DeviceGroup)
}

// UpdateDeviceGroup
// @Summary 更新一个设备组
// @Description 更新一个设备组
// @Tags DeviceGroups
// @Accept json
// @Produce json
// @Param id path int true "设备组id"
// @Param DeviceGroup body models.DeviceGroup true "设备组"
// @Success 200 {object}  servlet.JSONResult{data=models.DeviceGroup} "设备组"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "设备组未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /device_group/update [post]
func (api *DeviceGroupApi) UpdateDeviceGroup(c *gin.Context) {
	var req models.DeviceGroup
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.DeviceGroup
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "DeviceGroup not found")
		return
	}

	var newV models.DeviceGroup
	newV = old
	newV.Name = req.Name
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageDeviceGroup
// @Summary 分页查询设备组
// @Description 分页查询设备组
// @Tags DeviceGroups
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.DeviceGroup}} "设备组"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /device_group/page [get]
func (api *DeviceGroupApi) PageDeviceGroup(c *gin.Context) {
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

	data, err := deviceGroupBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteDeviceGroup
// @Tags      DeviceGroups
// @Summary   删除设备组
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /device_group/delete/:id [post]
func (api *DeviceGroupApi) DeleteDeviceGroup(c *gin.Context) {
	var DeviceGroup models.DeviceGroup

	param := c.Param("id")

	result := glob.GDb.First(&DeviceGroup, param)
	if result.Error != nil {
		servlet.Error(c, "DeviceGroup not found")

		return
	}

	if result := glob.GDb.Delete(&DeviceGroup); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdDeviceGroup
// @Tags      DeviceGroups
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /device_group/:id [get]
func (api *DeviceGroupApi) ByIdDeviceGroup(c *gin.Context) {
	var DeviceGroup models.DeviceGroup

	param := c.Param("id")

	result := glob.GDb.First(&DeviceGroup, param)
	if result.Error != nil {
		servlet.Error(c, "DeviceGroup not found")

		return
	}

	servlet.Resp(c, DeviceGroup)
}
