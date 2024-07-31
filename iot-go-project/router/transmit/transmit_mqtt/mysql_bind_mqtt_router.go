package transmit_mqtt

import (
	"github.com/gin-gonic/gin"
	transmit "igp/biz/transmit/mqtt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type MySQLTransmitBindApi struct{}

var MySQLTransmitBindBiz = transmit.MySQLTransmitBindBiz{}

// CreateMySQLTransmitBind
// @Summary 创建MySQLTransmitBind
// @Description 创建MySQLTransmitBind
// @Tags MySQLTransmitBinds
// @Accept json
// @Produce json
// @Param MySQLTransmitBind body models.MySQLTransmitBind true "MySQLTransmitBind"
// @Success 201 {object} servlet.JSONResult{data=models.MySQLTransmitBind} "创建成功的MySQLTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /MySQLTransmitBind/create [post]
func (api *MySQLTransmitBindApi) CreateMySQLTransmitBind(c *gin.Context) {
	var MySQLTransmitBind models.MySQLTransmitBind
	if err := c.ShouldBindJSON(&MySQLTransmitBind); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 MySQLTransmitBind 是否被正确初始化
	if MySQLTransmitBind.Table == "" {
		servlet.Error(c, "表不能为空")
		return
	}

	result := glob.GDb.Create(&MySQLTransmitBind)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	MySQLTransmitBindBiz.HandlerRedis(MySQLTransmitBind)
	// 返回创建成功的MySQLTransmitBind
	servlet.Resp(c, MySQLTransmitBind)
}

// UpdateMySQLTransmitBind
// @Summary 更新一个MySQLTransmitBind
// @Description 更新一个MySQLTransmitBind
// @Tags MySQLTransmitBinds
// @Accept json
// @Produce json
// @Param MySQLTransmitBind body models.MySQLTransmitBind true "MySQLTransmitBind"
// @Success 200 {object}  servlet.JSONResult{data=models.MySQLTransmitBind} "MySQLTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "MySQLTransmitBind未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /MySQLTransmitBind/update [post]
func (api *MySQLTransmitBindApi) UpdateMySQLTransmitBind(c *gin.Context) {
	var req models.MySQLTransmitBind
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.MySQLTransmitBind
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "MySQLTransmitBind not found")
		return
	}

	var newV models.MySQLTransmitBind
	newV = old
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	MySQLTransmitBindBiz.HandlerRedis(newV)
	servlet.Resp(c, old)
}

// PageMySQLTransmitBind
// @Summary 分页查询MySQLTransmitBind
// @Description 分页查询MySQLTransmitBind
// @Tags MySQLTransmitBinds
// @Accept json
// @Produce json
// @Param page table string false "表名" default(0)
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.MySQLTransmitBind}} "MySQLTransmitBind"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /MySQLTransmitBind/page [get]
func (api *MySQLTransmitBindApi) PageMySQLTransmitBind(c *gin.Context) {
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

	data, err := MySQLTransmitBindBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteMySQLTransmitBind
// @Tags      MySQLTransmitBinds
// @Summary   删除MySQLTransmitBind
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /MySQLTransmitBind/delete/:id [post]
func (api *MySQLTransmitBindApi) DeleteMySQLTransmitBind(c *gin.Context) {
	var MySQLTransmitBind models.MySQLTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&MySQLTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "MySQLTransmitBind not found")

		return
	}

	if result := glob.GDb.Delete(&MySQLTransmitBind); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	MySQLTransmitBind.Enable = false
	MySQLTransmitBindBiz.HandlerRedis(MySQLTransmitBind)
	servlet.Resp(c, "删除成功")
}

// ByIdMySQLTransmitBind
// @Tags      MySQLTransmitBinds
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /MySQLTransmitBind/:id [get]
func (api *MySQLTransmitBindApi) ByIdMySQLTransmitBind(c *gin.Context) {
	var MySQLTransmitBind models.MySQLTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&MySQLTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "MySQLTransmitBind not found")

		return
	}

	servlet.Resp(c, MySQLTransmitBind)
}




// MockScript
// @Tags      MySQLTransmits
// @Summary   模拟脚本
// @Param MySQLTransmit body servlet.TransmitScriptParam true "执行参数"
// @Produce   application/json
// @Router    /MySQLTransmitBind/mockScript [post]
func (api *MySQLTransmitBindApi) MockScript(c *gin.Context) {
	var req servlet.TransmitScriptParam
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}
	script := MySQLTransmitBindBiz.MockScript(req.DataRowList, req.Script)
	servlet.Resp(c, script)
}

