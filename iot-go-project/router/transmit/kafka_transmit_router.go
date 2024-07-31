package transmit

import (
	"github.com/gin-gonic/gin"
	"igp/biz/transmit"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type KafkaTransmitApi struct{}

var KafkaTransmitBiz = transmit.KafkaTransmitBiz{}

// CreateKafkaTransmit
// @Summary 创建Kafka数据库管理
// @Description 创建Kafka数据库管理
// @Tags KafkaTransmits
// @Accept json
// @Produce json
// @Param KafkaTransmit body models.KafkaTransmit true "Kafka数据库管理"
// @Success 201 {object} servlet.JSONResult{data=models.KafkaTransmit} "创建成功的Kafka数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /KafkaTransmit/create [post]
func (api *KafkaTransmitApi) CreateKafkaTransmit(c *gin.Context) {
	var KafkaTransmit models.KafkaTransmit
	if err := c.ShouldBindJSON(&KafkaTransmit); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 KafkaTransmit 是否被正确初始化
	if KafkaTransmit.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&KafkaTransmit)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	KafkaTransmitBiz.SetRedis(KafkaTransmit)
	// 返回创建成功的Kafka数据库管理
	servlet.Resp(c, KafkaTransmit)
}

// UpdateKafkaTransmit
// @Summary 更新一个Kafka数据库管理
// @Description 更新一个Kafka数据库管理
// @Tags KafkaTransmits
// @Accept json
// @Produce json
// @Param KafkaTransmit body models.KafkaTransmit true "Kafka数据库管理"
// @Success 200 {object}  servlet.JSONResult{data=models.KafkaTransmit} "Kafka数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Kafka数据库管理未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /KafkaTransmit/update [post]
func (api *KafkaTransmitApi) UpdateKafkaTransmit(c *gin.Context) {
	var req models.KafkaTransmit
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.KafkaTransmit
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "KafkaTransmit not found")
		return
	}

	var newV models.KafkaTransmit
	newV = old
	newV.Name = req.Name
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	KafkaTransmitBiz.SetRedis(newV)
	servlet.Resp(c, old)
}

// PageKafkaTransmit
// @Summary 分页查询Kafka数据库管理
// @Description 分页查询Kafka数据库管理
// @Tags KafkaTransmits
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.KafkaTransmit}} "Kafka数据库管理"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /KafkaTransmit/page [get]
func (api *KafkaTransmitApi) PageKafkaTransmit(c *gin.Context) {
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

	data, err := KafkaTransmitBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteKafkaTransmit
// @Tags      KafkaTransmits
// @Summary   删除Kafka数据库管理
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /KafkaTransmit/delete/:id [post]
func (api *KafkaTransmitApi) DeleteKafkaTransmit(c *gin.Context) {
	var KafkaTransmit models.KafkaTransmit

	param := c.Param("id")

	result := glob.GDb.First(&KafkaTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "KafkaTransmit not found")

		return
	}

	if result := glob.GDb.Delete(&KafkaTransmit); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	KafkaTransmitBiz.DeleteRedis(KafkaTransmit)
	servlet.Resp(c, "删除成功")
}

// ByIdKafkaTransmit
// @Tags      KafkaTransmits
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /KafkaTransmit/:id [get]
func (api *KafkaTransmitApi) ByIdKafkaTransmit(c *gin.Context) {
	var KafkaTransmit models.KafkaTransmit

	param := c.Param("id")

	result := glob.GDb.First(&KafkaTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "KafkaTransmit not found")

		return
	}

	servlet.Resp(c, KafkaTransmit)
}

//// MockScript
//// @Tags      KafkaTransmits
//// @Summary   模拟脚本
//// @Param KafkaTransmit body servlet.TransmitScriptParam true "执行参数"
//// @Produce   application/json
//// @Router    /KafkaTransmit/mockScript [post]
//func (api *KafkaTransmitApi) MockScript(c *gin.Context) {
//	var req servlet.TransmitScriptParam
//	if err := c.ShouldBindJSON(&req); err != nil {
//
//		servlet.Error(c, err.Error())
//		return
//	}
//	script := KafkaTransmit.MockScript(req.DataRowList, req.Script)
//	servlet.Resp(c, script)
//}
