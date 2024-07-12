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

func (biz *TransmitCacheBiz) Run(redis *redis.Client, mqttClientId string, dataRowList []common.DataRowList) {
	var mysqlCache = biz.findMysql(redis, mqttClientId)
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

	var mongoCache = biz.findMongo(redis, mqttClientId)
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

	var cassandraCache = biz.findCassandra(redis, mqttClientId)
	for _, cache := range cassandraCache {
		getCassandra, err := cassandra.GetCassandra([]string{fmt.Sprintf("%s:%d", cache.Host, cache.Port)}, cache.Username, cache.Password, cache.ID)
		if err != nil {
			zap.S().Errorf("cassandra 连接异常", zap.Error(err))
		}
		cassandraOp.HandleDataRowLists(cache.Database, cache.Table, cache.Script, dataRowList, getCassandra)
	}

	var clickhouseCache = biz.findClickhouse(redis, mqttClientId)
	for _, cache := range clickhouseCache {
		println(cache.Script)
		house1, _ := clickhouse.GetClickHouse(mqttClientId, []string{fmt.Sprintf("%s:%d", cache.Host, cache.Port)}, cache.Database, cache.Username, cache.Password)

		err := clickhouseOp.HandleDataRowLists(cache.Table, cache.Script, dataRowList, house1)
		if err != nil {
			zap.S().Errorf("clickhouse 执行异常", zap.Error(err))
		}
	}

	var influxdbCache = biz.findInfluxdb(redis, mqttClientId)
	for _, cache := range influxdbCache {
		db := influxdb2.GetInfluxDb(cache.Host, cache.Token, cache.Port, cache.ID)
		err := influxdbOp.HandleDataRowLists(cache.Bucket, cache.Org, cache.Measurement, cache.Script, dataRowList, db)
		if err != nil {
			zap.S().Errorf("influxdb 执行异常", zap.Error(err))
		}
	}

}

func (biz *TransmitCacheBiz) findMysql(redis *redis.Client, mqttClientId string) []MySQLTransmitCache {
	val := redis.LRange(context.Background(), "transmit:mysql:"+mqttClientId, 0, -1).Val()
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
func (biz *TransmitCacheBiz) findMongo(redis *redis.Client, mqttClientId string) []MongoTransmitCache {
	val := redis.LRange(context.Background(), "transmit:mongo:"+mqttClientId, 0, -1).Val()
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

func (biz *TransmitCacheBiz) findCassandra(redis *redis.Client, mqttClientId string) []CassandraTransmitCache {
	val := redis.LRange(context.Background(), "transmit:cassandra:"+mqttClientId, 0, -1).Val()
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

func (biz *TransmitCacheBiz) findClickhouse(redis *redis.Client, mqttClientId string) []ClickhouseTransmitCache {
	val := redis.LRange(context.Background(), "transmit:clickhouse:"+mqttClientId, 0, -1).Val()
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

func (biz *TransmitCacheBiz) findInfluxdb(redis *redis.Client, mqttClientId string) []InfluxTransmitCache {
	val := redis.LRange(context.Background(), "transmit:influxdb:"+mqttClientId, 0, -1).Val()
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
