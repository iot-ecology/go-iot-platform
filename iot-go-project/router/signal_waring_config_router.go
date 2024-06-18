package router

import (
	"context"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"log"
	"strconv"
)

type SignalWaringConfigApi struct{}

// CreateSignalWaringConfig
// @Summary 创建一个新的信号报警配置
// @Description 创建一个新的信号报警配置记录。
// @Tags signal-waring-configs
// @Accept json
// @Produce json
// @Param config body models.SignalWaringConfig true "信号报警配置信息"
// @Success 201 {object} servlet.JSONResult{data=models.SignalWaringConfig}  "创建成功的信号报警配置"
// @Failure 400 {string} string "请求数据错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /signal-waring-config/create [post]
func (api *SignalWaringConfigApi) CreateSignalWaringConfig(c *gin.Context) {
	var config models.SignalWaringConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		servlet.Error(c, err.Error())

		return
	}
	if config.SignalId < 1 {
		servlet.Error(c, "信号ID必填")
		return
	}
	result := glob.GDb.Create(&config)
	bizSignal.SetSignalWaringCache(config)
	if result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}
	servlet.Resp(c, config)
}

// UpdateSignalWaringConfig
// @Summary 更新一个已有的信号报警配置
// @Description 根据ID更新一个已有的信号报警配置记录。
// @Tags signal-waring-configs
// @Accept json
// @Produce json
// @Param config body models.SignalWaringConfig true "更新的信号报警配置信息"
// @Success 200 {object} servlet.JSONResult{data=models.SignalWaringConfig}  "更新成功的信号报警配置"
// @Failure 400 {string} string "请求数据错误"
// @Failure 404 {string} string "信号报警配置未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /signal-waring-config/update [post]
func (api *SignalWaringConfigApi) UpdateSignalWaringConfig(c *gin.Context) {
	var req models.SignalWaringConfig

	if err := c.ShouldBindJSON(&req); err != nil {
		servlet.Error(c, err.Error())
		return
	}
	var old models.SignalWaringConfig

	result := glob.GDb.First(&old, req.ID)
	if result.Error != nil {
		servlet.Error(c, "SignalWaringConfig not found")
		return
	}
	bizSignal.RemoveSignalWaringCache(old)

	var newV models.SignalWaringConfig
	newV = old
	newV.Min = req.Min
	newV.Max = req.Max
	newV.InOrOut = req.InOrOut

	m := structs.Map(newV)
	result = glob.GDb.Table("signal_waring_configs").Where("id = ?", newV.ID).Updates(m)

	if result.Error != nil {
		servlet.Error(c, result.Error.Error())

		return
	}
	bizSignal.SetSignalWaringCache(newV)
	servlet.Resp(c, old)
}

// PageSignalWaringConfig
// @Summary 分页查询信号报警配置
// @Description 根据条件分页查询信号报警配置列表。
// @Tags signal-waring-configs
// @Accept json
// @Produce json
// @Param signal_id query int false "信号ID"
// @Param mqtt_client_id query int false "mqtt客户端表id"
// @Param page query int true "页码" default(1)
// @Param page_size query int true "每页大小" default(10)
// @Success 200 {object} servlet.JSONResult{data=servlet.PaginationQ{data=models.SignalWaringConfig}} "信号报警配置列表"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "查询异常"
// @Router /signal-waring-config/page [get]
func (api *SignalWaringConfigApi) PageSignalWaringConfig(c *gin.Context) {
	var err error

	value := c.Query("signal_id")
	mqttClientId := c.Query("mqtt_client_id")

	atoi, err := strconv.Atoi(value)
	if err != nil {
		servlet.Error(c, "无效的signal_id")

		return
	}

	var page = c.DefaultQuery("page", "1")
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

	data, err := bizSignal.PageSignalWaringConfig(atoi, mqttClientId, parseUint, u)
	if err != nil {
		servlet.Error(c, "查询异常")
		return
	}
	servlet.Resp(c, data)
}

// DeleteSignalWaringConfig
// @Tags signal-waring-configs
// @Summary   删除信号报警配置
// @Param id path int true "主键"
// @Produce   application/json
// @Router    /signal-waring-config/delete/:id [post]
func (api *SignalWaringConfigApi) DeleteSignalWaringConfig(c *gin.Context) {
	var config models.SignalWaringConfig

	param := c.Param("id")

	result := glob.GDb.First(&config, param)
	if result.Error != nil {
		servlet.Error(c, "SignalWaringConfig not found")
		return
	}
	bizSignal.RemoveSignalWaringCache(config)

	if result := glob.GDb.Delete(&config); result.Error != nil {
		servlet.Error(c, result.Error.Error())
		return
	}

	servlet.Resp(c, "删除成功")
}

// QueryWaringList
// @Tags signal-waring-configs
// @Summary   查询报警历史
// @Param config body servlet.WaringRowQuery true "查询参数"
// @Produce   application/json
// @Router    /signal-waring-config/query-row [post]
func (api *SignalWaringConfigApi) QueryWaringList(c *gin.Context) {

	var req servlet.WaringRowQuery

	if err := c.ShouldBindJSON(&req); err != nil {
		servlet.Error(c, err.Error())
		return
	}

	servlet.Resp(c, query(req))

}

func query(req servlet.WaringRowQuery) []bson.M {
	database := glob.GMongoClient.Database(glob.GConfig.MongoConfig.Db)
	collection := database.Collection(glob.GConfig.MongoConfig.WaringCollection)

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
			zap.S().Errorf("err %+v", err)

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
