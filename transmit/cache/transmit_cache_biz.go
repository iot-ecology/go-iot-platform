package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"iot-transmit/cassandra"
	"iot-transmit/clickhouse"
	"iot-transmit/common"
	"iot-transmit/influxdb2"
	"iot-transmit/mongo"
	"iot-transmit/mysql"
)

type TransmitCacheBiz struct {
}

var (
	clickhouseOp = clickhouse.ClickhouseOp{}
	mysqlOp      = mysql.MysqlOp{}
	mongoOp      = mongo.MongoOp{}
	cassandraOp  = cassandra.CassandraOp{}
	influxdbOp   = influxdb2.InfluxDbOp{}
)

// Run 是 TransmitCacheBiz 结构体的一个方法，用于并发执行多个数据库的处理操作
// 参数：
// redis: Redis 客户端实例，用于查询数据库缓存配置信息
// dataRowList: 需要处理的数据行列表
// 返回值：
//
//	无返回值
func (biz *TransmitCacheBiz) Run(redis *redis.Client, dataRowList []common.DataRowList,Protocol string ) {
	go biz.workMysql(redis, dataRowList ,Protocol)
	go biz.workMongo(redis, dataRowList ,Protocol)
	go biz.workCassandra(redis, dataRowList ,Protocol)
	go biz.workClickhouse(redis, dataRowList ,Protocol)
	go biz.workInfluxdb(redis, dataRowList ,Protocol)
}

// workInfluxdb 是TransmitCacheBiz结构体的一个方法，用于处理InfluxDB缓存的相关操作
// 参数：
// redis: Redis客户端实例，用于查找InfluxDB的缓存配置信息
// mqttClientId: MQTT客户端ID，作为缓存查询的标识
// dataRowList: 需要处理的数据行列表
// 返回值：
//
//	无返回值
func (biz *TransmitCacheBiz) workInfluxdb(redis *redis.Client, dataRowList []common.DataRowList, protocol string) {
	for _, row := range dataRowList {

		var influxdbCache = biz.findInfluxdb(redis, row.DeviceUid, row.IdentificationCode,protocol)
		for _, cache := range influxdbCache {
			db := influxdb2.GetInfluxDb(cache.Host, cache.Token, cache.Port, cache.ID)
			err := influxdbOp.HandleDataRowLists(cache.Bucket, cache.Org, cache.Measurement, cache.Script, dataRowList, db)
			if err != nil {
				zap.S().Errorf("influxdb 执行异常", zap.Error(err))
			}
		}
	}

}

// workClickhouse 是TransmitCacheBiz结构体的一个方法，用于处理ClickHouse缓存的相关操作
// 参数：
// redis: Redis客户端实例
// mqttClientId: MQTT客户端的ID
// dataRowList: 需要处理的数据行列表
// 返回值：
//
//	无返回值
func (biz *TransmitCacheBiz) workClickhouse(redis *redis.Client, dataRowList []common.DataRowList, protocol string) {
	for _, row := range dataRowList {
		var clickhouseCache = biz.findClickhouse(redis, row.DeviceUid, row.IdentificationCode,protocol)
		for _, cache := range clickhouseCache {
			println(cache.Script)
			house1, _ := clickhouse.GetClickHouse(row.DeviceUid, []string{fmt.Sprintf("%s:%d", cache.Host, cache.Port)}, cache.Database, cache.Username, cache.Password)

			err := clickhouseOp.HandleDataRowLists(cache.Table, cache.Script, dataRowList, house1)
			if err != nil {
				zap.S().Errorf("clickhouse 执行异常", zap.Error(err))
			}
		}
	}
}

// workCassandra 是TransmitCacheBiz结构体的一个方法，用于处理Cassandra缓存的相关操作
// 参数：
// redis: Redis客户端实例
// mqttClientId: MQTT客户端的ID
// dataRowList: 需要处理的数据行列表
// 返回值：
//
//	无返回值
func (biz *TransmitCacheBiz) workCassandra(redis *redis.Client, dataRowList []common.DataRowList, protocol string) {
	for _, row := range dataRowList {
		var cassandraCache = biz.findCassandra(redis, row.DeviceUid, row.IdentificationCode,protocol)
		for _, cache := range cassandraCache {
			getCassandra, err := cassandra.GetCassandra([]string{fmt.Sprintf("%s:%d", cache.Host, cache.Port)}, cache.Username, cache.Password, cache.ID)
			if err != nil {
				zap.S().Errorf("cassandra 连接异常", zap.Error(err))
			}
			err = cassandraOp.HandleDataRowLists(cache.Database, cache.Table, cache.Script, dataRowList, getCassandra)
			if err != nil {
				zap.S().Errorf("cassandra 执行异常", zap.Error(err))
			}
		}
	}

}

// workMongo 是 TransmitCacheBiz 结构体的一个方法，用于处理 MongoDB 缓存的相关操作
// 参数：
// redis: Redis 客户端实例
// mqttClientId: MQTT 客户端的 ID
// dataRowList: 需要处理的数据行列表
// 返回值：
//
//	无返回值
func (biz *TransmitCacheBiz) workMongo(redis *redis.Client, dataRowList []common.DataRowList, protocol string) {
	for _, row := range dataRowList {
		var mongoCache = biz.findMongo(redis, row.DeviceUid, row.IdentificationCode,protocol)
		for _, cache := range mongoCache {
			client, err := mongo.GetMongoDBClient(cache.Host, cache.Username, cache.Password, cache.Database, cache.Port, cache.ID)
			if err != nil {
				zap.S().Errorf("mongo 连接异常", zap.Error(err))
			}
			err = mongoOp.HandleDataRowLists(cache.Database, cache.Collection, cache.Script, dataRowList, client)
			if err != nil {
				zap.S().Errorf("mongo 执行异常", zap.Error(err))
			}
		}
	}
}

// workMysql 是一个TransmitCacheBiz的方法，用于处理与MySQL相关的业务逻辑
//
// 参数：
//
//	biz: TransmitCacheBiz实例指针，表示当前操作的业务对象
//	redis: *redis.Client，Redis客户端实例指针，用于从Redis中获取MySQL缓存信息
//	mqttClientId: string，MQTT客户端ID，用于从Redis中定位特定的MySQL缓存信息
//	dataRowList: []common.DataRowList，待处理的数据行列表
//
// 返回值：
//
//	无返回值
func (biz *TransmitCacheBiz) workMysql(redis *redis.Client, dataRowList []common.DataRowList, protocol string) {

	for _, row := range dataRowList {
		var mysqlCache = biz.findMysql(redis, row.DeviceUid, row.IdentificationCode,protocol)
		for _, cache := range mysqlCache {
			connection, err := mysql.InitMySQLConnection(cache.Username, cache.Host, cache.Password, cache.Database, cache.Port, cache.ID)
			if err != nil {
				zap.S().Errorf("mysql 连接异常", zap.Error(err))
			}
			err = mysqlOp.HandleDataRowLists(cache.Table, cache.Script, dataRowList, connection)
			if err != nil {
				zap.S().Errorf("mysql 执行异常", zap.Error(err))
			}
		}
	}

}

// findMysql 从Redis中获取与mqttClientId相关的MySQLTransmitCache列表
//
// 参数：
// biz: TransmitCacheBiz 实例指针
// redis: Redis客户端实例指针
// mqttClientId: MQTT客户端ID
//
// 返回值：
// []MySQLTransmitCache: MySQLTransmitCache的切片
func (biz *TransmitCacheBiz) findMysql(redis *redis.Client, DeviceUid, IdentificationCode, protocol string) []MySQLTransmitCache {
	val := redis.LRange(context.Background(), "transmit:mysql:"+protocol+":"+DeviceUid+":"+IdentificationCode, 0,
		-1).Val()
	var c []MySQLTransmitCache

	for _, s := range val {
		var e MySQLTransmitCache
		err := json.Unmarshal([]byte(s), &e)
		if err != nil {
			continue
		}
		c = append(c, e)

	}
	return c
}

// findMongo 从Redis中查找与mqttClientId相关的MongoTransmitCache列表
//
// 参数：
//
//	biz *TransmitCacheBiz：TransmitCacheBiz实例指针
//	redis *redis.Client：Redis客户端实例指针
//	mqttClientId string：MQTT客户端ID
//
// 返回值：
//
//	[]MongoTransmitCache：MongoTransmitCache列表
func (biz *TransmitCacheBiz) findMongo(redis *redis.Client, DeviceUid, IdentificationCode, protocol string) []MongoTransmitCache {
	val := redis.LRange(context.Background(), "transmit:mongo:"+protocol+":"+DeviceUid+":"+IdentificationCode, 0, -1).Val()
	var c []MongoTransmitCache

	for _, s := range val {
		var e MongoTransmitCache
		err := json.Unmarshal([]byte(s), &e)
		if err != nil {
			continue
		}
		c = append(c, e)

	}
	return c
}

// findCassandra 从Redis中根据mqttClientId获取CassandraTransmitCache的切片
// 参数：
// biz: TransmitCacheBiz类型指针，表示业务逻辑处理对象
// redis: *redis.Client类型指针，表示Redis客户端对象
// mqttClientId: string类型，表示MQTT客户端ID
//
// 返回值:
// []CassandraTransmitCache类型切片，包含从Redis中获取的CassandraTransmitCache对象
func (biz *TransmitCacheBiz) findCassandra(redis *redis.Client, DeviceUid, IdentificationCode, protocol string) []CassandraTransmitCache {
	val := redis.LRange(context.Background(), "transmit:cassandra:"+protocol+":"+DeviceUid+":"+IdentificationCode, 0,
		-1).Val()
	var c []CassandraTransmitCache

	for _, s := range val {
		var e CassandraTransmitCache
		err := json.Unmarshal([]byte(s), &e)
		if err != nil {
			continue
		}
		c = append(c, e)

	}
	return c
}

// findClickhouse 从Redis中根据mqttClientId获取对应的ClickhouseTransmitCache切片
//
// 参数：
// biz: *TransmitCacheBiz - TransmitCacheBiz类型的指针，业务逻辑处理对象
// redis: *redis.Client - redis.Client类型的指针，Redis客户端对象
// mqttClientId: string - MQTT客户端ID
//
// 返回值：
// []ClickhouseTransmitCache - ClickhouseTransmitCache类型的切片，包含从Redis中获取的ClickhouseTransmitCache对象
func (biz *TransmitCacheBiz) findClickhouse(redis *redis.Client, DeviceUid, IdentificationCode, protocol string) []ClickhouseTransmitCache {
	val := redis.LRange(context.Background(), "transmit:clickhouse:" +protocol+":"+DeviceUid+":"+IdentificationCode, 0,
		-1).Val()
	var c []ClickhouseTransmitCache

	for _, s := range val {
		var e ClickhouseTransmitCache
		err := json.Unmarshal([]byte(s), &e)
		if err != nil {
			continue
		}
		c = append(c, e)

	}
	return c
}

// findInfluxdb 从Redis中根据mqttClientId获取对应的InfluxTransmitCache切片
//
// 参数：
//
//	biz *TransmitCacheBiz：TransmitCacheBiz类型的指针，业务逻辑处理对象
//	redis *redis.Client：redis.Client类型的指针，Redis客户端对象
//	mqttClientId string：string类型，MQTT客户端ID
//
// 返回值：
//
//	[]InfluxTransmitCache：InfluxTransmitCache类型的切片，包含从Redis中获取的InfluxTransmitCache对象
func (biz *TransmitCacheBiz) findInfluxdb(redis *redis.Client, DeviceUid, IdentificationCode, protocol string) []InfluxTransmitCache {
	val := redis.LRange(context.Background(), "transmit:influxdb:"+protocol+":"+DeviceUid+":"+IdentificationCode, 0,
		-1).Val()
	var c []InfluxTransmitCache

	for _, s := range val {
		var e InfluxTransmitCache
		err := json.Unmarshal([]byte(s), &e)
		if err != nil {
			continue
		}
		c = append(c, e)

	}
	return c
}
