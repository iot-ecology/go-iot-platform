package router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"igp/biz"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"log"
	"strconv"
)

type SignalDelayWaringApi struct {
}

var SignalDelayWaringBiz = biz.SignalDelayWaringBiz{}

// CreateSignalDelayWaring
// @Summary 创建脚本报警
// @Description 创建脚本报警
// @Tags signal-delay-waring
// @Accept json
// @Produce json
// @Param SignalDelayWaring body models.SignalDelayWaring true "脚本报警"
// @Success 201 {object} servlet.JSONResult{data=models.SignalDelayWaring} "创建成功的脚本报警"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /signal-delay-waring/create [post]
func (api *SignalDelayWaringApi) CreateSignalDelayWaring(c *gin.Context) {
	var SignalDelayWaring models.SignalDelayWaring
	if err := c.ShouldBindJSON(&SignalDelayWaring); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	// 检查 SignalDelayWaring 是否被正确初始化
	if SignalDelayWaring.Name == "" {
		servlet.Error(c, "名称不能为空")
		return
	}

	result := glob.GDb.Create(&SignalDelayWaring)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	jsonData, err := json.Marshal(SignalDelayWaring)
	if err != nil {
		glob.GLog.Error("json 序列化异常", zap.Any("err", err))
	}
	glob.GRedis.HSet(context.Background(), "signal_delay_config", strconv.Itoa(int(SignalDelayWaring.ID)), jsonData)

	// 返回创建成功的脚本报警
	servlet.Resp(c, SignalDelayWaring)
}

// UpdateSignalDelayWaring
// @Summary 更新一个脚本报警
// @Description 更新一个脚本报警
// @Tags signal-delay-waring
// @Accept json
// @Produce json
// @Param id path int true "脚本报警id"
// @Param SignalDelayWaring body models.SignalDelayWaring true "脚本报警"
// @Success 200 {object}  servlet.JSONResult{data=models.SignalDelayWaring} "脚本报警"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "脚本报警未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /signal-delay-waring/update [post]
func (api *SignalDelayWaringApi) UpdateSignalDelayWaring(c *gin.Context) {
	var req models.SignalDelayWaring
	if err := c.ShouldBindJSON(&req); err != nil {

		servlet.Error(c, err.Error())
		return
	}

	var old models.SignalDelayWaring
	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {

		servlet.Error(c, "SignalDelayWaring not found")
		return
	}

	var newV models.SignalDelayWaring
	newV = old
	newV.Name = req.Name
	newV.Script = req.Script
	result = glob.GDb.Model(&newV).Updates(newV)

	jsonData, err := json.Marshal(newV)
	if err != nil {
		glob.GLog.Error("json 序列化异常", zap.Any("err", err))
	}
	glob.GRedis.HSet(context.Background(), "signal_delay_config", strconv.Itoa(int(newV.ID)), jsonData)

	if result.Error != nil {

		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, old)
	return
}

// PageSignalDelayWaring
// @Summary 分页查询脚本报警
// @Description 分页查询脚本报警
// @Tags signal-delay-waring
// @Accept json
// @Produce json
// @Param name query string false "名称" maxlength(100)
// @Param page query int false "页码" default(0)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.SignalDelayWaring}} "脚本报警"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /signal-delay-waring/page [get]
func (api *SignalDelayWaringApi) PageSignalDelayWaring(c *gin.Context) {
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

	data, err := SignalDelayWaringBiz.PageData(name, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
	return
}

// DeleteSignalDelayWaring
// @Tags      signal-delay-waring
// @Summary   删除脚本报警
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /signal-delay-waring/delete/:id [post]
func (api *SignalDelayWaringApi) DeleteSignalDelayWaring(c *gin.Context) {
	var SignalDelayWaring models.SignalDelayWaring

	param := c.Param("id")

	result := glob.GDb.First(&SignalDelayWaring, param)
	if result.Error != nil {
		servlet.Error(c, "SignalDelayWaring not found")

		return
	}

	if result := glob.GDb.Delete(&SignalDelayWaring); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// Mock
// @Tags      signal-delay-waring
// @Summary   模拟执行
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /signal-delay-waring/Mock/:id [post]
func (api *SignalDelayWaringApi) Mock(c *gin.Context) {

	param := c.Param("id")
	i, err := strconv.Atoi(param)
	if err != nil {
		fmt.Println("转换错误:", err)
		return
	}
	mock := SignalDelayWaringBiz.Mock(i)
	servlet.Resp(c, mock)
}

// GenParam
// @Tags      signal-delay-waring
// @Summary   生产模拟参数
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /signal-delay-waring/GenParam/:id [post]
func (api *SignalDelayWaringApi) GenParam(c *gin.Context) {

	param := c.Param("id")
	i, err := strconv.Atoi(param)
	if err != nil {
		fmt.Println("转换错误:", err)
		return
	}
	_, m := SignalDelayWaringBiz.GenParam(i)
	servlet.Resp(c, m)
}

// QueryWaringList
// @Tags signal-delay-waring
// @Summary   查询报警历史
// @Param config body servlet.WaringRowQuery true "查询参数"
// @Produce   application/json
// @Router    /signal-delay-waring/query-row [post]
func (api *SignalDelayWaringApi) QueryWaringList(c *gin.Context) {

	var req servlet.WaringRowQuery

	if err := c.ShouldBindJSON(&req); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	servlet.Resp(c, query2(req))
	return

}
func query2(req servlet.WaringRowQuery) []bson.M {
	database := glob.GMongoClient.Database(glob.GConfig.MongoConfig.Db)
	collection := database.Collection(glob.GConfig.MongoConfig.ScriptWaringCollection)

	filter := bson.M{
		"rule_id": req.ID,
		"up_time": bson.M{
			"$gte": req.UpTimeStart, // 大于等于开始时间
			"$lte": req.UpTimeEnd,   // 小于等于结束时间
		},
	}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
		}
	}(cur, context.TODO())
	var c []bson.M

	for cur.Next(context.TODO()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		c = append(c, result)
	}
	return c
}
