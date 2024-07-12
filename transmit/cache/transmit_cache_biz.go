package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"iot-transmit/clickhouse"
	"iot-transmit/common"
)

type TransmitCacheBiz struct {
}

var (
	clickhouseOp = clickhouse.ClickhouseOp{}
)

func (biz *TransmitCacheBiz) Run(redis *redis.Client, mqttClientId string, dataRowList []common.DataRowList) {
	var mysqlCache = biz.findMysql(redis, mqttClientId)
	for _, cache := range mysqlCache {
		println(cache.Script)
	}
	biz.findMongo(redis, mqttClientId)
	biz.findCassandra(redis, mqttClientId)
	var clickhouseCache = biz.findClickhouse(redis, mqttClientId)
	for _, cache := range clickhouseCache {
		println(cache.Script)
		house1, _ := clickhouse.GetClickHouse(mqttClientId, []string{"127.0.0.1:9000"}, cache.Database,
			cache.Username, cache.Password)

		clickhouseOp.HandleDataRowLists(cache.Table, cache.Script, dataRowList, house1)
	}
	biz.findInfluxdb(redis, mqttClientId)

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
