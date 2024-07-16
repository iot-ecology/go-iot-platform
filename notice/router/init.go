package router

import "github.com/gin-gonic/gin"

func InitRouter(g *gin.RouterGroup) {
	g.GET("/notice/beat",Beat)
}
