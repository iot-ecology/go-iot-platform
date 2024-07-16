package notice

import (
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"igp/biz/notice"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type DingDingApi struct{}

var DingDingBiz = notice.DingDingBiz{}

// CreateDingDing
// @Summary 创建钉钉通道
// @Description 创建钉钉通道
// @Tags DingDings
// @Accept json
// @Produce json
// @Param DingDing body models.DingDing true "钉钉通道"
// @Success 201 {object} servlet.JSONResult{data=models.DingDing} "创建成功的钉钉通道"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /DingDing/create [post]
func (api *DingDingApi) CreateDingDing(c *gin.Context) {
	var DingDing models.DingDing
	if err := c.ShouldBindJSON(&DingDing); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 DingDing 是否被正确初始化
	if DingDing.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&DingDing)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的钉钉通道
	servlet.Resp(c, DingDing)
}

// UpdateDingDing
// @Summary 更新一个钉钉通道
// @Description 更新一个钉钉通道
// @Tags DingDings
// @Accept json
// @Produce json
// @Param DingDing body models.DingDing true "钉钉通道"
// @Success 200 {object}  servlet.JSONResult{data=models.DingDing} "钉钉通道"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "钉钉通道未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /DingDing/update [post]
func (api *DingDingApi) UpdateDingDing(c *gin.Context) {
	var req models.DingDing
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.DingDing
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "DingDing not found")
		return
	}

	var newV models.DingDing
	newV = old
	newV.Name = req.Name

	m := structs.Map(newV)
	result = glob.GDb.Table("DingDings").Where("id = ?", newV.ID).Updates(m)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageDingDing
// @Summary 分页查询钉钉通道
// @Description 分页查询钉钉通道
// @Tags DingDings
// @Accept json
// @Produce json
// @Param name query string false "名称"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.DingDing}} "钉钉通道"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /DingDing/page [get]
func (api *DingDingApi) PageDingDing(c *gin.Context) {
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

	data, err := DingDingBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteDingDing
// @Tags      DingDings
// @Summary   删除钉钉通道
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /DingDing/delete/:id [post]
func (api *DingDingApi) DeleteDingDing(c *gin.Context) {
	var DingDing models.DingDing

	param := c.Param("id")

	result := glob.GDb.First(&DingDing, param)
	if result.Error != nil {
		servlet.Error(c, "DingDing not found")

		return
	}

	if result := glob.GDb.Delete(&DingDing); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdDingDing
// @Tags      DingDings
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /DingDing/:id [get]
func (api *DingDingApi) ByIdDingDing(c *gin.Context) {
	var DingDing models.DingDing

	param := c.Param("id")

	result := glob.GDb.First(&DingDing, param)
	if result.Error != nil {
		servlet.Error(c, "DingDing not found")

		return
	}

	servlet.Resp(c, DingDing)
}

// Bind
// @Tags      DingDings
// @Summary   绑定钉钉通道
// @Param DingDing body []models.DingDingBindProduct true "钉钉通道"
// @Produce   application/json
// @Router    /DingDing/bind [post]
func (api *DingDingApi) Bind(c *gin.Context) {
	var req = []models.DingDingBindProduct{}
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	if DingDingBiz.Bind(req) {

		servlet.Resp(c, "绑定成功")
		return
	} else {
		servlet.Error(c, "绑定失败")
		return
	}
}
