package router

import (
	"context"
	"encoding/json"
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
		DeviceInfo.WarrantyExpiry = WarrantyExpiry
	}

	if DeviceInfo.Source == 2 {
		DeviceInfo.ManufacturingDate = DeviceInfo.ProcurementDate
		WarrantyExpiry := DeviceInfo.ManufacturingDate.AddDate(0, 0, Product.WarrantyPeriod)
		DeviceInfo.WarrantyExpiry = WarrantyExpiry
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
	newV.ProductId = req.ProductId
	newV.Source = req.Source
	newV.SN = req.SN
	newV.ManufacturingDate = req.ManufacturingDate
	newV.ProcurementDate = req.ProcurementDate
	newV.WarrantyExpiry = req.WarrantyExpiry
	newV.PushInterval = req.PushInterval
	newV.ErrorRate = req.ErrorRate


	var Product models.Product
	result = glob.GDb.First(&Product, newV.ProductId)
	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	if newV.Source == 2 {
		newV.ManufacturingDate = req.ProcurementDate
		WarrantyExpiry := newV.ManufacturingDate.AddDate(0, 0, Product.WarrantyPeriod)
		newV.WarrantyExpiry = WarrantyExpiry
	}

	if !newV.ManufacturingDate.IsZero() {
		WarrantyExpiry := newV.ManufacturingDate.AddDate(0, 0, Product.WarrantyPeriod)
		newV.WarrantyExpiry = WarrantyExpiry
	}
	result = glob.GDb.Model(&newV).Updates(newV)


	if result.Error != nil {
		zap.S().Errorw("更新 DeviceInfo 失败", "error", result.Error)
		servlet.Error(c, result.Error.Error())
		return
	}
	deviceInfoBiz.SetRedis(newV)

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

// QueryBindHttp
// @Tags      DeviceInfos
// @Summary   查询绑定HTTP客户端
// @Accept json
// @Produce json
// @Param device_info_id path int true "主键"
// @Router    /DeviceInfo/QueryBindHTTP [get]
func (api *DeviceInfoApi) QueryBindHttp(c *gin.Context) {
	param := c.Param("device_info_id")

	var res []models.DeviceBindTcpHandler

	// 使用 Where 和 Find 方法查询记录
	result := glob.GDb.Where("`device_info_id` = ?", param).Find(&res)
	if result.Error != nil {
		zap.S().Infoln("Error occurred during query:", result.Error)
		servlet.Error(c, "暂无数据")
		return
	}
	servlet.Resp(c, res)
}
// QueryBindCoap
// @Tags      DeviceInfos
// @Summary   查询绑定HTTP客户端
// @Accept json
// @Produce json
// @Param device_info_id path int true "主键"
// @Router    /DeviceInfo/QueryBindCoap [get]
func (api *DeviceInfoApi) QueryBindCoap(c *gin.Context) {
	param := c.Param("device_info_id")

	var res []models.CoapHandler

	// 使用 Where 和 Find 方法查询记录
	result := glob.GDb.Where("`device_info_id` = ?", param).Find(&res)
	if result.Error != nil {
		zap.S().Infoln("Error occurred during query:", result.Error)
		servlet.Error(c, "暂无数据")
		return
	}
	servlet.Resp(c, res)
}

// QueryBindTcp
// @Tags      DeviceInfos
// @Summary   查询绑定tcp客户端
// @Accept json
// @Produce json
// @Param device_info_id path int true "主键"
// @Router    /DeviceInfo/QueryBindTcp [get]
func (api *DeviceInfoApi) QueryBindTcp(c *gin.Context) {
	param := c.Param("device_info_id")

	var res []models.DeviceBindTcpHandler

	// 使用 Where 和 Find 方法查询记录
	result := glob.GDb.Where("`device_info_id` = ?", param).Find(&res)
	if result.Error != nil {
		zap.S().Infoln("Error occurred during query:", result.Error)
		servlet.Error(c, "暂无数据")
		return
	}
	servlet.Resp(c, res)
}

// BindMqtt
// @Tags      DeviceInfos
// @Summary   绑定mqtt客户端
// @Accept json
// @Produce json
// @Param DeviceGroup body servlet.DeviceBindMqttClientParam true "绑定参数"
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

// BindTcp
// @Tags      DeviceInfos
// @Summary   绑定tcp处理器
// @Accept json
// @Produce json
// @Param DeviceGroup body servlet.DeviceBindTcpParam true "绑定参数"
// @Router    /DeviceInfo/BindTcp [post]
func (api *DeviceInfoApi) BindTcp(c *gin.Context) {
	var param servlet.DeviceBindTcpParam
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
	var toDel []models.DeviceBindTcpHandler

	tx.Where("`device_info_id` = ?", param.DeviceId).Find(toDel)

	result := tx.Where("`device_info_id` = ?", param.DeviceId).Delete(&models.DeviceBindTcpHandler{})

	if result.Error != nil {
		// 如果出现错误，回滚事务
		tx.Rollback()
		servlet.Error(c, "Error occurred during deletion")
		return
	}

	var deviceBindTcpHandlers []models.DeviceBindTcpHandler
	for _, tcpId := range param.TcpHandlerId {
		deviceBindTcpHandlers = append(deviceBindTcpHandlers, models.DeviceBindTcpHandler{
			DeviceInfoId: uint(param.DeviceId),
			TcpHandlerId: uint(tcpId),
		})
	}

	result = tx.Model(&models.DeviceBindTcpHandler{}).CreateInBatches(deviceBindTcpHandlers, len(deviceBindTcpHandlers))
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

	// redis 中建立 tcp 与 device_info_id 的映射

	var DeviceInfo models.DeviceInfo

	first := tx.First(&DeviceInfo, param.DeviceId)
	if first.Error != nil {
		servlet.Error(c, "DeviceInfo not found")
		return
	}

	for _, client := range toDel {
		glob.GRedis.Del(context.Background(), "tcp_bind_product:"+strconv.Itoa(int(client.TcpHandlerId)))
	}

	for _, item := range param.TcpHandlerId {
		glob.GRedis.LPush(context.Background(), "tcp_bind_product:"+strconv.Itoa(item), DeviceInfo.ProductId)
	}

	for _, client := range toDel {
		glob.GRedis.Del(context.Background(), "tcp_bind_device_info:"+strconv.Itoa(int(client.TcpHandlerId)))
	}

	for _, item := range param.TcpHandlerId {
		glob.GRedis.LPush(context.Background(), "tcp_bind_device_info:"+strconv.Itoa(item), DeviceInfo.ID)

	}

	servlet.Resp(c, "绑定成功")

}

// BindHTTP
// @Tags      DeviceInfos
// @Summary   绑定http处理器
// @Accept json
// @Produce json
// @Param DeviceGroup body models.HttpHandler true "绑定参数"
// @Router    /DeviceInfo/BindHTTP [post]
func (api *DeviceInfoApi) BindHTTP(c *gin.Context) {
	var param models.HttpHandler
	if err := c.ShouldBindJSON(&param); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	if param.ID != 0 {
		var old models.HttpHandler

		result := glob.GDb.First(&old, param.ID)
		if result.Error != nil {

			servlet.Error(c, "HttpHandler not found")
			return
		}
		var newV models.HttpHandler
		newV = old
		newV.DeviceInfoId = param.DeviceInfoId
		newV.Name = param.Name
		newV.Username = param.Username
		newV.Password = param.Password
		newV.Script = param.Script
		// 更新记录
		result = glob.GDb.Model(&newV).Updates(newV)
		setHttpHandlerRedis(newV)
	} else {
		// 新增
		glob.GDb.Model(models.HttpHandler{}).Create(&param)
		setHttpHandlerRedis(param)

	}

	glob.GRedis.LPush(context.Background(), "http_bind_device_info:"+strconv.Itoa(int(param.ID)), param.DeviceInfoId)

	servlet.Resp(c, "绑定成功")

}

// BindHCoap
// @Tags      DeviceInfos
// @Summary   绑定coap处理器
// @Accept json
// @Produce json
// @Param DeviceGroup body models.CoapHandler true "绑定参数"
// @Router    /DeviceInfo/BindHCoap [post]
func (api *DeviceInfoApi) BindHCoap(c *gin.Context) {
	var param models.CoapHandler
	if err := c.ShouldBindJSON(&param); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	if param.ID != 0 {
		var old models.CoapHandler

		result := glob.GDb.First(&old, param.ID)
		if result.Error != nil {

			servlet.Error(c, "HttpHandler not found")
			return
		}
		var newV models.CoapHandler
		newV = old
		newV.DeviceInfoId = param.DeviceInfoId
		newV.Name = param.Name
		newV.Username = param.Username
		newV.Password = param.Password
		newV.Script = param.Script
		// 更新记录
		result = glob.GDb.Model(&newV).Updates(newV)
		setCoapHandlerRedis(newV)
	} else {
		// 新增
		glob.GDb.Model(models.CoapHandler{}).Create(&param)
		setCoapHandlerRedis(param)

	}

	glob.GRedis.LPush(context.Background(), "coap_bind_device_info:"+strconv.Itoa(int(param.ID)), param.DeviceInfoId)

	servlet.Resp(c, "绑定成功")

}

func setHttpHandlerRedis(config models.HttpHandler) {
	jsonData, _ := json.Marshal(config)
	glob.GRedis.HSet(context.Background(), "auth:http", strconv.Itoa(int(config.DeviceInfoId)), jsonData)
}
func setCoapHandlerRedis(config models.CoapHandler) {
	jsonData, _ := json.Marshal(config)
	glob.GRedis.HSet(context.Background(), "auth:coap", strconv.Itoa(int(config.DeviceInfoId)), jsonData)
}
