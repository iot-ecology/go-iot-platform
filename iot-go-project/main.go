package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"igp/docs"
	"igp/glob"
	"igp/initialize"
	"igp/servlet"
	"igp/task"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

// Beat
// 健康检查
// @Summary      健康检查
// @Tags         beat
// @Produce      json
// @Success      20000  {object}  string
// @Router       /beat [get]
func Beat(g *gin.Context) {
	result := servlet.JSONResult{}
	result.Message = "操作成功"
	result.Code = 20000
	result.Data = "beat"
	g.JSON(http.StatusOK, result)

}

// @title Go语言物联网
// @version 1.0
// @description 基于Go语言的物联网项目

// @contact.name Zen Huifer
// @contact.url https://github.com/huifer
// @contact.email huifer97@163.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host http://127.0.0.1:8080/
// @BasePath /
func main() {

	r := gin.Default()

	r.Use(CORSMiddleware())
	r.Use(ExceptionMiddleware)
	initialize.InitAll(r)
	task.InitTask()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/beat", Beat)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":" + strconv.Itoa(glob.GConfig.NodeInfo.Port))

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func ExceptionMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 简单返回友好提示，具体可自定义发生错误后处理逻辑
			glob.GLog.Sugar().Error(err)

			servlet.Error(c, err)
			c.Abort()
		}
	}()
	c.Next()
}
