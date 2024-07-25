package initialize

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"igp/glob"
	"igp/models"
	"igp/router"
	"igp/router/notice"
	"igp/router/transmit"
	"log"
	"net/url"
	"os"
	"syscall"
	"time"
)

var (
	mqttApi                   = router.MqttApi{}
	signalApi                 = router.SignalApi{}
	signalWaringConfigApi     = router.SignalWaringConfigApi{}
	influxdbApi               = router.InfluxDbApi{}
	dashboardApi              = router.DashboardApi{}
	calcRuleApi               = router.CalcRuleApi{}
	calcParamApi              = router.CalcParamApi{}
	signalDelayWaringParamApi = router.SignalDelayWaringParamApi{}
	signalDelayWaringApi      = router.SignalDelayWaringApi{}
	productApi                = router.ProductApi{}
	deviceGroupApi            = router.DeviceGroupApi{}
	deviceInfoApi             = router.DeviceInfoApi{}
	productionPlanApi         = router.ProductionPlanApi{}
	repairRecordApi           = router.RepairRecordApi{}
	fileApi                   = router.FileApi{}
	userApi                   = router.UserApi{}
	deptApi                   = router.DeptApi{}
	roleApi                   = router.RoleApi{}
	shipmentRecordApi         = router.ShipmentRecordApi{}
	loginApi                  = router.LoginApi{}
	messageListApi            = router.MessageListApi{}
	simCardApi                = router.SimCardApi{}
	mysqlTransmitApi          = transmit.MySQLTransmitApi{}
	mongoTransmitApi          = transmit.MongoTransmitApi{}
	influxdbTransmitApi       = transmit.InfluxdbTransmitApi{}
	clickTransmitApi          = transmit.ClickhouseTransmitApi{}
	cassandraTransmitApi      = transmit.CassandraTransmitApi{}

	feishuApi   = notice.FeiShuApi{}
	dingdingApi = notice.DingDingApi{}

	tcpHandlerApi  = router.TcpHandlerApi{}
	httpHandlerApi = router.HttpHandlerApi{}
)

func initTable() {
	if !glob.GDb.Migrator().HasTable(&models.MqttClient{}) {

		err := glob.GDb.AutoMigrate(&models.MqttClient{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.Signal{}) {

		err := glob.GDb.AutoMigrate(&models.Signal{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)

		}
	}
	if !glob.GDb.Migrator().HasTable(&models.SignalWaringConfig{}) {

		err := glob.GDb.AutoMigrate(&models.SignalWaringConfig{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)

		}
	}
	if !glob.GDb.Migrator().HasTable(&models.SignalDelayWaring{}) {

		err := glob.GDb.AutoMigrate(&models.SignalDelayWaring{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)

		}
	}

	if !glob.GDb.Migrator().HasTable(&models.SignalDelayWaringParam{}) {

		err := glob.GDb.AutoMigrate(&models.SignalDelayWaringParam{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)

		}
	}

	if !glob.GDb.Migrator().HasTable(&models.CalcRule{}) {

		err := glob.GDb.AutoMigrate(&models.CalcRule{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)

		}
	}
	if !glob.GDb.Migrator().HasTable(&models.CalcParam{}) {

		err := glob.GDb.AutoMigrate(&models.CalcParam{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)

		}
	}
	if !glob.GDb.Migrator().HasTable(&models.Dashboard{}) {

		err := glob.GDb.AutoMigrate(&models.Dashboard{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.Product{}) {

		err := glob.GDb.AutoMigrate(&models.Product{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.DeviceInfo{}) {

		err := glob.GDb.AutoMigrate(&models.DeviceInfo{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.DeviceGroup{}) {

		err := glob.GDb.AutoMigrate(&models.DeviceGroup{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.DeviceGroupDevice{}) {

		err := glob.GDb.AutoMigrate(&models.DeviceGroupDevice{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.RepairRecord{}) {

		err := glob.GDb.AutoMigrate(&models.RepairRecord{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.ProductionPlan{}) {

		err := glob.GDb.AutoMigrate(&models.ProductionPlan{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.ProductPlan{}) {

		err := glob.GDb.AutoMigrate(&models.ProductPlan{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}

	if !glob.GDb.Migrator().HasTable(&models.User{}) {

		err := glob.GDb.AutoMigrate(&models.User{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.Role{}) {

		err := glob.GDb.AutoMigrate(&models.Role{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.UserRole{}) {

		err := glob.GDb.AutoMigrate(&models.UserRole{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.Dept{}) {

		err := glob.GDb.AutoMigrate(&models.Dept{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.ShipmentRecord{}) {

		err := glob.GDb.AutoMigrate(&models.ShipmentRecord{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.ShipmentProductDetail{}) {

		err := glob.GDb.AutoMigrate(&models.ShipmentProductDetail{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.UserBindDeviceInfo{}) {

		err := glob.GDb.AutoMigrate(&models.UserBindDeviceInfo{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.DeviceBindMqttClient{}) {

		err := glob.GDb.AutoMigrate(&models.DeviceBindMqttClient{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.DeviceGroupBindMqttClient{}) {

		err := glob.GDb.AutoMigrate(&models.DeviceGroupBindMqttClient{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.MessageTypeBindRole{}) {

		err := glob.GDb.AutoMigrate(&models.MessageTypeBindRole{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.MessageList{}) {

		err := glob.GDb.AutoMigrate(&models.MessageList{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.SimCard{}) {

		err := glob.GDb.AutoMigrate(&models.SimCard{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.SimUseHistory{}) {

		err := glob.GDb.AutoMigrate(&models.SimUseHistory{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}

	if !glob.GDb.Migrator().HasTable(&models.MySQLTransmit{}) {

		err := glob.GDb.AutoMigrate(&models.MySQLTransmit{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.MySQLTransmitBind{}) {

		err := glob.GDb.AutoMigrate(&models.MySQLTransmitBind{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.MongoTransmit{}) {

		err := glob.GDb.AutoMigrate(&models.MongoTransmit{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.MongoTransmitBind{}) {

		err := glob.GDb.AutoMigrate(&models.MongoTransmitBind{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.InfluxdbTransmit{}) {

		err := glob.GDb.AutoMigrate(&models.InfluxdbTransmit{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.InfluxdbTransmitBind{}) {

		err := glob.GDb.AutoMigrate(&models.InfluxdbTransmitBind{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.ClickhouseTransmit{}) {

		err := glob.GDb.AutoMigrate(&models.ClickhouseTransmit{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.ClickhouseTransmitBind{}) {

		err := glob.GDb.AutoMigrate(&models.ClickhouseTransmitBind{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.DingDing{}) {

		err := glob.GDb.AutoMigrate(&models.DingDing{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.TcpHandler{}) {

		err := glob.GDb.AutoMigrate(&models.TcpHandler{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.DeviceBindTcpHandler{}) {

		err := glob.GDb.AutoMigrate(&models.DeviceBindTcpHandler{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
	if !glob.GDb.Migrator().HasTable(&models.HttpHandler{}) {

		err := glob.GDb.AutoMigrate(&models.HttpHandler{})
		if err != nil {
			zap.S().Errorf("数据库表创建失败 %+v", err)
		}
	}
}

func initDb() {

	username := glob.GConfig.MySQLConfig.Username //账号
	password := glob.GConfig.MySQLConfig.Password //密码
	host := glob.GConfig.MySQLConfig.Host         //数据库地址，可以是Ip或者域名
	port := glob.GConfig.MySQLConfig.Port         //数据库端口
	Dbname := glob.GConfig.MySQLConfig.DBName     //数据库名
	timeout := "10s"                              //连接超时，10秒

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		glob.GLog.Sugar().Errorf("数据库链接异常 ", err)

	}
	glob.GDb = db
}

func initMongo() {
	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%d", url.QueryEscape(glob.GConfig.MongoConfig.Username), url.QueryEscape(glob.GConfig.MongoConfig.Password), glob.GConfig.MongoConfig.Host, glob.GConfig.MongoConfig.Port)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	glob.GMongoClient = client

	fmt.Println("Connected to MongoDB!")

}

var myTimeEncoder = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 按照 "2006-01-02 15:04:05" 的格式编码时间
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
})

func initLog() {
	encoderConfig := zapcore.EncoderConfig{
		// 使用自定义的时间编码器
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码日志级别
		EncodeTime:     myTimeEncoder,                 // 使用自定义的时间编码器
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码调用者
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), // 使用 Console 编码器
		zapcore.AddSync(os.Stdout),          // 输出到标准输出
		zap.NewAtomicLevelAt(zap.InfoLevel), // 设置日志级别为 Debug
	)

	lg := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // 替换全局 Logger

	// 确保日志被刷新
	defer func(lg *zap.Logger) {
		err := lg.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) {
			zap.S().Errorf("日志同步失败 %+v", err)
		}
	}(lg)

	// 记录一条日志作为示例
	lg.Debug("这是一个调试级别的日志")
	glob.GLog = lg
}

func initRouter(r *gin.RouterGroup) {
	r.Use(router.JwtCheck())
	r.GET("/p/metrics", gin.WrapH(promhttp.Handler()))
	r.POST("/mqtt/create", mqttApi.CreateMqtt)
	r.GET("/mqtt/page", mqttApi.PageMqtt)
	r.GET("/mqtt/start", mqttApi.StartMqtt)
	r.GET("/mqtt/stop", mqttApi.StopMqtt)
	r.POST("/mqtt/update", mqttApi.UpdateMqtt)
	r.POST("/mqtt/delete/:id", mqttApi.DeleteMqtt)
	r.GET("/mqtt/node-using-status", mqttApi.NodeUsingStatus)
	r.POST("/mqtt/set-script", mqttApi.SetScript)
	r.POST("/mqtt/check-script", mqttApi.CheckScript)
	r.POST("/mqtt/send", mqttApi.SendMqttMessage)
	r.POST("/query/influxdb", influxdbApi.QueryInfluxdb)
	r.POST("/query/str-influxdb", influxdbApi.QueryInfluxdbString)

	r.POST("/signal/create", signalApi.CreateSignal)
	r.POST("/signal/update", signalApi.UpdateSignal)
	r.POST("/signal/delete/:id", signalApi.DeleteSignal)
	r.GET("/signal/page", signalApi.PageSignal)

	r.POST("/signal-waring-config/create", signalWaringConfigApi.CreateSignalWaringConfig)
	r.POST("/signal-waring-config/delete/:id", signalWaringConfigApi.DeleteSignalWaringConfig)
	r.POST("/signal-waring-config/query-row", signalWaringConfigApi.QueryWaringList)
	r.POST("/signal-waring-config/update", signalWaringConfigApi.UpdateSignalWaringConfig)
	r.GET("/signal-waring-config/page", signalWaringConfigApi.PageSignalWaringConfig)

	r.POST("/dashboard/create", dashboardApi.CreateDashboard)
	r.POST("/dashboard/update", dashboardApi.UpdateDashboard)
	r.GET("/dashboard/:id", dashboardApi.ByIdDashboard)
	r.GET("/dashboard/page", dashboardApi.PageDashboard)
	r.POST("/dashboard/delete/:id", dashboardApi.DeleteDashboard)

	r.POST("/MySQLTransmit/create", mysqlTransmitApi.CreateMySQLTransmit)
	r.POST("/MySQLTransmit/update", mysqlTransmitApi.UpdateMySQLTransmit)
	r.GET("/MySQLTransmit/:id", mysqlTransmitApi.ByIdMySQLTransmit)
	r.GET("/MySQLTransmit/page", mysqlTransmitApi.PageMySQLTransmit)
	r.POST("/MySQLTransmit/delete/:id", mysqlTransmitApi.DeleteMySQLTransmit)
	r.POST("/MySQLTransmit/mockScript", mysqlTransmitApi.MockScript)

	r.POST("/MongoTransmitApi/create", mongoTransmitApi.CreateMongoTransmit)
	r.POST("/MongoTransmitApi/update", mongoTransmitApi.UpdateMongoTransmit)
	r.GET("/MongoTransmitApi/:id", mongoTransmitApi.ByIdMongoTransmit)
	r.GET("/MongoTransmitApi/page", mongoTransmitApi.PageMongoTransmit)
	r.POST("/MongoTransmitApi/delete/:id", mongoTransmitApi.DeleteMongoTransmit)
	r.POST("/MongoTransmitApi/mockScript", mongoTransmitApi.MockScript)

	r.POST("/InfluxdbTransmitApi/create", influxdbTransmitApi.CreateInfluxdbTransmit)
	r.POST("/InfluxdbTransmitApi/update", influxdbTransmitApi.UpdateInfluxdbTransmit)
	r.GET("/InfluxdbTransmitApi/:id", influxdbTransmitApi.ByIdInfluxdbTransmit)
	r.GET("/InfluxdbTransmitApi/page", influxdbTransmitApi.PageInfluxdbTransmit)
	r.POST("/InfluxdbTransmitApi/delete/:id", influxdbTransmitApi.DeleteInfluxdbTransmit)
	r.POST("/InfluxdbTransmitApi/mockScript", influxdbTransmitApi.MockScript)

	r.POST("/ClickTransmitApi/create", clickTransmitApi.CreateClickhouseTransmit)
	r.POST("/ClickTransmitApi/update", clickTransmitApi.UpdateClickhouseTransmit)
	r.GET("/ClickTransmitApi/:id", clickTransmitApi.ByIdClickhouseTransmit)
	r.GET("/ClickTransmitApi/page", clickTransmitApi.PageClickhouseTransmit)
	r.POST("/ClickTransmitApi/delete/:id", clickTransmitApi.DeleteClickhouseTransmit)
	r.POST("/ClickTransmitApi/mockScript", clickTransmitApi.MockScript)

	r.POST("/CassandraTransmitApi/create", cassandraTransmitApi.CreateCassandraTransmit)
	r.POST("/CassandraTransmitApi/update", cassandraTransmitApi.UpdateCassandraTransmit)
	r.GET("/CassandraTransmitApi/:id", cassandraTransmitApi.ByIdCassandraTransmit)
	r.GET("/CassandraTransmitApi/page", cassandraTransmitApi.PageCassandraTransmit)
	r.POST("/CassandraTransmitApi/delete/:id", cassandraTransmitApi.DeleteCassandraTransmit)
	r.POST("/CassandraTransmitApi/mockScript", cassandraTransmitApi.MockScript)

	r.POST("/product/create", productApi.CreateProduct)
	r.POST("/product/update", productApi.UpdateProduct)
	r.GET("/product/:id", productApi.ByIdProduct)
	r.GET("/product/page", productApi.PageProduct)
	r.POST("/product/delete/:id", productApi.DeleteProduct)

	r.POST("/device_group/create", deviceGroupApi.CreateDeviceGroup)
	r.POST("/device_group/update", deviceGroupApi.UpdateDeviceGroup)
	r.GET("/device_group/:id", deviceGroupApi.ByIdDeviceGroup)
	r.GET("/device_group/page", deviceGroupApi.PageDeviceGroup)
	r.POST("/device_group/delete/:id", deviceGroupApi.DeleteDeviceGroup)

	r.GET("/device_group/query_bind_device", deviceGroupApi.QueryBindDeviceInfo)
	r.POST("/device_group/bind_device", deviceGroupApi.BindDeviceInfo)
	r.POST("/device_group/BindMqtt", deviceGroupApi.BindMqtt)
	r.POST("/device_group/QueryBindMqtt", deviceGroupApi.QueryBindMqtt)

	r.POST("/DeviceInfo/create", deviceInfoApi.CreateDeviceInfo)
	r.POST("/DeviceInfo/update", deviceInfoApi.UpdateDeviceInfo)
	r.GET("/DeviceInfo/:id", deviceInfoApi.ByIdDeviceInfo)
	r.GET("/DeviceInfo/page", deviceInfoApi.PageDeviceInfo)
	r.POST("/DeviceInfo/delete/:id", deviceInfoApi.DeleteDeviceInfo)
	r.POST("/DeviceInfo/BindMqtt", deviceInfoApi.BindMqtt)
	r.POST("/DeviceInfo/BindTcp", deviceInfoApi.BindTcp)
	r.POST("/DeviceInfo/BindHTTP", deviceInfoApi.BindHTTP)
	r.POST("/DeviceInfo/BindHCoap", deviceInfoApi.BindHCoap)
	r.GET("/DeviceInfo/QueryBindMqtt", deviceInfoApi.QueryBindMqtt)
	r.GET("/DeviceInfo/QueryBindTcp", deviceInfoApi.QueryBindTcp)
	r.GET("/DeviceInfo/QueryBindHTTP", deviceInfoApi.QueryBindHttp)
	r.GET("/DeviceInfo/QueryBindCoap", deviceInfoApi.QueryBindCoap)

	r.POST("/ProductionPlan/create", productionPlanApi.CreateProductionPlan)
	r.POST("/ProductionPlan/update", productionPlanApi.UpdateProductionPlan)
	r.GET("/ProductionPlan/:id", productionPlanApi.ByIdProductionPlan)
	r.GET("/ProductionPlan/page", productionPlanApi.PageProductionPlan)
	r.POST("/ProductionPlan/delete/:id", productionPlanApi.DeleteProductionPlan)

	r.POST("/SimCard/create", simCardApi.CreateSimCard)
	r.POST("/SimCard/update", simCardApi.UpdateSimCard)
	r.GET("/SimCard/page", simCardApi.PageSimCard)
	r.POST("/SimCard/delete/:id", simCardApi.DeleteSimCard)
	r.GET("/SimCard/:id", simCardApi.ByIdSimCard)
	r.POST("/SimCard/BindDeviceInfo", simCardApi.BindDeviceInfo)

	r.POST("/RepairRecord/create", repairRecordApi.CreateRepairRecord)
	r.POST("/RepairRecord/update", repairRecordApi.UpdateRepairRecord)
	r.GET("/RepairRecord/:id", repairRecordApi.ByIdRepairRecord)
	r.GET("/RepairRecord/page", repairRecordApi.PageRepairRecord)
	r.POST("/RepairRecord/delete/:id", repairRecordApi.DeleteRepairRecord)

	r.POST("/calc-rule/create", calcRuleApi.CreateCalcRule)
	r.POST("/calc-rule/update", calcRuleApi.UpdateCalcRule)
	r.GET("/calc-rule/page", calcRuleApi.PageCalcRule)
	r.POST("/calc-rule/delete/:id", calcRuleApi.DeleteCalcRule)
	r.POST("/calc-rule/start/:id", calcRuleApi.StartCalcRule)
	r.POST("/calc-rule/stop/:id", calcRuleApi.StopCalcRule)
	r.POST("/calc-rule/refresh/:id", calcRuleApi.Refresh)
	r.POST("/calc-rule/mock", calcRuleApi.MockCalcRule)
	r.GET("/calc-rule/rd", calcRuleApi.CalcRuleResult)

	r.POST("/calc-param/create", calcParamApi.CreateCalcParam)
	r.POST("/calc-param/update", calcParamApi.UpdateCalcParam)
	r.GET("/calc-param/page", calcParamApi.PageCalcParam)
	r.POST("/calc-param/delete/:id", calcParamApi.DeleteCalcParam)

	r.POST("/User/create", userApi.CreateUser)
	r.POST("/User/update", userApi.UpdateUser)
	r.GET("/User/page", userApi.PageUser)
	r.POST("/User/delete/:id", userApi.DeleteUser)
	r.GET("/User/:id", userApi.ByIdUser)
	r.GET("/User/list", userApi.ListUser)
	r.POST("/User/BindRole", userApi.BindRole)
	r.GET("/User/QueryBindRole", userApi.QueryBindRole)
	r.POST("/User/BindDeviceInfo", userApi.BindDeviceInfo)
	r.POST("/User/QueryBindDeviceInfo", userApi.QueryBindDeviceInfo)

	r.POST("/Dept/create", deptApi.CreateDept)
	r.POST("/Dept/update", deptApi.UpdateDept)
	r.GET("/Dept/page", deptApi.PageDept)
	r.POST("/Dept/delete/:id", deptApi.DeleteDept)
	r.GET("/Dept/:id", deptApi.ByIdDept)
	r.GET("/Dept/subs", deptApi.FindByIdSubs)

	r.POST("/Role/create", roleApi.CreateRole)
	r.POST("/Role/update", roleApi.UpdateRole)
	r.GET("/Role/page", roleApi.PageRole)
	r.POST("/Role/delete/:id", roleApi.DeleteRole)
	r.GET("/Role/:id", roleApi.ByIdRole)
	r.GET("/Role/list", roleApi.ListRole)

	r.POST("/ShipmentRecord/create", shipmentRecordApi.CreateShipmentRecord)
	r.POST("/ShipmentRecord/update", shipmentRecordApi.UpdateShipmentRecord)
	r.GET("/ShipmentRecord/page", shipmentRecordApi.PageShipmentRecord)
	r.POST("/ShipmentRecord/delete/:id", shipmentRecordApi.DeleteShipmentRecord)
	r.GET("/ShipmentRecord/:id", shipmentRecordApi.ByIdShipmentRecord)

	r.POST("/signal-delay-waring-param/create", signalDelayWaringParamApi.CreateSignalDelayWaring)
	r.POST("/signal-delay-waring-param/update", signalDelayWaringParamApi.UpdateSignalDelayWaring)
	r.GET("/signal-delay-waring-param/page", signalDelayWaringParamApi.PageSignalDelayWaring)
	r.POST("/signal-delay-waring-param/delete/:id", signalDelayWaringParamApi.DeleteSignalDelayWaring)

	r.POST("/signal-delay-waring/create", signalDelayWaringApi.CreateSignalDelayWaring)
	r.POST("/signal-delay-waring/update", signalDelayWaringApi.UpdateSignalDelayWaring)
	r.GET("/signal-delay-waring/page", signalDelayWaringApi.PageSignalDelayWaring)
	r.POST("/signal-delay-waring/delete/:id", signalDelayWaringApi.DeleteSignalDelayWaring)
	r.POST("/signal-delay-waring/Mock/:id", signalDelayWaringApi.Mock)
	r.POST("/signal-delay-waring/GenParam/:id", signalDelayWaringApi.GenParam)
	r.POST("/signal-delay-waring/query-row", signalDelayWaringApi.QueryWaringList)

	r.POST("/file/update", fileApi.UpdateFile)
	r.GET("/file/download", fileApi.DownloadFile)

	r.POST("/userinfo", loginApi.UserInfo)

	r.GET("/MessageList/page", messageListApi.PageMessageList)

	r.POST("/DingDing/create", dingdingApi.CreateDingDing)
	r.POST("/DingDing/update", dingdingApi.UpdateDingDing)
	r.GET("/DingDing/:id", dingdingApi.ByIdDingDing)
	r.GET("/DingDing/page", dingdingApi.PageDingDing)
	r.POST("/DingDing/delete/:id", dingdingApi.DeleteDingDing)
	r.POST("/DingDing/bind", dingdingApi.Bind)

	r.POST("/FeiShuId/create", feishuApi.CreateFeiShu)
	r.POST("/FeiShuId/update", feishuApi.UpdateFeiShu)
	r.GET("/FeiShuId/:id", feishuApi.ByIdFeiShu)
	r.GET("/FeiShuId/page", feishuApi.PageFeiShu)
	r.POST("/FeiShuId/delete/:id", feishuApi.DeleteFeiShu)
	r.POST("/FeiShuId/bind", feishuApi.Bind)

	r.POST("/TcpHandler/create", tcpHandlerApi.CreateTcpHandler)
	r.POST("/TcpHandler/update", tcpHandlerApi.UpdateTcpHandler)
	r.GET("/TcpHandler/:id", tcpHandlerApi.ByIdTcpHandler)
	r.GET("/TcpHandler/page", tcpHandlerApi.PageTcpHandler)
	r.POST("/TcpHandler/delete/:id", tcpHandlerApi.DeleteTcpHandler)

	r.POST("/HttpHandler/create", httpHandlerApi.CreateHttpHandler)
	r.POST("/HttpHandler/update", httpHandlerApi.UpdateHttpHandler)
	r.GET("/HttpHandler/:id", httpHandlerApi.ByIdHttpHandler)
	r.GET("/HttpHandler/page", httpHandlerApi.PageHttpHandler)
	r.POST("/HttpHandler/delete/:id", httpHandlerApi.DeleteHttpHandler)

}
func initGlobalRedisClient() {

	add := fmt.Sprintf("%s:%d", glob.GConfig.RedisConfig.Host, glob.GConfig.RedisConfig.Port)
	glob.GRedis = redis.NewClient(&redis.Options{
		Addr:     add,
		Password: glob.GConfig.RedisConfig.Password, // 如果没有设置密码，就留空字符串
		DB:       glob.GConfig.RedisConfig.Db,       // 使用默认数据库
	})

	// 检查连接是否成功
	if err := glob.GRedis.Ping(context.Background()).Err(); err != nil {
		glob.GLog.Sugar().Fatalf("Could not connect to Redis: %v", err)
	}

}
func InitConfig() {
	var configPath string
	flag.StringVar(&configPath, "config", "app-node1.yml", "Path to the config file")
	flag.Parse()

	yfile, err := os.ReadFile(configPath)
	if err != nil {
		zap.S().Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(yfile, &glob.GConfig)
	if err != nil {
		zap.S().Fatalf("error: %v", err)
	}

}

func initTableData() {

	glob.GDb.Model(models.User{}).FirstOrInit(&models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "admin",
		Password: "admin",
		Email:    "",
		Status:   "active",
	})

	glob.GDb.Model(models.User{}).FirstOrInit(&models.User{
		Model: gorm.Model{
			ID: 2,
		},
		Username: "p1",
		Password: "p1",
		Email:    "",
		Status:   "active",
	})

	glob.GDb.Model(models.User{}).FirstOrInit(&models.User{
		Model: gorm.Model{
			ID: 3,
		},
		Username: "p1-1",
		Password: "p1-1",
		Email:    "",
		Status:   "active",
	})

	glob.GDb.Model(models.User{}).FirstOrInit(&models.User{
		Model: gorm.Model{
			ID: 4,
		},
		Username: "p2",
		Password: "p2",
		Email:    "",
		Status:   "active",
	})
	glob.GDb.Model(models.User{}).FirstOrInit(&models.User{
		Model: gorm.Model{
			ID: 5,
		},
		Username: "p2-1",
		Password: "p2-1",
		Email:    "",
		Status:   "active",
	})

	glob.GDb.Model(models.UserRole{}).FirstOrInit(&models.UserRole{
		Model: gorm.Model{
			ID: 1,
		},
		UserId: 1,
		RoleId: 1,
	})

	glob.GDb.Model(models.UserRole{}).FirstOrInit(&models.UserRole{
		Model: gorm.Model{
			ID: 2,
		},
		UserId: 2,
		RoleId: 2,
	})

	glob.GDb.Model(models.UserRole{}).FirstOrInit(&models.UserRole{
		Model: gorm.Model{
			ID: 3,
		},
		UserId: 3,
		RoleId: 3,
	})
	glob.GDb.Model(models.UserRole{}).FirstOrInit(&models.UserRole{
		Model: gorm.Model{
			ID: 4,
		},
		UserId: 4,
		RoleId: 4,
	})
	glob.GDb.Model(models.UserRole{}).FirstOrInit(&models.UserRole{
		Model: gorm.Model{
			ID: 5,
		},
		UserId: 5,
		RoleId: 5,
	})

	glob.GDb.Model(models.Role{}).FirstOrInit(&models.Role{
		Model: gorm.Model{
			ID: 1,
		},
		Name:        "超级管理员",
		Description: "超级管理员",
		CanDel:      false,
	})
	glob.GDb.Model(models.Role{}).FirstOrInit(&models.Role{
		Model: gorm.Model{
			ID: 2,
		},
		Name:        "生产管理员",
		Description: "负责完成生产计划的定制",
		CanDel:      false,
	})
	glob.GDb.Model(models.Role{}).FirstOrInit(&models.Role{
		Model: gorm.Model{
			ID: 3,
		},
		Name:        "生产人员",
		Description: "负责完成生产",
		CanDel:      false,
	})
	glob.GDb.Model(models.Role{}).FirstOrInit(&models.Role{
		Model: gorm.Model{
			ID: 4,
		},
		Name:        "维修员",
		Description: "负责现场维修",
		CanDel:      false,
	})
	glob.GDb.Model(models.Role{}).FirstOrInit(&models.Role{
		Model: gorm.Model{
			ID: 5,
		},
		Name:        "售后员",
		Description: "负责处理售后问题",
		CanDel:      false,
	})

	// 生产管理员 消息类型
	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 1,
		},
		MessageType: int(glob.StartNotification),
		RoleId:      2,
	})

	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 2,
		},
		MessageType: int(glob.DueSoonNotification),
		RoleId:      2,
	})

	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 3,
		},
		MessageType: int(glob.DueNotification),
		RoleId:      2,
	})

	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 4,
		},
		MessageType: int(glob.ProductionStartNotification),
		RoleId:      2,
	})

	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 5,
		},
		MessageType: int(glob.ProductionCompleteNotification),
		RoleId:      2,
	})

	// 生产人员 消息类型
	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 6,
		},
		MessageType: int(glob.StartNotification),
		RoleId:      3,
	})

	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 7,
		},
		MessageType: int(glob.DueSoonNotification),
		RoleId:      3,
	})

	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 8,
		},
		MessageType: int(glob.DueNotification),
		RoleId:      3,
	})

	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 9,
		},
		MessageType: int(glob.ProductionStartNotification),
		RoleId:      3,
	})

	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 10,
		},
		MessageType: int(glob.ProductionCompleteNotification),
		RoleId:      3,
	})

	//维修员 消息类型
	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 11,
		},
		MessageType: int(glob.MaintenanceNotification),
		RoleId:      4,
	})
	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 12,
		},
		MessageType: int(glob.MaintenanceStartNotification),
		RoleId:      4,
	})
	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 13,
		},
		MessageType: int(glob.MaintenanceEndNotification),
		RoleId:      4,
	})
	//售后员 消息类型

	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 14,
		},
		MessageType: int(glob.MaintenanceNotification),
		RoleId:      5,
	})
	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 15,
		},
		MessageType: int(glob.MaintenanceStartNotification),
		RoleId:      5,
	})
	glob.GDb.Model(&models.MessageTypeBindRole{}).FirstOrInit(&models.MessageTypeBindRole{
		Model: gorm.Model{
			ID: 16,
		},
		MessageType: int(glob.MaintenanceEndNotification),
		RoleId:      5,
	})

}
func InitAll(r *gin.RouterGroup) {
	InitConfig()
	createFileUpload()

	initLog()
	glob.GLog.Info("日志初始化完成")
	initDb()
	glob.GLog.Info("数据库已链接")
	initTable()
	initTableData()
	glob.GLog.Info("数据库表已生成")
	initGlobalRedisClient()
	glob.GLog.Info("redis 客户端连接成功")
	InitRabbitCon()
	initMongo()

	initRouter(r)
	//go biz.InitRedisExpireHandler(glob.GRedis)
	InitInfluxDbClient()
}

func InitInfluxDbClient() {
	glob.GInfluxdb = influxdb2.NewClient(fmt.Sprintf("http://%s:%d", glob.GConfig.InfluxConfig.Host, glob.GConfig.InfluxConfig.Port), glob.GConfig.InfluxConfig.Token)

}

func InitRabbitCon() {
	conn, err := amqp.Dial(genUrl())
	if err != nil {
		zap.S().Fatalf("Failed to connect to RabbitMQ  %v", err)
	}

	glob.GRabbitMq = conn

	CreateRabbitQueue("calc_queue")

}

func genUrl() string {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", glob.GConfig.MQConfig.Username, glob.GConfig.MQConfig.Password, glob.GConfig.MQConfig.Host, glob.GConfig.MQConfig.Port)
	return connStr
}

func CreateRabbitQueue(queueName string) {

	ch, err := glob.GRabbitMq.Channel()
	if err != nil {
		zap.S().Fatalf("Failed to open a channel %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		zap.S().Fatalf("创建queue异常 %s", queueName)
	}
}

func createFileUpload() {
	// 设置文件夹路径
	dirPath := "./fileupdate"

	// 检查文件夹是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err := os.MkdirAll(dirPath, 0755) // 0755是文件夹权限设置，可根据需要调整
		if err != nil {
			fmt.Println("创建文件夹失败：", err)
			return
		}
		fmt.Println("文件夹已创建：", dirPath)
	} else {
		// 文件夹已存在
		fmt.Println("文件夹已存在：", dirPath)
	}
}
