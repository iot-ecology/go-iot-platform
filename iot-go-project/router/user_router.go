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

type UserApi struct{}

var userBiz = biz.UserBiz{}

// CreateUser
// @Summary 创建用户
// @Description 创建用户
// @Tags Users
// @Accept json
// @Produce json
// @Param User body models.User true "用户"
// @Success 201 {object} servlet.JSONResult{data=models.User} "创建成功的用户"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /User/create [post]
func (api *UserApi) CreateUser(c *gin.Context) {
	var User models.User
	if err := c.ShouldBindJSON(&User); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 User 是否被正确初始化
	if User.Username == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&User)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	// 返回创建成功的用户
	servlet.Resp(c, User)
}

// UpdateUser
// @Summary 更新一个用户
// @Description 更新一个用户
// @Tags Users
// @Accept json
// @Produce json
// @Param User body models.User true "用户"
// @Success 200 {object}  servlet.JSONResult{data=models.User} "用户"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "用户未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /User/update [post]
func (api *UserApi) UpdateUser(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.User
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "User not found")
		return
	}

	var newV models.User
	newV = old
	newV.Password = req.Password
	result = glob.GDb.Model(&newV).Updates(newV)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, old)
}

// PageUser
// @Summary 分页查询用户
// @Description 分页查询用户
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.User}} "用户"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /User/page [get]
func (api *UserApi) PageUser(c *gin.Context) {
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

	data, err := userBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteUser
// @Tags      Users
// @Summary   删除用户
// @Produce   application/json
// @Param id path int true "主键"
// @Router    /User/delete/:id [post]
func (api *UserApi) DeleteUser(c *gin.Context) {
	var User models.User

	param := c.Param("id")

	result := glob.GDb.First(&User, param)
	if result.Error != nil {
		servlet.Error(c, "User not found")

		return
	}

	if result := glob.GDb.Delete(&User); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// ByIdUser
// @Tags      Users
// @Summary   单个详情
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /User/:id [get]
func (api *UserApi) ByIdUser(c *gin.Context) {
	var User models.User

	param := c.Param("id")

	result := glob.GDb.First(&User, param)
	if result.Error != nil {
		servlet.Error(c, "User not found")

		return
	}

	servlet.Resp(c, User)
}

// ListUser
// @Tags      Users
// @Summary   用户列表
// @Produce   application/json
// @Router    /User/list [get]
func (api *UserApi) ListUser(c *gin.Context) {
	var users []models.User
	result := glob.GDb.Find(&users)
	if result.Error != nil {
		zap.S().Errorln("Error occurred during querying users:", result.Error)

		servlet.Error(c, "User not found")
		return
	}
	servlet.Resp(c, users)
}

// BindRole
// @Tags      Users
// @Summary   用户绑定角色
// @Param User body servlet.UserBindRoleParam true "绑定参数"
// @Produce   application/json
// @Router    /User/BindRole [post]
func (api *UserApi) BindRole(c *gin.Context) {

	var param servlet.UserBindRoleParam

	if err := c.ShouldBindJSON(&param); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	tx := glob.GDb.Begin()
	if tx.Error != nil {
		servlet.Error(c, "Failed to begin transaction")
		return
	}

	result := tx.Where("`user_id` = ?", param.UserId).Delete(&models.UserRole{})
	if result.Error != nil {
		// 如果出现错误，回滚事务
		tx.Rollback()
		servlet.Error(c, "Error occurred during deletion")
		return
	}

	var userRoles []models.UserRole
	for _, roleId := range param.RoleIds {
		userRoles = append(userRoles, models.UserRole{
			UserId: uint(param.UserId),
			RoleId: uint(roleId),
		})
	}

	result = tx.Model(&models.UserRole{}).CreateInBatches(userRoles, len(userRoles))
	if result.Error != nil {
		tx.Rollback()
		zap.S().Infoln("Error occurred during creation:", result.Error)
		servlet.Error(c, "Error occurred during creation")
		return
	}
	if err := tx.Commit().Error; err != nil {
		servlet.Error(c, "Failed to commit transaction")
		return
	}

	servlet.Resp(c, "绑定成功")

}

// QueryBindRole
// @Tags      Users
// @Summary   查询绑定角色
// @Param user_id path int true "主键"
// @Produce   application/json
// @Router    /User/QueryBindRole [post]
func (api *UserApi) QueryBindRole(c *gin.Context) {
	param := c.Param("user_id")

	var userRoles []models.UserRole

	// 使用 Where 和 Find 方法查询记录
	result := glob.GDb.Where("`user_id` = ?", param).Find(&userRoles)
	if result.Error != nil {
		zap.S().Infoln("Error occurred during query:", result.Error)
		servlet.Error(c, "暂无数据")
		return
	}
	servlet.Resp(c, userRoles)
}

// QueryBindDeviceInfo
// @Tags      Users
// @Summary   查询绑定设备
// @Param user_id path int true "主键"
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.UserBindDeviceInfo}} "绑定关系"
// @Produce   application/json
// @Router    /User/QueryBindDeviceInfo [post]
func (api *UserApi) QueryBindDeviceInfo(c *gin.Context) {
	param := c.Param("user_id")

	var bindDeviceInfos []models.UserBindDeviceInfo

	// 使用 Where 和 Find 方法查询记录
	result := glob.GDb.Where("`user_id` = ?", param).Find(&bindDeviceInfos)
	if result.Error != nil {
		zap.S().Infoln("Error occurred during query:", result.Error)
		servlet.Error(c, "暂无数据")
		return
	}
	servlet.Resp(c, bindDeviceInfos)
}

// BindDeviceInfo
// @Tags      Users
// @Summary   用户绑定设备
// @Param User body servlet.UserBindDeviceInfoParam true "绑定参数"
// @Produce   application/json
// @Router    /User/BindDeviceInfo [post]
func (api *UserApi) BindDeviceInfo(c *gin.Context) {

	var param servlet.UserBindDeviceInfoParam

	if err := c.ShouldBindJSON(&param); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	tx := glob.GDb.Begin()
	if tx.Error != nil {
		servlet.Error(c, "Failed to begin transaction")
		return
	}

	result := tx.Where("`user_id` = ?", param.UserId).Delete(&models.UserBindDeviceInfo{})
	if result.Error != nil {
		// 如果出现错误，回滚事务
		tx.Rollback()
		servlet.Error(c, "Error occurred during deletion")
		return
	}

	var bindDeviceInfos []models.UserBindDeviceInfo
	for _, roleId := range param.DeviceInfoIds {
		bindDeviceInfos = append(bindDeviceInfos, models.UserBindDeviceInfo{
			UserId:   uint(param.UserId),
			DeviceId: uint(roleId),
		})
	}

	result = tx.Model(&models.UserBindDeviceInfo{}).CreateInBatches(bindDeviceInfos, len(bindDeviceInfos))
	if result.Error != nil {
		tx.Rollback()
		zap.S().Infoln("Error occurred during creation:", result.Error)
		servlet.Error(c, "Error occurred during creation")
		return
	}
	if err := tx.Commit().Error; err != nil {
		servlet.Error(c, "Failed to commit transaction")
		return
	}

	servlet.Resp(c, "绑定成功")

}
