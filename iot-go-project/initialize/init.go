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
	"log"
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
	if !glob.GDb.Migrator().HasTable(&models.ProductionBatch{}) {

		err := glob.GDb.AutoMigrate(&models.ProductionBatch{})
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
	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%d", glob.GConfig.MongoConfig.Username, glob.GConfig.MongoConfig.Password, glob.GConfig.MongoConfig.Host, glob.GConfig.MongoConfig.Port)
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
	r.POST("/DeviceInfo/QueryBindMqtt", deviceInfoApi.QueryBindMqtt)

	r.POST("/ProductionPlan/create", productionPlanApi.CreateProductionPlan)
	r.POST("/ProductionPlan/update", productionPlanApi.UpdateProductionPlan)
	r.GET("/ProductionPlan/:id", productionPlanApi.ByIdProductionPlan)
	r.GET("/ProductionPlan/page", productionPlanApi.PageProductionPlan)
	r.POST("/ProductionPlan/delete/:id", productionPlanApi.DeleteProductionPlan)

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
	r.POST("/User/QueryBindRole", userApi.QueryBindRole)
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
	r.POST("/file/download", fileApi.DownloadFile)

	r.POST("/login", loginApi.Login)
	r.POST("/userinfo", loginApi.UserInfo)

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

func InitAll(r *gin.RouterGroup) {
	InitConfig()
	initLog()
	glob.GLog.Info("日志初始化完成")
	initDb()
	glob.GLog.Info("数据库已链接")
	initTable()
	glob.GLog.Info("数据库表已生成")
	initGlobalRedisClient()
	glob.GLog.Info("redis 客户端连接成功")
	InitRabbitCon()
	initMongo()

	initRouter(r)
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
