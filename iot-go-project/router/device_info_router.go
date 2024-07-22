package router

import (
	"context"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		WarrantyExpiry := DeviceInfo.ManufacturingDate.AddDate(0, 0, Product.WarrantyPeriod)
		DeviceInfo.WarrantyExpiry = &WarrantyExpiry
	}

	m := structs.Map(DeviceInfo)

	result = glob.GDb.Model(models.DeviceInfo{}).Create(m)

	if result.Error != nil {
		zap.S().Errorw("创建 DeviceInfo 失败", "error", result.Error)
		servlet.Error(c, result.Error.Error())
		return
	}
	deviceInfoBiz.SetRedis(DeviceInfo)

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
		WarrantyExpiry := newV.ManufacturingDate.AddDate(0, 0, Product.WarrantyPeriod)
		newV.WarrantyExpiry = &WarrantyExpiry
	}
	result = glob.GDb.Model(&newV).Updates(newV)

	deviceInfoBiz.SetRedis(newV)
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
// @Param sn query string false "SN"
// @Param manufacturingDateStart query string false "制造日期开始"
// @Param manufacturingDateEnd query string false "制造日期结束"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.DeviceInfo}} "设备详情"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /DeviceInfo/page [get]
func (api *DeviceInfoApi) PageDeviceInfo(c *gin.Context) {
	var sn = c.Query("sn")
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

	data, err := deviceInfoBiz.PageData(sn, parseUint, u)
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

	deviceInfoBiz.RemoveRedis(DeviceInfo.ID)
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

// QueryBindMqtt
// @Tags      DeviceInfos
// @Summary   查询绑定MQTT客户端
// @Accept json
// @Produce json
// @Param device_info_id path int true "主键"
// @Router    /DeviceInfo/QueryBindMqtt [get]
func (api *DeviceInfoApi) QueryBindMqtt(c *gin.Context) {
	param := c.Param("device_info_id")

	var deviceBindMqttClients []models.DeviceBindMqttClient

	// 使用 Where 和 Find 方法查询记录
	result := glob.GDb.Where("`device_info_id` = ?", param).Find(&deviceBindMqttClients)
	if result.Error != nil {
		zap.S().Infoln("Error occurred during query:", result.Error)
		servlet.Error(c, "暂无数据")
		return
	}
	servlet.Resp(c, deviceBindMqttClients)
}

// BindMqtt
// @Tags      DeviceInfos
// @Summary   绑定mqtt客户端
// @Accept json
// @Produce json
// @Param DeviceGroup body servlet.DeviceGroupCreateParam true "绑定参数"
// @Router    /DeviceInfo/BindMqtt [post]
func (api *DeviceInfoApi) BindMqtt(c *gin.Context) {
	var param servlet.DeviceBindMqttClientParam
	if err := c.ShouldBindJSON(&param); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	// 开启事务
	tx := glob.GDb.Begin()
	if tx.Error != nil {
		servlet.Error(c, "Failed to begin transaction")
		return
	}
	var toDel []models.DeviceBindMqttClient

	tx.Where("`device_info_id` = ?", param.DeviceId).Find(toDel)

	result := tx.Where("`device_info_id` = ?", param.DeviceId).Delete(&models.DeviceBindMqttClient{})

	if result.Error != nil {
		// 如果出现错误，回滚事务
		tx.Rollback()
		servlet.Error(c, "Error occurred during deletion")
		return
	}

	var groupBindMqttClients []models.DeviceBindMqttClient
	for _, mqttClientId := range param.MqttClientId {
		groupBindMqttClients = append(groupBindMqttClients, models.DeviceBindMqttClient{
			DeviceInfoId: uint(param.DeviceId),
			MqttClientId: uint(mqttClientId),
		})
	}

	result = tx.Model(&models.DeviceBindMqttClient{}).CreateInBatches(groupBindMqttClients, len(groupBindMqttClients))
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

	// redis 中建立 mqtt_client_id 与 device_info_id 的映射

	var DeviceInfo models.DeviceInfo

	first := tx.First(&DeviceInfo, param.DeviceId)
	if first.Error != nil {
		servlet.Error(c, "DeviceInfo not found")
		return
	}

	for _, client := range toDel {
		glob.GRedis.Del(context.Background(), "mqtt_client_id_bind_product:"+strconv.Itoa(int(client.MqttClientId)))
	}

	for _, item := range param.MqttClientId {
		glob.GRedis.LPush(context.Background(), "mqtt_client_id_bind_product:"+strconv.Itoa(item), DeviceInfo.ProductId)
	}


	for _, client := range toDel {
		glob.GRedis.Del(context.Background(), "mqtt_client_id_bind_device_info:"+strconv.Itoa(int(client.MqttClientId)))
	}

	for _, item := range param.MqttClientId {
		glob.GRedis.LPush(context.Background(), "mqtt_client_id_bind_device_info:"+strconv.Itoa(item), DeviceInfo.ID)

	}


	servlet.Resp(c, "绑定成功")

}
