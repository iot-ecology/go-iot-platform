package transmit_mqtt

import (
	"github.com/gin-gonic/gin"
	transmit "igp/biz/transmit/mqtt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type MongoTransmitBindApi struct{}

var MongoTransmitBindBiz = transmit.MongoTransmitBindBiz{}

// CreateMongoTransmitBind
// @Summary 创建MongoTransmitBind
// @Description 创建MongoTransmitBind
// @Tags MongoTransmitBinds
// @Accept json
// @Produce json
// @Param MongoTransmitBind body models.MongoTransmitBind true "MongoTransmitBind"
// @Success 200 {object} servlet.JSONResult{data=models.MongoTransmitBind} "创建成功的MongoTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /MongoTransmitBind/create [post]
func (api *MongoTransmitBindApi) CreateMongoTransmitBind(c *gin.Context) {
	var MongoTransmitBind models.MongoTransmitBind
	if err := c.ShouldBindJSON(&MongoTransmitBind); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 MongoTransmitBind 是否被正确初始化
	if MongoTransmitBind.Collection == "" {
		servlet.Error(c, "Collection不能为空")
		return
	}

	result := glob.GDb.Create(&MongoTransmitBind)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	MongoTransmitBindBiz.HandlerRedis(MongoTransmitBind)
	// 返回创建成功的MongoTransmitBind
	servlet.Resp(c, MongoTransmitBind)
}

// UpdateMongoTransmitBind
// @Summary 更新一个MongoTransmitBind
// @Description 更新一个MongoTransmitBind
// @Tags MongoTransmitBinds
// @Accept json
// @Produce json
// @Param MongoTransmitBind body models.MongoTransmitBind true "MongoTransmitBind"
// @Success 200 {object}  servlet.JSONResult{data=models.MongoTransmitBind} "MongoTransmitBind"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "MongoTransmitBind未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /MongoTransmitBind/update [post]
func (api *MongoTransmitBindApi) UpdateMongoTransmitBind(c *gin.Context) {
	var req models.MongoTransmitBind
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.MongoTransmitBind
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "MongoTransmitBind not found")
		return
	}

	var newV models.MongoTransmitBind
	newV = old
	newV.DeviceUid = req.DeviceUid
	newV.Protocol = req.Protocol
	newV.IdentificationCode = req.IdentificationCode
	newV.MongoTransmitId = req.MongoTransmitId
	newV.Collection = req.Collection
	newV.Database = req.Database
	newV.Script = req.Script
	newV.Enable = req.Enable
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	MongoTransmitBindBiz.HandlerRedis(newV)
	servlet.Resp(c, old)
}

// PageMongoTransmitBind
// @Summary 分页查询MongoTransmitBind
// @Description 分页查询MongoTransmitBind
// @Tags MongoTransmitBinds
// @Accept json
// @Produce json
// @Param collection query string false "collection" default(0)
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.MongoTransmitBind}} "MongoTransmitBind"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /MongoTransmitBind/page [get]
func (api *MongoTransmitBindApi) PageMongoTransmitBind(c *gin.Context) {
	var name = c.Query("collection")
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

	data, err := MongoTransmitBindBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteMongoTransmitBind
// @Tags      MongoTransmitBinds
// @Summary   删除MongoTransmitBind
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /MongoTransmitBind/delete/:id [post]
// @Success 200 {object}  servlet.JSONResult{data=string} 
func (api *MongoTransmitBindApi) DeleteMongoTransmitBind(c *gin.Context) {
	var MongoTransmitBind models.MongoTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&MongoTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "MongoTransmitBind not found")

		return
	}

	if result := glob.GDb.Delete(&MongoTransmitBind); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	MongoTransmitBind.Enable = false
	MongoTransmitBindBiz.HandlerRedis(MongoTransmitBind)
	servlet.Resp(c, "删除成功")
}

// ByIdMongoTransmitBind
// @Tags      MongoTransmitBinds
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /MongoTransmitBind/:id [get]
// @Success 200 {object}  servlet.JSONResult{data=models.MongoTransmitBind} 
func (api *MongoTransmitBindApi) ByIdMongoTransmitBind(c *gin.Context) {
	var MongoTransmitBind models.MongoTransmitBind

	param := c.Param("id")

	result := glob.GDb.First(&MongoTransmitBind, param)
	if result.Error != nil {
		servlet.Error(c, "MongoTransmitBind not found")

		return
	}

	servlet.Resp(c, MongoTransmitBind)
}



// MockScript
// @Tags      MongoTransmits
// @Summary   模拟脚本
// @Param MongoTransmit body servlet.TransmitScriptParam true "执行参数"
// @Produce   application/json
// @Router    /MongoTransmitBind/mockScript [post]
// @Success 200 {object}  servlet.JSONResult{data=string} 
func (api *MongoTransmitBindApi) MockScript(c *gin.Context) {
	var req servlet.TransmitScriptParam
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}
	script := 	MongoTransmitBindBiz.MockScript(req.DataRowList, req.Script)
	servlet.Resp(c, script)
}
