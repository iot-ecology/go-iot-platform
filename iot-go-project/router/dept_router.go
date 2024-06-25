package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type DeptApi struct{}

var deptBiz = biz.DeptBiz{}

// CreateDept
// @Summary 创建部门
// @Description 创建部门
// @Tags Depts
// @Accept json
// @Produce json
// @Param Dept body models.Dept true "部门"
// @Success 201 {object} servlet.JSONResult{data=models.Dept} "创建成功的部门"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /Dept/create [post]
func (api *DeptApi) CreateDept(c *gin.Context) {
	var Dept models.Dept
	if err := c.ShouldBindJSON(&Dept); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 Dept 是否被正确初始化
	if Dept.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&Dept)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的部门
	servlet.Resp(c, Dept)
}

// UpdateDept
// @Summary 更新一个部门
// @Description 更新一个部门
// @Tags Depts
// @Accept json
// @Produce json
// @Param Dept body models.Dept true "部门"
// @Success 200 {object}  servlet.JSONResult{data=models.Dept} "部门"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "部门未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /Dept/update [post]
func (api *DeptApi) UpdateDept(c *gin.Context) {
	var req models.Dept
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.Dept
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "Dept not found")
		return
	}

	var newV models.Dept
	newV = old
	newV.Name = req.Name
	newV.ParentId = req.ParentId
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageDept
// @Summary 分页查询部门
// @Description 分页查询部门
// @Tags Depts
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.Dept}} "部门"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /Dept/page [get]
func (api *DeptApi) PageDept(c *gin.Context) {
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

	data, err := deptBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteDept
// @Tags      Depts
// @Summary   删除部门
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /Dept/delete/:id [post]
func (api *DeptApi) DeleteDept(c *gin.Context) {
	var Dept models.Dept

	param := c.Param("id")

	result := glob.GDb.First(&Dept, param)
	if result.Error != nil {
		servlet.Error(c, "Dept not found")

		return
	}

	if result := glob.GDb.Delete(&Dept); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdDept
// @Tags      Depts
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /Dept/:id [get]
func (api *DeptApi) ByIdDept(c *gin.Context) {
	var Dept models.Dept

	param := c.Param("id")

	result := glob.GDb.First(&Dept, param)
	if result.Error != nil {
		servlet.Error(c, "Dept not found")

		return
	}

	servlet.Resp(c, Dept)
}

// FindByIdSubs
// @Tags      Depts
// @Summary   根据ID查询下级
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /Dept/subs [get]
func (api *DeptApi) FindByIdSubs(c *gin.Context) {
	param := c.Param("id")
	var subDepts []models.Dept
	// 使用手动查询来找到所有父ID为parentID的部门
	result := glob.GDb.Model(&models.Dept{}).Where("parent_id = ?", param).Find(&subDepts)
	if result.Error != nil {
		zap.S().Errorw("查询异常", "error", result.Error)
		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, subDepts)
}
