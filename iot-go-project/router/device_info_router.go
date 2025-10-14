/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package router

import (
	"context"
	"encoding/json"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
// @Success 200 {object} servlet.JSONResult{data=models.DeviceInfo} "创建成功的设备详情"
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

	sn := deviceInfoBiz.FindBySn(DeviceInfo.SN)
	if sn.ID > 0 {
		servlet.Error(c, "设备已存在")
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


	result = glob.GDb.Model(models.DeviceInfo{}).Create(&DeviceInfo)

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
	newV.ManufacturingDate = req.ManufacturingDate
	newV.ProcurementDate = req.ProcurementDate
	newV.WarrantyExpiry = req.WarrantyExpiry
	newV.PushInterval = req.PushInterval
	newV.ErrorRate = req.ErrorRate
	newV.Protocol = req.Protocol

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
	var protocol = c.Query("protocol")
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

	data, err := deviceInfoBiz.PageData(sn, protocol, parseUint, u)
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
// @Success 200 {object}  servlet.JSONResult{data=string}
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
// @Success 200 {object}  servlet.JSONResult{data=models.DeviceInfo}
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




// ListDeviceInfo
// @Tags      DeviceInfos
// @Summary   设备列表
// @Accept json
// @Produce json
// @Router    /DeviceInfo/list [get]
// @Success 200 {object}  servlet.JSONResult{data=models.DeviceInfo[]}
func (api *DeviceInfoApi) ListDeviceInfo(c *gin.Context) {
	var res []models.DeviceInfo
	glob.GDb.Find(&res)
	servlet.Resp(c, res)

}

func SetHttpHandlerRedis(config models.HttpHandler) {
	jsonData, _ := json.Marshal(config)
	glob.GRedis.HSet(context.Background(), "auth:http", strconv.Itoa(int(config.DeviceInfoId)), jsonData)
}
func SetCoapHandlerRedis(config models.CoapHandler) {
	jsonData, _ := json.Marshal(config)
	glob.GRedis.HSet(context.Background(), "auth:coap", strconv.Itoa(int(config.DeviceInfoId)), jsonData)
}
func SetWebsocketHandlerRedis(config models.WebsocketHandler) {
	jsonData, _ := json.Marshal(config)
	glob.GRedis.HSet(context.Background(), "auth:ws", strconv.Itoa(int(config.DeviceInfoId)), jsonData)
}
func SetTcpIpHandlerRedis(config models.TcpHandler) {
	jsonData, _ := json.Marshal(config)
	glob.GRedis.HSet(context.Background(), "auth:tcp", strconv.Itoa(int(config.DeviceInfoId)), jsonData)
}


// BindHandlers
// @Summary 绑定处理器
// @Description 绑定处理器
// @Tags DeviceInfos
// @Accept json
// @Produce json
// @Param DeviceInfo body models.DeviceInfo true "设备详情"
// @Success 200 {object} servlet.JSONResult{data=models.DeviceInfo} "创建成功的设备详情"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /DeviceInfo/BindHandlers [post]
func (api *DeviceInfoApi) BindHandlers(c *gin.Context) {
	type param struct {
		DeviceInfoId       uint      `json:"device_info_id" structs:"device_info_id"`   // 设备表的外键ID
		Protocol   string `json:"protocol"` // 协议
		IdentificationCode string `json:"identification_code"` // 设备标识码
		HandlerId uint `json:"handler_id"` // mqtt client 表主键 、TcpHandler Ws .... id
	}

	var p param
	if err := c.ShouldBindJSON(&p); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// Validate required fields
	if p.DeviceInfoId == 0 || p.Protocol == "" || p.IdentificationCode == "" || p.HandlerId == 0 {
		servlet.Error(c, "missing required fields")
		return
	}

	// Start transaction
	tx := glob.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing bindings
	if err := tx.Where("device_info_id = ?", p.DeviceInfoId).Delete(&models.DeviceBindHandler{}).Error; err != nil {
		tx.Rollback()
		servlet.Error(c, "failed to delete existing bindings")
		return
	}

	// Create new binding
	newBinding := models.DeviceBindHandler{
		DeviceInfoId: p.DeviceInfoId,
		Protocol: p.Protocol,
		IdentificationCode: p.IdentificationCode,
		HandlerId: p.HandlerId,
	}

	if err := tx.Create(&newBinding).Error; err != nil {
		tx.Rollback()
		servlet.Error(c, "failed to create new binding")
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		servlet.Error(c, "failed to commit transaction")
		return
	}

	servlet.Resp(c, newBinding)
}