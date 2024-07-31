package transmit_mqtt

import (
	"github.com/gin-gonic/gin"
	transmit "igp/biz/transmit/mqtt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type ClickhouseTransmitBindApi struct{}

var clickhouseTransmitBindBiz = transmit.ClickhouseTransmitBindBiz{}

// CreateClickhouseTransmitBind
// @Summary 创建ClickhouseTransmitBind
// @Description 创建ClickhouseTransmitBind
// @Tags ClickhouseTransmitBinds
// @Accept json
// @Produce json
// @Param ClickhouseTransmitBind body models.ClickhouseTransmitBind true "ClickhouseTransmitBind"
// @Success 201 {object} servlet.JSONResult{data=models.ClickhouseTransmitBind} "创建成功的ClickhouseTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ClickhouseTransmitBind/create [post]
func (api *ClickhouseTransmitBindApi) CreateClickhouseTransmitBind(c *gin.Context) {
	var ClickhouseTransmitBind models.ClickhouseTransmitBind
	if err := c.ShouldBindJSON(&ClickhouseTransmitBind); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 ClickhouseTransmitBind 是否被正确初始化
	if ClickhouseTransmitBind.Table == "" {
		servlet.Error(c, "表不能为空")
		return
	}

	result := glob.GDb.Create(&ClickhouseTransmitBind)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	clickhouseTransmitBindBiz.HandlerRedis(ClickhouseTransmitBind)
	// 返回创建成功的ClickhouseTransmitBind
	servlet.Resp(c, ClickhouseTransmitBind)
}

// UpdateClickhouseTransmitBind
// @Summary 更新一个ClickhouseTransmitBind
// @Description 更新一个ClickhouseTransmitBind
// @Tags ClickhouseTransmitBinds
// @Accept json
// @Produce json
// @Param ClickhouseTransmitBind body models.ClickhouseTransmitBind true "ClickhouseTransmitBind"
// @Success 200 {object}  servlet.JSONResult{data=models.ClickhouseTransmitBind} "ClickhouseTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "ClickhouseTransmitBind未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /ClickhouseTransmitBind/update [post]
func (api *ClickhouseTransmitBindApi) UpdateClickhouseTransmitBind(c *gin.Context) {
	var req models.ClickhouseTransmitBind
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.ClickhouseTransmitBind
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "ClickhouseTransmitBind not found")
		return
	}

	var newV models.ClickhouseTransmitBind
	newV = old
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	clickhouseTransmitBindBiz.HandlerRedis(newV)
	servlet.Resp(c, old)
}

// PageClickhouseTransmitBind
// @Summary 分页查询ClickhouseTransmitBind
// @Description 分页查询ClickhouseTransmitBind
// @Tags ClickhouseTransmitBinds
// @Accept json
// @Produce json
// @Param page table string false "表名" default(0)
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.ClickhouseTransmitBind}} "ClickhouseTransmitBind"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /ClickhouseTransmitBind/page [get]
func (api *ClickhouseTransmitBindApi) PageClickhouseTransmitBind(c *gin.Context) {
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

	data, err := clickhouseTransmitBindBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteClickhouseTransmitBind
// @Tags      ClickhouseTransmitBinds
// @Summary   删除ClickhouseTransmitBind
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /ClickhouseTransmitBind/delete/:id [post]
func (api *ClickhouseTransmitBindApi) DeleteClickhouseTransmitBind(c *gin.Context) {
	var ClickhouseTransmitBind models.ClickhouseTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&ClickhouseTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "ClickhouseTransmitBind not found")

		return
	}

	if result := glob.GDb.Delete(&ClickhouseTransmitBind); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	ClickhouseTransmitBind.Enable = false
	clickhouseTransmitBindBiz.HandlerRedis(ClickhouseTransmitBind)
	servlet.Resp(c, "删除成功")
}

// ByIdClickhouseTransmitBind
// @Tags      ClickhouseTransmitBinds
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /ClickhouseTransmitBind/:id [get]
func (api *ClickhouseTransmitBindApi) ByIdClickhouseTransmitBind(c *gin.Context) {
	var ClickhouseTransmitBind models.ClickhouseTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&ClickhouseTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "ClickhouseTransmitBind not found")

		return
	}

	servlet.Resp(c, ClickhouseTransmitBind)
}

// MockScript
// @Tags      ClickhouseTransmitBinds
// @Summary   模拟脚本
// @Param ClickhouseTransmit body servlet.TransmitScriptParam true "执行参数"
// @Produce   application/json
// @Router    /ClickhouseTransmitBind/mockScript [post]
func (api *ClickhouseTransmitBindApi) MockScript(c *gin.Context) {
	var req servlet.TransmitScriptParam
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}
	script := clickhouseTransmitBindBiz.MockScript(req.DataRowList, req.Script)
	servlet.Resp(c, script)
}
