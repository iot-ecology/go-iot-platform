package router

import (
	"github.com/gin-gonic/gin"
	"igp/servlet"
	"strconv"
)

type MessageListApi struct{}

// PageMessageList
// @Summary 分页查询消息
// @Description 分页查询消息
// @Tags MessageLists
// @Accept json
// @Produce json
// @Param page query int false "页码" default(0)
// @Param message_type_id query int false "消息类型"
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.MessageList}} "消息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /MessageList/page [get]
func (api *MessageListApi) PageMessageList(c *gin.Context) {
	var messageTypeId = c.Query("message_type_id")
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

	data, err := dashBiz.PageData(messageTypeId, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}
