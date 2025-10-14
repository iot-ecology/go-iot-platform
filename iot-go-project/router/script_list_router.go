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
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type ScriptListApi struct{}

var ScriptListBiz = biz.ScriptListBiz{}

// CreateScriptList
// @Summary 创建脚本
// @Description 创建脚本
// @Tags ScriptLists
// @Accept json
// @Produce json
// @Param ScriptListId body models.ScriptList true "脚本"
// @Success 201 {object} servlet.JSONResult{data=models.ScriptList} "创建成功的脚本"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ScriptListId/create [post]
func (api *ScriptListApi) CreateScriptList(c *gin.Context) {
	var ScriptList models.ScriptList
	if err := c.ShouldBindJSON(&ScriptList); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 ScriptListId 是否被正确初始化
	if ScriptList.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&ScriptList)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的脚本
	servlet.Resp(c, ScriptList)
}

// UpdateScriptList
// @Summary 更新一个脚本
// @Description 更新一个脚本
// @Tags ScriptLists
// @Accept json
// @Produce json
// @Param ScriptListId body models.ScriptList true "脚本"
// @Success 200 {object}  servlet.JSONResult{data=models.ScriptList} "脚本"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "脚本未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ScriptListId/update [post]
func (api *ScriptListApi) UpdateScriptList(c *gin.Context) {
	var req models.ScriptList
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.ScriptList
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "ScriptListId not found")
		return
	}

	var newV models.ScriptList
	newV = old
	newV.Name = req.Name
	newV.Content = req.Content

	m := structs.Map(newV)
	result = glob.GDb.Model(models.ScriptList{}).Where("id = ?", newV.ID).Updates(m)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, old)
}

// PageScriptList
// @Summary 分页查询脚本
// @Description 分页查询脚本
// @Tags ScriptLists
// @Accept json
// @Produce json
// @Param name query string false "名称"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.ScriptList}} "脚本"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /ScriptListId/page [get]
func (api *ScriptListApi) PageScriptList(c *gin.Context) {
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

	data, err := ScriptListBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteScriptList
// @Tags      ScriptLists
// @Summary   删除脚本
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /ScriptListId/delete/:id [post]
func (api *ScriptListApi) DeleteScriptList(c *gin.Context) {
	var ScriptList models.ScriptList

	param := c.Param("id")

	result := glob.GDb.First(&ScriptList, param)
	if result.Error != nil {
		servlet.Error(c, "ScriptListId not found")

		return
	}

	if result := glob.GDb.Delete(&ScriptList); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdScriptList
// @Tags      ScriptLists
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /ScriptListId/:id [get]
func (api *ScriptListApi) ByIdScriptList(c *gin.Context) {
	var ScriptList models.ScriptList

	param := c.Param("id")

	result := glob.GDb.First(&ScriptList, param)
	if result.Error != nil {
		servlet.Error(c, "ScriptListId not found")

		return
	}

	servlet.Resp(c, ScriptList)
}
