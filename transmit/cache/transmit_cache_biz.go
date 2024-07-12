package cache

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

type TransmitCacheBiz struct {
}

func (biz *TransmitCacheBiz) Run(redis *redis.Client, mqttClientId string) {
	biz.findMysql(redis, mqttClientId)
	biz.findMongo(redis, mqttClientId)
	biz.findCassandra(redis, mqttClientId)
	biz.findClickhouse(redis, mqttClientId)
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
