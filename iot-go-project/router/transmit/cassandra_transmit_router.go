package transmit

import (
	"github.com/gin-gonic/gin"
	"igp/biz/transmit"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type CassandraTransmitApi struct{}

var CassandraTransmitBiz = transmit.CassandraTransmitBiz{}

// CreateCassandraTransmit
// @Summary 创建Cassandra数据库管理
// @Description 创建Cassandra数据库管理
// @Tags CassandraTransmits
// @Accept json
// @Produce json
// @Param CassandraTransmit body models.CassandraTransmit true "Cassandra数据库管理"
// @Success 200 {object} servlet.JSONResult{data=models.CassandraTransmit} "创建成功的Cassandra数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /CassandraTransmit/create [post]
func (api *CassandraTransmitApi) CreateCassandraTransmit(c *gin.Context) {
	var CassandraTransmit models.CassandraTransmit
	if err := c.ShouldBindJSON(&CassandraTransmit); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	result := glob.GDb.Create(&CassandraTransmit)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	CassandraTransmitBiz.SetRedis(CassandraTransmit)
	// 返回创建成功的Cassandra数据库管理
	servlet.Resp(c, CassandraTransmit)
}

// UpdateCassandraTransmit
// @Summary 更新一个Cassandra数据库管理
// @Description 更新一个Cassandra数据库管理
// @Tags CassandraTransmits
// @Accept json
// @Produce json
// @Param CassandraTransmit body models.CassandraTransmit true "Cassandra数据库管理"
// @Success 200 {object}  servlet.JSONResult{data=models.CassandraTransmit} "Cassandra数据库管理"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "Cassandra数据库管理未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /CassandraTransmit/update [post]
func (api *CassandraTransmitApi) UpdateCassandraTransmit(c *gin.Context) {
	var req models.CassandraTransmit
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.CassandraTransmit
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "CassandraTransmit not found")
		return
	}

	var newV models.CassandraTransmit
	newV = old
	newV.Name = req.Name
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	CassandraTransmitBiz.SetRedis(newV)

	servlet.Resp(c, old)
}

// PageCassandraTransmit
// @Summary 分页查询Cassandra数据库管理
// @Description 分页查询Cassandra数据库管理
// @Tags CassandraTransmits
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.CassandraTransmit}} "Cassandra数据库管理"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /CassandraTransmit/page [get]
func (api *CassandraTransmitApi) PageCassandraTransmit(c *gin.Context) {
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

	data, err := CassandraTransmitBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteCassandraTransmit
// @Tags      CassandraTransmits
// @Summary   删除Cassandra数据库管理
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /CassandraTransmit/delete/:id [post]
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=string}

func (api *CassandraTransmitApi) DeleteCassandraTransmit(c *gin.Context) {
	var CassandraTransmit models.CassandraTransmit

	param := c.Param("id")

	result := glob.GDb.First(&CassandraTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "CassandraTransmit not found")

		return
	}

	if result := glob.GDb.Delete(&CassandraTransmit); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	CassandraTransmitBiz.DeleteRedis(CassandraTransmit)
	servlet.Resp(c, "删除成功")
}

// ByIdCassandraTransmit
// @Tags      CassandraTransmits
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /CassandraTransmit/:id [get]
// @Success 200 {object}  servlet.JSONResult{data=models.CassandraTransmit}
func (api *CassandraTransmitApi) ByIdCassandraTransmit(c *gin.Context) {
	var CassandraTransmit models.CassandraTransmit

	param := c.Param("id")

	result := glob.GDb.First(&CassandraTransmit, param)
	if result.Error != nil {
		servlet.Error(c, "CassandraTransmit not found")

		return
	}

	servlet.Resp(c, CassandraTransmit)
}

// ListCassandraTransmit
// @Tags      CassandraTransmits
// @Summary   列表
// @Produce   application/json
// @Router    /CassandraTransmit/list [get]
// @Success 200 {object}  servlet.JSONResult{data=models.CassandraTransmit[]}
func (api *CassandraTransmitApi) ListCassandraTransmit(c *gin.Context) {
	var CassandraTransmit []models.CassandraTransmit

	glob.GDb.Find(&CassandraTransmit)

	servlet.Resp(c, CassandraTransmit)
}
