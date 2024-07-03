package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type SimCardApi struct{}

var simCardBiz = biz.SimCardBiz{}

// CreateSimCard
// @Summary 创建SIM卡
// @Description 创建SIM卡
// @Tags SimCards
// @Accept json
// @Produce json
// @Param SimCard body models.SimCard true "SIM卡"
// @Success 201 {object} servlet.JSONResult{data=models.SimCard} "创建成功的SIM卡"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /SimCard/create [post]
func (api *SimCardApi) CreateSimCard(c *gin.Context) {
	var SimCard models.SimCard
	if err := c.ShouldBindJSON(&SimCard); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	err := SimCard.Validate()
	if err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := glob.GDb.Create(&SimCard)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的SIM卡
	servlet.Resp(c, SimCard)
}

// UpdateSimCard
// @Summary 更新一个SIM卡
// @Description 更新一个SIM卡
// @Tags SimCards
// @Accept json
// @Produce json
// @Param SimCard body models.SimCard true "SIM卡"
// @Success 200 {object}  servlet.JSONResult{data=models.SimCard} "SIM卡"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "SIM卡未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /SimCard/update [post]
func (api *SimCardApi) UpdateSimCard(c *gin.Context) {
	var req models.SimCard
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.SimCard
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "SimCard not found")
		return
	}

	var newV models.SimCard
	newV = old
	newV.AccessNumber = req.AccessNumber
	newV.ICCID = req.ICCID
	newV.IMSI = req.IMSI
	newV.Operator = req.Operator
	newV.Expiration = req.Expiration
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageSimCard
// @Summary 分页查询SIM卡
// @Description 分页查询SIM卡
// @Tags SimCards
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Param AccessNumber query string false 接入号
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.SimCard}} "SIM卡"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /SimCard/page [get]
func (api *SimCardApi) PageSimCard(c *gin.Context) {
	var accessNumber = c.Query("AccessNumber")
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

	data, err := simCardBiz.PageData(accessNumber, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteSimCard
// @Tags      SimCards
// @Summary   删除SIM卡
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /SimCard/delete/:id [post]
func (api *SimCardApi) DeleteSimCard(c *gin.Context) {
	var SimCard models.SimCard

	param := c.Param("id")

	result := glob.GDb.First(&SimCard, param)
	if result.Error != nil {
		servlet.Error(c, "SimCard not found")

		return
	}

	if result := glob.GDb.Delete(&SimCard); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdSimCard
// @Tags      SimCards
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /SimCard/:id [get]
func (api *SimCardApi) ByIdSimCard(c *gin.Context) {
	var SimCard models.SimCard

	param := c.Param("id")

	result := glob.GDb.First(&SimCard, param)
	if result.Error != nil {
		servlet.Error(c, "SimCard not found")

		return
	}

	servlet.Resp(c, SimCard)
}

// BindDeviceInfo
// @Summary SIM卡绑定设备
// @Description SIM卡绑定设备
// @Tags SimCards
// @Accept json
// @Produce json
// @Param SimCard body models.SimUseHistory true "SIM卡历史"
// @Success 201 {object} servlet.JSONResult{data=models.SimUseHistory} ""
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /SimCard/BindDeviceInfo [post]
func (api *SimCardApi) BindDeviceInfo(c *gin.Context) {
	var SimUseHistory models.SimUseHistory
	if err := c.ShouldBindJSON(&SimUseHistory); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := glob.GDb.Create(&SimUseHistory)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, SimUseHistory)
}

// HistoryList
// @Summary SIM卡绑定设备历史查询
// @Description SIM卡绑定设备历史查询
// @Tags SimCards
// @Accept json
// @Produce json
// @Param sim_id query int false "sim卡id"
// @Param page_size query int false "每页大小" default(10)
// @Param page_size query int false "每页大小" default(10)
// @Success 201 {object} servlet.JSONResult{data=servlet.SimUseHistoryResp} ""
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /SimCard/history [get]
func (api *SimCardApi) HistoryList(c *gin.Context) {
	var simId = c.Query("sim_id")
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

	data, err := simCardBiz.PageHistory(simId, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}
