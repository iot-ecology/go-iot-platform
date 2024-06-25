package models

import (
	"gorm.io/gorm"
	"time"
)

type MqttClient struct {
	Host     string `json:"host"`      // 主机
	Port     int    `json:"port"`      // 端口
	ClientId string `json:"client_id"` // 客户端id
	Username string `json:"username"`  // 账号
	Password string `json:"password"`  // 密码
	Subtopic string `json:"subtopic"`  // 订阅的主题
	Start    bool   `json:"start"`     // 是否启动

	Script     string `json:"script" gorm:"type:text"` // 数据处理脚本
	gorm.Model `structs:"-"`
}

type Signal struct {
	MqttClientId   int    `json:"mqtt_client_id"`                  // MQTT客户端表的外键ID
	Name           string `json:"name"`                            // 信号的名称，用于标识不同的信号
	Alias          string `json:"alias" structs:"alias"`           // 信号的别名，用于显示
	Type           string `json:"type"`                            // 信号的数据类型，如整数、字符串等
	MqttClientName string `gorm:"-" json:"mqtt_client_name"`       // MQTT客户端的名称，不存储在数据库中
	Unit           string `json:"unit" structs:"unit"`             // 单位
	CacheSize      int    `json:"cache_size" structs:"cache_size"` // 缓存大小
	gorm.Model     `structs:"-"`
}
type SignalWaringConfig struct {
	SignalId     int     `gorm:"signal_id"  json:"signal_id" structs:"signal_id"` // 信号表的外键ID
	Min          float64 `gorm:"min"  json:"min" structs:"min"`                   // 范围,小值
	Max          float64 `gorm:"max"  json:"max" structs:"max"`                   // 范围,大值
	InOrOut      int     `gorm:"in_or_out"  json:"in_or_out" structs:"in_or_out"` //范围内报警,范围外报警 1 范围内报警 0 范围外报警
	MqttClientId int     `json:"mqtt_client_id" structs:"mqtt_client_id"`         // MQTT客户端表的外键ID
	gorm.Model   `structs:"-"`
	Signal       Signal `gorm:"-" json:"signal" structs:"-"`
}

// SignalDelayWaring : 允许信号数据做一定条数的暂存
type SignalDelayWaring struct {
	Name       string `json:"name" structs:"name"`     // 名称
	Script     string `json:"script" gorm:"type:text"` // 数据处理脚本
	gorm.Model `structs:"-"`
}

type SignalDelayWaringParam struct {
	MqttClientName      string `gorm:"-" json:"mqtt_client_name"`                             // MQTT客户端的名称，不存储在数据库中
	MqttClientId        int    `json:"mqtt_client_id"`                                        // MQTT客户端表的外键ID
	Name                string `json:"name"`                                                  // 参数名称
	SignalName          string `gorm:"signal_name"  json:"signal_name" structs:"signal_name"` // 信号表 name
	SignalId            int    `gorm:"signal_id"  json:"signal_id" structs:"signal_id"`       // 信号表的外键ID
	SignalDelayWaringId int    `json:"signal_delay_waring_id"`                                // SignalDelayWaring 主键
	gorm.Model          `structs:"-"`
}

type NodeInfo struct {
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	Port int    `json:"port,omitempty" yaml:"port,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	Size int64  `json:"size,omitempty" yaml:"size,omitempty"` // 最大处理数量
}

// CalcRule 计算规则
type CalcRule struct {
	Name       string `json:"name" structs:"name"`                              // 任务名称
	Cron       string `json:"cron"`                                             // 定时任务执行周期
	Script     string `json:"script" gorm:"type:text"`                          // 数据处理脚本
	Offset     int64  `json:"offset" structs:"offset"`                          // 执行时间往前推移几秒
	Start      bool   `json:"start" structs:"start"`                            // 是否启动
	MockValue  string `json:"mock_value" structs:"mock_value" gorm:"type:text"` // 模拟数据
	gorm.Model `structs:"-"`
}

// CalcParam 计算参数
type CalcParam struct {
	MqttClientId   int    `json:"mqtt_client_id"`                                        // MQTT客户端表的外键ID
	Name           string `json:"name"`                                                  // 参数名称
	SignalName     string `gorm:"signal_name"  json:"signal_name" structs:"signal_name"` // 信号表 name
	SignalId       int    `gorm:"signal_id"  json:"signal_id" structs:"signal_id"`       // 信号表的外键ID
	Reduce         string `json:"reduce"`                                                // 数据聚合方式 1. mean 2. sum 3. max 4. min
	CalcRuleId     int    `json:"calc_rule_id"`                                          // CalcRule 主键
	MqttClientName string `gorm:"-" json:"mqtt_client_name"`                             // MQTT客户端的名称，不存储在数据库中
	gorm.Model     `structs:"-"`
}

// Dashboard 可视化板
type Dashboard struct {
	Name       string `json:"name"`                    // 名称
	Config     string `json:"config" gorm:"type:text"` // 配置
	gorm.Model `structs:"-"`
}

type Product struct {
	Name           string  `json:"name"`                                      // 产品名称
	Description    string  `json:"description" structs:"description"`         // 产品描述
	SKU            string  `json:"sku" structs:"sku"`                         // 库存单位
	Price          float64 `json:"price" structs:"price"`                     // 销售价格
	Cost           float64 `json:"cost" structs:"cost"`                       // 成本
	Quantity       int     `json:"quantity" structs:"quantity"`               // 库存数量
	MinimumStock   int     `json:"minimum_stock" structs:"minimum_stock"`     // 最低库存量
	WarrantyPeriod int     `json:"warranty_period" structs:"warranty_period"` // 质保时间（天）
	Status         string  `json:"status" structs:"status"`                   // 产品状态
	Tags           string  `json:"tags" structs:"tags"`                       // 标签
	ImageURL       string  `json:"image_url" structs:"image_url"`             // 图片URL
	gorm.Model     `structs:"-"`
}
type DeviceInfo struct {
	ProductId         uint      `json:"product_id" structs:"product_id"`                   // 产品ID
	ProductionBatchId uint      `json:"production_batch_id" structs:"production_batch_id"` // 产品批次ID
	SN                string    `json:"sn" structs:"sn"`                                   // 设备编号
	ManufacturingDate time.Time `json:"manufacturing_date" structs:"manufacturing_date"`   // 制造日期
	ProcurementDate   time.Time `json:"procurement_date" structs:"procurement_date"`       // 采购日期
	Source            int       `json:"source" structs:"source"`                           // 设备来源,1: 内部,2: 外源
	WarrantyExpiry    time.Time `json:"warranty_expiry" structs:"warranty_expiry"`         // 保修截止日期
	gorm.Model        `structs:"-"`
}

type DeviceGroup struct {
	Name       string `json:"name" structs:"name"` // 名称
	gorm.Model `structs:"-"`
}

type DeviceGroupDevice struct {
	DeviceInfoId uint `json:"device_info_id" structs:"device_info_id"` // 设备表的外键ID
	GroupId      uint `json:"group_id" structs:"group_id"`             // 设备组表的外键ID
	gorm.Model   `structs:"-"`
}

type DeviceInstallRecord struct {
	gorm.Model   `structs:"-"`
	DeviceInfoId uint      `json:"device_info_id" structs:"device_info_id"` // 设备表的外键ID
	InstallDate  time.Time `json:"install_date" structs:"install_date"`     // 安装日期
	Technician   string    `json:"technician" structs:"technician"`         // 安装人员
	Description  string    `json:"description" structs:"description"`       // 安装描述
	PhotoURL     string    `json:"photo_url" structs:"photo_url"`           // 照片URL
}

type RepairRecord struct {
	DeviceInfoId uint      `json:"device_info_id" structs:"device_info_id"` // 设备表的外键ID
	RepairDate   time.Time `json:"repair_date" structs:"repair_date"`       // 维修日期
	Technician   string    `json:"technician" structs:"technician"`         // 维修人员
	Cost         float64   `json:"cost" structs:"cost"`                     // 维修成本
	Description  string    `json:"description" structs:"description"`       // 维修描述
	gorm.Model   `structs:"-"`
}

// ShipmentRecord 发货记录
type ShipmentRecord struct {
	gorm.Model      `structs:"-"`
	ShipmentDate    time.Time `json:"shipment_date" structs:"shipment_date"`       // 发货日期
	Technician      string    `json:"technician" structs:"technician"`             // 发货人员
	CustomerName    string    `json:"customer_name" structs:"customer_name"`       // 客户名称
	CustomerPhone   string    `json:"customer_phone" structs:"customer_phone"`     // 客户手机
	CustomerAddress string    `json:"customer_address" structs:"customer_address"` // 客户地址
	TrackingNumber  string    `json:"tracking_number" structs:"tracking_number"`   // 跟踪号码
	Status          string    `json:"status" structs:"status"`                     // 发货状态（例如：pending, shipped, delivered）
	Description     string    `json:"description" structs:"description"`           // 发货描述
}

// ShipmentProductDetail 发货记录中的具体产品
type ShipmentProductDetail struct {
	gorm.Model       `structs:"-"`
	ShipmentRecordId uint `json:"shipment_record_id" structs:"shipment_record_id"` // 发货记录ID
	ProductID        uint `json:"product_id" structs:"product_id"`                 // 关联的产品ID
	Quantity         int  `json:"quantity" structs:"quantity"`                     // 发货数量
}

// ProductionPlan 表示生产计划
type ProductionPlan struct {
	gorm.Model  `structs:"-"`
	Name        string    `json:"name" structs:"name"`               // 生产计划名称
	StartDate   time.Time `json:"start_date" structs:"start_date"`   // 生产计划开始日期
	EndDate     time.Time `json:"end_date" structs:"end_date"`       // 生产计划结束日期
	Description string    `json:"description" structs:"description"` // 生产计划描述
	Status      string    `json:"status" structs:"status"`           // 计划状态（准备中,进行中, 已完成）
}

// ProductPlan 表示生产计划中的具体产品计划
type ProductPlan struct {
	gorm.Model       `structs:"-"`
	ProductionPlanID uint `json:"production_plan_id" structs:"production_plan_id"` // 关联的生产计划ID
	ProductID        uint `json:"product_id" structs:"product_id"`                 // 关联的产品ID
	Quantity         int  `json:"quantity" structs:"quantity"`                     // 计划生产数量
}

type User struct {
	gorm.Model `structs:"-"`
	Username   string `json:"username" structs:"username"` // 用户名
	Password   string `json:"password" structs:"password"` // 密码
	Email      string `json:"email" structs:"email"`       // 电子邮箱
	Status     string `json:"status" structs:"status"`     // 用户状态（例如：active, inactive）
}

type Role struct {
	gorm.Model  `structs:"-"`
	Name        string `json:"name" structs:"name"`               // 角色名
	Description string `json:"description" structs:"description"` // 角色描述
}

type UserRole struct {
	gorm.Model `structs:"-"`
	UserId     uint `json:"user_id" structs:"user_id"` // 用户ID
	RoleId     uint `json:"role_id" structs:"role_id"` // 角色ID
}
type Dept struct {
	gorm.Model `structs:"-"`
	Name       string `json:"name" structs:"name"`                     // 部门名
	ParentId   uint   `json:"parent_id,omitempty" structs:"parent_id"` // 父部门ID
}

type UserBindDeviceInfo struct {
	gorm.Model `structs:"-"`
	UserId     uint `json:"user_id" structs:"user_id"`     // 用户ID
	DeviceId   uint `json:"device_id" structs:"device_id"` // 设备ID
}

type DeviceBindMqttClient struct {
	gorm.Model   `structs:"-"`
	DeviceInfoId uint `json:"device_info_id" structs:"device_info_id"` // 设备ID
	MqttClientId uint `json:"mqtt_client_id"`                          // MQTT客户端表的外键ID
}

type DeviceGroupBindMqttClient struct {
	gorm.Model    `structs:"-"`
	DeviceGroupId uint `json:"device_group_id" structs:"device_group_id"` // 设备组ID
	MqttClientId  uint `json:"mqtt_client_id"`                            // MQTT客户端表的外键ID
}

type MessageTypeBindRole struct {
	gorm.Model  `structs:"-"`
	MessageType string `json:"message_type" structs:"message_type"` // 消息类型
	RoleId      uint   `json:"role_id" structs:"role_id"`           // 角色ID
}

type MessageList struct {
	gorm.Model    `structs:"-"`
	Content       string `json:"content" structs:"content"`                 // 消息内容
	MessageTypeId uint   `json:"message_type_id" structs:"message_type_id"` // 消息类型ID
}
