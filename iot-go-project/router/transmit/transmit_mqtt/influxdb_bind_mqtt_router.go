package transmit_mqtt

import (
	"github.com/gin-gonic/gin"
	transmit "igp/biz/transmit/mqtt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type InfluxdbTransmitBindApi struct{}

var InfluxdbTransmitBindBiz = transmit.InfluxdbTransmitBindBiz{}

// CreateInfluxdbTransmitBind
// @Summary 创建InfluxdbTransmitBind
// @Description 创建InfluxdbTransmitBind
// @Tags InfluxdbTransmitBinds
// @Accept json
// @Produce json
// @Param InfluxdbTransmitBind body models.InfluxdbTransmitBind true "InfluxdbTransmitBind"
// @Success 200 {object} servlet.JSONResult{data=models.InfluxdbTransmitBind} "创建成功的InfluxdbTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /InfluxdbTransmitBind/create [post]
func (api *InfluxdbTransmitBindApi) CreateInfluxdbTransmitBind(c *gin.Context) {
	var InfluxdbTransmitBind models.InfluxdbTransmitBind
	if err := c.ShouldBindJSON(&InfluxdbTransmitBind); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 InfluxdbTransmitBind 是否被正确初始化
	if InfluxdbTransmitBind.Org == "" {
		servlet.Error(c, "表不能为空")
		return
	}

	result := glob.GDb.Create(&InfluxdbTransmitBind)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	InfluxdbTransmitBindBiz.HandlerRedis(InfluxdbTransmitBind)
	// 返回创建成功的InfluxdbTransmitBind
	servlet.Resp(c, InfluxdbTransmitBind)
}

// UpdateInfluxdbTransmitBind
// @Summary 更新一个InfluxdbTransmitBind
// @Description 更新一个InfluxdbTransmitBind
// @Tags InfluxdbTransmitBinds
// @Accept json
// @Produce json
// @Param InfluxdbTransmitBind body models.InfluxdbTransmitBind true "InfluxdbTransmitBind"
// @Success 200 {object}  servlet.JSONResult{data=models.InfluxdbTransmitBind} "InfluxdbTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "InfluxdbTransmitBind未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /InfluxdbTransmitBind/update [post]
func (api *InfluxdbTransmitBindApi) UpdateInfluxdbTransmitBind(c *gin.Context) {
	var req models.InfluxdbTransmitBind
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.InfluxdbTransmitBind
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "InfluxdbTransmitBind not found")
		return
	}

	var newV models.InfluxdbTransmitBind
	newV = old
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	InfluxdbTransmitBindBiz.HandlerRedis(newV)
	servlet.Resp(c, old)
}

// PageInfluxdbTransmitBind
// @Summary 分页查询InfluxdbTransmitBind
// @Description 分页查询InfluxdbTransmitBind
// @Tags InfluxdbTransmitBinds
// @Accept json
// @Produce json
// @Param org query string false "组织" default(0)
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.InfluxdbTransmitBind}} "InfluxdbTransmitBind"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /InfluxdbTransmitBind/page [get]
func (api *InfluxdbTransmitBindApi) PageInfluxdbTransmitBind(c *gin.Context) {
	var name = c.Query("org")
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

	data, err := InfluxdbTransmitBindBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteInfluxdbTransmitBind
// @Tags      InfluxdbTransmitBinds
// @Summary   删除InfluxdbTransmitBind
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /InfluxdbTransmitBind/delete/:id [post]
// @Success 200 {object}  servlet.JSONResult{data=string} 
func (api *InfluxdbTransmitBindApi) DeleteInfluxdbTransmitBind(c *gin.Context) {
	var InfluxdbTransmitBind models.InfluxdbTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&InfluxdbTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "InfluxdbTransmitBind not found")

		return
	}

	if result := glob.GDb.Delete(&InfluxdbTransmitBind); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	InfluxdbTransmitBind.Enable = false
	InfluxdbTransmitBindBiz.HandlerRedis(InfluxdbTransmitBind)
	servlet.Resp(c, "删除成功")
}

// ByIdInfluxdbTransmitBind
// @Tags      InfluxdbTransmitBinds
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /InfluxdbTransmitBind/:id [get]
// @Success 200 {object}  servlet.JSONResult{data=models.InfluxdbTransmitBind} 
func (api *InfluxdbTransmitBindApi) ByIdInfluxdbTransmitBind(c *gin.Context) {
	var InfluxdbTransmitBind models.InfluxdbTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&InfluxdbTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "InfluxdbTransmitBind not found")

		return
	}

	servlet.Resp(c, InfluxdbTransmitBind)
}


// MockScript
// @Tags      InfluxdbTransmits
// @Summary   模拟脚本
// @Param InfluxdbTransmit body servlet.TransmitScriptParam true "执行参数"
// @Produce   application/json
// @Router    /InfluxdbTransmitBind/mockScript [post]
// @Success 200 {object}  servlet.JSONResult{data=string} 
func (api *InfluxdbTransmitBindApi) MockScript(c *gin.Context) {
	var req servlet.TransmitScriptParam
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}
	script := InfluxdbTransmitBindBiz.MockScript(req.DataRowList, req.Script)
	servlet.Resp(c, script)
}
