package transmit

import (
	"github.com/gin-gonic/gin"
	"igp/biz/transmit"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type MySQLTransmitApi struct{}

var mySQLTransmit = transmit.MySQLTransmitBiz{}

// CreateMySQLTransmit
// @Summary 创建MySql数据库管理
// @Description 创建MySql数据库管理
// @Tags MySQLTransmits
// @Accept json
// @Produce json
// @Param MySQLTransmit body models.MySQLTransmit true "MySql数据库管理"
// @Success 201 {object} servlet.JSONResult{data=models.MySQLTransmit} "创建成功的MySql数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /MySQLTransmit/create [post]
func (api *MySQLTransmitApi) CreateMySQLTransmit(c *gin.Context) {
	var MySQLTransmit models.MySQLTransmit
	if err := c.ShouldBindJSON(&MySQLTransmit); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 MySQLTransmit 是否被正确初始化
	if MySQLTransmit.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&MySQLTransmit)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的MySql数据库管理
	servlet.Resp(c, MySQLTransmit)
}

// UpdateMySQLTransmit
// @Summary 更新一个MySql数据库管理
// @Description 更新一个MySql数据库管理
// @Tags MySQLTransmits
// @Accept json
// @Produce json
// @Param MySQLTransmit body models.MySQLTransmit true "MySql数据库管理"
// @Success 200 {object}  servlet.JSONResult{data=models.MySQLTransmit} "MySql数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "MySql数据库管理未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /MySQLTransmit/update [post]
func (api *MySQLTransmitApi) UpdateMySQLTransmit(c *gin.Context) {
	var req models.MySQLTransmit
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.MySQLTransmit
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "MySQLTransmit not found")
		return
	}

	var newV models.MySQLTransmit
	newV = old
	newV.Name  = req.Name
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageMySQLTransmit
// @Summary 分页查询MySql数据库管理
// @Description 分页查询MySql数据库管理
// @Tags MySQLTransmits
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.MySQLTransmit}} "MySql数据库管理"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /MySQLTransmit/page [get]
func (api *MySQLTransmitApi) PageMySQLTransmit(c *gin.Context) {
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

	data, err := mySQLTransmit.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteMySQLTransmit
// @Tags      MySQLTransmits
// @Summary   删除MySql数据库管理
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /MySQLTransmit/delete/:id [post]
func (api *MySQLTransmitApi) DeleteMySQLTransmit(c *gin.Context) {
	var MySQLTransmit models.MySQLTransmit

	param := c.Param("id")

	result := glob.GDb.First(&MySQLTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "MySQLTransmit not found")

		return
	}

	if result := glob.GDb.Delete(&MySQLTransmit); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdMySQLTransmit
// @Tags      MySQLTransmits
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /MySQLTransmit/:id [get]
func (api *MySQLTransmitApi) ByIdMySQLTransmit(c *gin.Context) {
	var MySQLTransmit models.MySQLTransmit

	param := c.Param("id")

	result := glob.GDb.First(&MySQLTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "MySQLTransmit not found")

		return
	}

	servlet.Resp(c, MySQLTransmit)
}

// MockScript
// @Tags      MySQLTransmits
// @Summary   模拟脚本
// @Param MySQLTransmit body servlet.TransmitScriptParam true "执行参数"
// @Produce   application/json
// @Router    /MySQLTransmit/mockScript [post]
func (api *MySQLTransmitApi) MockScript(c *gin.Context) {
	var req servlet.TransmitScriptParam
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}
	script := mySQLTransmit.MockScript(req.DataRowList, req.Script)
	servlet.Resp(c, script)
}
