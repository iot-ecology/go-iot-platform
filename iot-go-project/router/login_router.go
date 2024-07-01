package router

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"igp/servlet"
	"net"
	"time"
)

type LoginApi struct{}

type MyClaims struct {
	Uid      int    // 用户id
	UserName string // 用户名
	RoleIds  []uint // 用户角色id
	jwt.StandardClaims
}

// 定义key todo: 改成配置文件
var jwtKey = []byte("a_secret_test123456")

// 过期时间24小时
var expireTime = time.Now().Add(24 * time.Hour).Unix()

// GenerateToken 函数根据用户ID和用户名生成JWT令牌字符串
//
// 参数：
// - uid：int类型，用户ID
// - username：string类型，用户名
//
// 返回值：
// - string类型，生成的JWT令牌字符串
// - error类型，如果生成令牌过程中发生错误，则返回非零的错误码；否则返回nil
func GenerateToken(uid int, username string, roleId []uint) (string, error) {
	// 创建MyClaims实例
	myClaims := MyClaims{
		uid,      // 用户id
		username, // 用户名
		roleId,   // 用户角色id
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "test",
		},
	}

	// 创建签名对象
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	// 使用指定的jwtKey签名获取完整编码字符串token
	tokenStr, err := tokenObj.SignedString(jwtKey)
	if err != nil {

		return "", err
	}

	// 返回token
	return tokenStr, nil
}

// GetAuthorizationToken 从gin.Context中获取并返回Authorization头部中的token信息
//
// 参数：
//
//	c *gin.Context：gin框架中的上下文对象，包含了HTTP请求和响应的信息
//
// 返回值：
//
//	string：从Authorization头部获取的token信息，如果未找到则返回空字符串
func GetAuthorizationToken(c *gin.Context) string {
	// 获取接口传递过来的Authorization
	tokenInfo := c.Request.Header.Get("Authorization")

	// 返回token
	return tokenInfo
}

// ParseToken 是一个用于解析JWT令牌并返回其Token对象、自定义声明MyClaims以及可能发生的错误的函数
//
// 参数：
// tokenStr string - 待解析的JWT令牌字符串
//
// 返回值：
// *jwt.Token - 解析后的JWT令牌对象，包括Header、Payload和Signature
// *MyClaims - 自定义声明结构体指针，用于存储解析后的JWT令牌中的用户信息
// error - 解析过程中可能出现的错误，如果解析成功则为nil
func ParseToken(tokenStr string) (*jwt.Token, *MyClaims, error) {
	myClaims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, myClaims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, myClaims, err
}

// Login
// @Tags      登录
// @Summary   登录
// @Produce   application/json
// @Param     data  body      servlet.LoginParam true "账号密码"
// @Router    /login [post]
func Login(c *gin.Context) {

	var param servlet.LoginParam

	if err := c.ShouldBindJSON(&param); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	user := userBiz.FindUser(param.UserName, param.Password)
	if user == nil {

		servlet.Error(c, "用户名或密码错误")
		return
	}
	// 走到这说明登录成功,生成token
	tokenStr, err := GenerateToken(int(user.ID), user.Username, roleBiz.FindByUserId(user.ID))

	if err != nil {
		servlet.Error(c, "登录失败")
		return
	}

	servlet.Resp(c, gin.H{
		"token":    tokenStr,
		"uid":      int(user.ID),
		"username": user.Username,
	})
}

// UserInfo
// @Tags      登录
// @Summary   获取用户信息
// @Produce   application/json
// @Param      Authorization  header  string  true  "Authorization token"d
// @Router    /userinfo [post]
func (v LoginApi) UserInfo(c *gin.Context) {
	// 从header获取Authorization
	authorization := GetAuthorizationToken(c)

	// 不存在token 校验token是否有效
	tokenObj, userInfos, err := ParseToken(authorization)

	// 校验authorization是否为空
	if authorization == "" || err != nil {
		servlet.Error(c, "未登录")
		return
	}

	servlet.Resp(c, gin.H{
		"uid":      userInfos.Uid,
		"username": userInfos.UserName,
		"roleIds":  userInfos.RoleIds,
		"token":    tokenObj.Raw,
	})

}

// isLocalIP 判断一个IP地址是否为本地IP
//
// 参数：
//
//	ipStr string - 待判断的IP地址字符串
//
// 返回值：
//
//	bool - 如果IP地址是本地IP，则返回true；否则返回false
func isLocalIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip.IsLoopback() {
		return true
	}
	ip = ip.To4()
	if ip == nil {
		return false
	}

	switch {
	case ip[0] == 10:
		return true
	case ip[0] == 192 && ip[1] == 168:
		return true
	case ip[0] == 172 && (ip[1] >= 16 && ip[1] <= 31):
		return true
	default:
		return false
	}
}

// JwtCheck 是一个返回gin.HandlerFunc的函数，用于在gin框架中作为中间件检查JWT
func JwtCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的IP地址
		clientIP := c.ClientIP()
		// 检查是否是局域网IP
		if isLocalIP(clientIP) {
			// 如果是局域网IP，则跳过JWT验证
			c.Next()
			return
		}

		authorization := GetAuthorizationToken(c)
		tokenObj, userInfos, err := ParseToken(authorization)
		// 校验authorization是否为空或解析token出错
		if authorization == "" || err != nil {
			servlet.Error(c, "未登录")
			c.Abort()
			return
		}
		c.Set("claims", tokenObj)
		c.Set("user_info", userInfos)
		c.Set("role_ids", userInfos.RoleIds)
		c.Next()
	}
}
