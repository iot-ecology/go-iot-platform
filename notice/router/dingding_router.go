package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Beat
// 健康检查
// @Summary      健康检查
// @Tags         beat
// @Produce      json
// @Success      20000  {object}  string
// @Router       /bb/beat [get]
func Beat(g *gin.Context) {
	result := JSONResult{}
	result.Message = "操作成功"
	result.Code = 20000
	result.Data = "beat"
	g.JSON(http.StatusOK, result)

}

type JSONResult struct {
	Code    int         `json:"code" `
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
