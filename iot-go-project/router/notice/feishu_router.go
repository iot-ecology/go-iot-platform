package notice

import (
	"context"
	"encoding/json"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"igp/biz/notice"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type FeiShuApi struct{}

var FeiShuBiz = notice.FeiShuBiz{}

// CreateFeiShu
// @Summary 创建飞书通道
// @Description 创建飞书通道
// @Tags FeiShus
// @Accept json
// @Produce json
// @Param FeiShuId body models.FeiShu true "飞书通道"
// @Success 201 {object} servlet.JSONResult{data=models.FeiShu} "创建成功的飞书通道"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /FeiShuId/create [post]
func (api *FeiShuApi) CreateFeiShu(c *gin.Context) {
	var FeiShu models.FeiShu
	if err := c.ShouldBindJSON(&FeiShu); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 FeiShuId 是否被正确初始化
	if FeiShu.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&FeiShu)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	jsonData, _ := json.Marshal(FeiShu)
	glob.GRedis.HSet(context.Background(), "message_channel_info:feishu", strconv.Itoa(int(FeiShu.ID)), jsonData)
	// 返回创建成功的飞书通道
	servlet.Resp(c, FeiShu)
}

// UpdateFeiShu
// @Summary 更新一个飞书通道
// @Description 更新一个飞书通道
// @Tags FeiShus
// @Accept json
// @Produce json
// @Param FeiShuId body models.FeiShu true "飞书通道"
// @Success 200 {object}  servlet.JSONResult{data=models.FeiShu} "飞书通道"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "飞书通道未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /FeiShuId/update [post]
func (api *FeiShuApi) UpdateFeiShu(c *gin.Context) {
	var req models.FeiShu
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.FeiShu
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "FeiShuId not found")
		return
	}

	var newV models.FeiShu
	newV = old
	newV.Name = req.Name

	m := structs.Map(newV)
	result = glob.GDb.Table("FeiShus").Where("id = ?", newV.ID).Updates(m)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	jsonData, _ := json.Marshal(newV)

	glob.GRedis.HSet(context.Background(), "message_channel_info:feishu", strconv.Itoa(int(newV.ID)), jsonData)
	servlet.Resp(c, old)
}

// PageFeiShu
// @Summary 分页查询飞书通道
// @Description 分页查询飞书通道
// @Tags FeiShus
// @Accept json
// @Produce json
// @Param name query string false "名称"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.FeiShu}} "飞书通道"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /FeiShuId/page [get]
func (api *FeiShuApi) PageFeiShu(c *gin.Context) {
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

	data, err := FeiShuBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteFeiShu
// @Tags      FeiShus
// @Summary   删除飞书通道
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /FeiShuId/delete/:id [post]
func (api *FeiShuApi) DeleteFeiShu(c *gin.Context) {
	var FeiShu models.FeiShu

	param := c.Param("id")

	result := glob.GDb.First(&FeiShu, param)
	if result.Error != nil {
		servlet.Error(c, "FeiShuId not found")

		return
	}

	if result := glob.GDb.Delete(&FeiShu); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	glob.GRedis.HDel(context.Background(), "message_channel_info:feishu", strconv.Itoa(int(FeiShu.ID)))
	servlet.Resp(c, "删除成功")
}

// ByIdFeiShu
// @Tags      FeiShus
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /FeiShuId/:id [get]
func (api *FeiShuApi) ByIdFeiShu(c *gin.Context) {
	var FeiShu models.FeiShu

	param := c.Param("id")

	result := glob.GDb.First(&FeiShu, param)
	if result.Error != nil {
		servlet.Error(c, "FeiShuId not found")

		return
	}

	servlet.Resp(c, FeiShu)
}

// Bind
// @Tags      FeiShus
// @Summary   绑定飞书通道
// @Param FeiShuId body []models.FeiShuBindProduct true "飞书通道"
// @Produce   application/json
// @Router    /FeiShuId/bind [post]
func (api *FeiShuApi) Bind(c *gin.Context) {
	var req []models.FeiShuBindProduct
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	if FeiShuBiz.Bind(req) {

		servlet.Resp(c, "绑定成功")
		return
	} else {
		servlet.Error(c, "绑定失败")
		return
	}
}
