package router

import (
	"errors"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type RoleApi struct{}

var roleBiz = biz.RoleBiz{}

// CreateRole
// @Summary 创建面板
// @Description 创建面板
// @Tags Roles
// @Accept json
// @Produce json
// @Param Role body models.Role true "面板"
// @Success 201 {object} servlet.JSONResult{data=models.Role} "创建成功的面板"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /Role/create [post]
func (api *RoleApi) CreateRole(c *gin.Context) {
	var Role models.Role
	if err := c.ShouldBindJSON(&Role); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 Role 是否被正确初始化
	if Role.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	var qRole models.Role
	err := glob.GDb.Where("name = ?", Role.Name).First(&qRole).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result := glob.GDb.Create(&Role)

		if result.Error != nil {
			servlet.Error(c, result.Error.Error())
			return
		}
		// 返回创建成功的面板
		servlet.Resp(c, Role)
	} else {
		if err != nil {
			servlet.Error(c, err.Error())
			return
		}
		if qRole.ID != 0 {
			servlet.Error(c, "角色已存在")
			return
		}

	}
}

// UpdateRole
// @Summary 更新一个面板
// @Description 更新一个面板
// @Tags Roles
// @Accept json
// @Produce json
// @Param Role body models.Role true "面板"
// @Success 200 {object}  servlet.JSONResult{data=models.Role} "面板"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "面板未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /Role/update [post]
func (api *RoleApi) UpdateRole(c *gin.Context) {
	var req models.Role
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.Role
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "Role not found")
		return
	}

	var newV models.Role
	newV = old
	newV.Name = req.Name
	newV.Description = req.Description

	m := structs.Map(newV)
	result = glob.GDb.Table("roles").Where("id = ?", newV.ID).Updates(m)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageRole
// @Summary 分页查询面板
// @Description 分页查询面板
// @Tags Roles
// @Accept json
// @Produce json
// @Param name query string false "名称"
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.Role}} "面板"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /Role/page [get]
func (api *RoleApi) PageRole(c *gin.Context) {
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

	data, err := roleBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteRole
// @Tags      Roles
// @Summary   删除面板
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /Role/delete/:id [post]
func (api *RoleApi) DeleteRole(c *gin.Context) {
	var Role models.Role

	param := c.Param("id")

	result := glob.GDb.First(&Role, param)
	if result.Error != nil {
		servlet.Error(c, "Role not found")

		return
	}

	if result := glob.GDb.Delete(&Role); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdRole
// @Tags      Roles
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /Role/:id [get]
func (api *RoleApi) ByIdRole(c *gin.Context) {
	var Role models.Role

	param := c.Param("id")

	result := glob.GDb.First(&Role, param)
	if result.Error != nil {
		servlet.Error(c, "Role not found")

		return
	}

	servlet.Resp(c, Role)
}

// ListRole
// @Tags      Roles
// @Summary   角色列表
// @Produce   application/json
// @Router    /Role/list [get]
func (api *RoleApi) ListRole(c *gin.Context) {
	var roles []models.Role
	result := glob.GDb.Find(&roles)
	if result.Error != nil {
		zap.S().Errorln("Error occurred during querying roles:", result.Error)

		servlet.Error(c, "Role not found")
		return
	}
	servlet.Resp(c, roles)
}
