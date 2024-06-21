package router

import (
	"github.com/gin-gonic/gin"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type ProductApi struct{}

var productBiz = biz.ProductBiz{}

// CreateProduct
// @Summary 创建产品
// @Description 创建产品
// @Tags Products
// @Accept json
// @Produce json
// @Param Product body models.Product true "产品"
// @Success 201 {object} servlet.JSONResult{data=models.Product} "创建成功的产品"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /product/create [post]
func (api *ProductApi) CreateProduct(c *gin.Context) {
	var Product models.Product
	if err := c.ShouldBindJSON(&Product); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 Product 是否被正确初始化
	if Product.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&Product)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的产品
	servlet.Resp(c, Product)
}

// UpdateProduct
// @Summary 更新一个产品
// @Description 更新一个产品
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "产品id"
// @Param Product body models.Product true "产品"
// @Success 200 {object}  servlet.JSONResult{data=models.Product} "产品"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "产品未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /product/update [post]
func (api *ProductApi) UpdateProduct(c *gin.Context) {
	var req models.Product
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.Product
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "Product not found")
		return
	}

	var newV models.Product
	newV = old
	newV.Name = req.Name
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageProduct
// @Summary 分页查询产品
// @Description 分页查询产品
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.Product}} "产品"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /product/page [get]
func (api *ProductApi) PageProduct(c *gin.Context) {
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

	data, err := productBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteProduct
// @Tags      Products
// @Summary   删除产品
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /product/delete/:id [post]
func (api *ProductApi) DeleteProduct(c *gin.Context) {
	var Product models.Product

	param := c.Param("id")

	result := glob.GDb.First(&Product, param)
	if result.Error != nil {
		servlet.Error(c, "Product not found")

		return
	}

	if result := glob.GDb.Delete(&Product); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdProduct
// @Tags      Products
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /product/:id [get]
func (api *ProductApi) ByIdProduct(c *gin.Context) {
	var Product models.Product

	param := c.Param("id")

	result := glob.GDb.First(&Product, param)
	if result.Error != nil {
		servlet.Error(c, "Product not found")

		return
	}

	servlet.Resp(c, Product)
}
