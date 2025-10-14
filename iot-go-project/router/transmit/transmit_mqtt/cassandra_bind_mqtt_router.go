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

package transmit_mqtt

import (
	"github.com/gin-gonic/gin"
	transmit "igp/biz/transmit/mqtt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type CassandraTransmitBindApi struct{}

var dashBiz = transmit.CassandraTransmitBindBiz{}

// CreateCassandraTransmitBind
// @Summary 创建CassandraTransmitBind
// @Description 创建CassandraTransmitBind
// @Tags CassandraTransmitBinds
// @Accept json
// @Produce json
// @Param CassandraTransmitBind body models.CassandraTransmitBind true "CassandraTransmitBind"
// @Success 200 {object} servlet.JSONResult{data=models.CassandraTransmitBind} "创建成功的CassandraTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /CassandraTransmitBind/create [post]
func (api *CassandraTransmitBindApi) CreateCassandraTransmitBind(c *gin.Context) {
	var CassandraTransmitBind models.CassandraTransmitBind
	if err := c.ShouldBindJSON(&CassandraTransmitBind); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 CassandraTransmitBind 是否被正确初始化
	if CassandraTransmitBind.Table == "" {
		servlet.Error(c, "表不能为空")
		return
	}

	result := glob.GDb.Create(&CassandraTransmitBind)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	dashBiz.HandlerRedis(CassandraTransmitBind)
	// 返回创建成功的CassandraTransmitBind
	servlet.Resp(c, CassandraTransmitBind)
}

// UpdateCassandraTransmitBind
// @Summary 更新一个CassandraTransmitBind
// @Description 更新一个CassandraTransmitBind
// @Tags CassandraTransmitBinds
// @Accept json
// @Produce json
// @Param CassandraTransmitBind body models.CassandraTransmitBind true "CassandraTransmitBind"
// @Success 200 {object}  servlet.JSONResult{data=models.CassandraTransmitBind} "CassandraTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "CassandraTransmitBind未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /CassandraTransmitBind/update [post]
func (api *CassandraTransmitBindApi) UpdateCassandraTransmitBind(c *gin.Context) {
	var req models.CassandraTransmitBind
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.CassandraTransmitBind
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "CassandraTransmitBind not found")
		return
	}

	var newV models.CassandraTransmitBind
	newV = old
	newV.DeviceUid = req.DeviceUid
	newV.Protocol = req.Protocol
	newV.IdentificationCode = req.IdentificationCode
	newV.CassandraTransmitId = req.CassandraTransmitId
	newV.Database = req.Database
	newV.Table = req.Table
	newV.Script = req.Script
	newV.Enable = req.Enable
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	dashBiz.HandlerRedis(newV)

	servlet.Resp(c, old)
}

// PageCassandraTransmitBind
// @Summary 分页查询CassandraTransmitBind
// @Description 分页查询CassandraTransmitBind
// @Tags CassandraTransmitBinds
// @Accept json
// @Produce json
// @Param table query string false "表名"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.CassandraTransmitBind}} "CassandraTransmitBind"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /CassandraTransmitBind/page [get]
func (api *CassandraTransmitBindApi) PageCassandraTransmitBind(c *gin.Context) {
	var name = c.Query("table")
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

	data, err := dashBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteCassandraTransmitBind
// @Tags      CassandraTransmitBinds
// @Summary   删除CassandraTransmitBind
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /CassandraTransmitBind/delete/:id [post]
// @Success 200 {object}  servlet.JSONResult{data=string} 
func (api *CassandraTransmitBindApi) DeleteCassandraTransmitBind(c *gin.Context) {
	var CassandraTransmitBind models.CassandraTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&CassandraTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "CassandraTransmitBind not found")

		return
	}

	if result := glob.GDb.Delete(&CassandraTransmitBind); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	CassandraTransmitBind.Enable = false
	dashBiz.HandlerRedis(CassandraTransmitBind)
	servlet.Resp(c, "删除成功")
}

// ByIdCassandraTransmitBind
// @Tags      CassandraTransmitBinds
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /CassandraTransmitBind/:id [get]
// @Success 200 {object}  servlet.JSONResult{data=models.CassandraTransmitBind} 
func (api *CassandraTransmitBindApi) ByIdCassandraTransmitBind(c *gin.Context) {
	var CassandraTransmitBind models.CassandraTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&CassandraTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "CassandraTransmitBind not found")

		return
	}

	servlet.Resp(c, CassandraTransmitBind)
}



// MockScript
// @Tags      CassandraTransmits
// @Summary   模拟脚本
// @Param CassandraTransmit body servlet.TransmitScriptParam true "执行参数"
// @Produce   application/json
// @Router    /CassandraTransmitBind/mockScript [post]
// @Success 200 {object}  servlet.JSONResult{data=string} 
func (api *CassandraTransmitBindApi) MockScript(c *gin.Context) {
	var req servlet.TransmitScriptParam
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}
	script := dashBiz.MockScript(req.DataRowList, req.Script)
	servlet.Resp(c, script)
}
