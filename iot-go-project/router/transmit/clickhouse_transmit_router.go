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

package transmit

import (
	"github.com/gin-gonic/gin"
	"igp/biz/transmit"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type ClickhouseTransmitApi struct{}

var ClickhouseTransmitBiz = transmit.ClickhouseTransmitBiz{}

// CreateClickhouseTransmit
// @Summary 创建Clickhouse数据库管理
// @Description 创建Clickhouse数据库管理
// @Tags ClickhouseTransmits
// @Accept json
// @Produce json
// @Param ClickhouseTransmit body models.ClickhouseTransmit true "Clickhouse数据库管理"
// @Success 200 {object} servlet.JSONResult{data=models.ClickhouseTransmit} "创建成功的Clickhouse数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ClickhouseTransmit/create [post]
func (api *ClickhouseTransmitApi) CreateClickhouseTransmit(c *gin.Context) {
	var ClickhouseTransmit models.ClickhouseTransmit
	if err := c.ShouldBindJSON(&ClickhouseTransmit); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := glob.GDb.Create(&ClickhouseTransmit)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	ClickhouseTransmitBiz.SetRedis(ClickhouseTransmit)
	// 返回创建成功的Clickhouse数据库管理
	servlet.Resp(c, ClickhouseTransmit)
}

// UpdateClickhouseTransmit
// @Summary 更新一个Clickhouse数据库管理
// @Description 更新一个Clickhouse数据库管理
// @Tags ClickhouseTransmits
// @Accept json
// @Produce json
// @Param ClickhouseTransmit body models.ClickhouseTransmit true "Clickhouse数据库管理"
// @Success 200 {object}  servlet.JSONResult{data=models.ClickhouseTransmit} "Clickhouse数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Clickhouse数据库管理未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ClickhouseTransmit/update [post]
func (api *ClickhouseTransmitApi) UpdateClickhouseTransmit(c *gin.Context) {
	var req models.ClickhouseTransmit
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.ClickhouseTransmit
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "ClickhouseTransmit not found")
		return
	}

	var newV models.ClickhouseTransmit
	newV = old

	newV.Name = req.Name
	newV.Host = req.Host
	newV.Port = req.Port
	newV.Username = req.Username
	newV.Password = req.Password
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	ClickhouseTransmitBiz.SetRedis(newV)
	servlet.Resp(c, old)
}

// PageClickhouseTransmit
// @Summary 分页查询Clickhouse数据库管理
// @Description 分页查询Clickhouse数据库管理
// @Tags ClickhouseTransmits
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.ClickhouseTransmit}} "Clickhouse数据库管理"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /ClickhouseTransmit/page [get]
func (api *ClickhouseTransmitApi) PageClickhouseTransmit(c *gin.Context) {
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

	data, err := ClickhouseTransmitBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteClickhouseTransmit
// @Tags      ClickhouseTransmits
// @Summary   删除Clickhouse数据库管理
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /ClickhouseTransmit/delete/:id [post]
// @Success 200 {object}  servlet.JSONResult{data=string} 
func (api *ClickhouseTransmitApi) DeleteClickhouseTransmit(c *gin.Context) {
	var ClickhouseTransmit models.ClickhouseTransmit

	param := c.Param("id")

	result := glob.GDb.First(&ClickhouseTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "ClickhouseTransmit not found")

		return
	}

	if result := glob.GDb.Delete(&ClickhouseTransmit); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	ClickhouseTransmitBiz.DeleteRedis(ClickhouseTransmit)
	servlet.Resp(c, "删除成功")
}

// ByIdClickhouseTransmit
// @Tags      ClickhouseTransmits
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /ClickhouseTransmit/:id [get]
// @Success 200 {object}  servlet.JSONResult{data=models.ClickhouseTransmit} 
func (api *ClickhouseTransmitApi) ByIdClickhouseTransmit(c *gin.Context) {
	var ClickhouseTransmit models.ClickhouseTransmit

	param := c.Param("id")

	result := glob.GDb.First(&ClickhouseTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "ClickhouseTransmit not found")

		return
	}

	servlet.Resp(c, ClickhouseTransmit)
}

// ListClickhouseTransmit
// @Tags      ClickhouseTransmits
// @Summary   单个详情
// @Produce   application/json
// @Router    /ClickhouseTransmit/list [get]
// @Success 200 {object}  servlet.JSONResult{data=models.ClickhouseTransmit[]}
func (api *ClickhouseTransmitApi) ListClickhouseTransmit(c *gin.Context) {
	var ClickhouseTransmit []models.ClickhouseTransmit

 glob.GDb.Find(&ClickhouseTransmit)


	servlet.Resp(c, ClickhouseTransmit)
}
