package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"iot-transmit/common"
	"strconv"
	"testing"
)

var (
	biz = TransmitCacheBiz{}
)

func TestBiz(t *testing.T) {
	add := fmt.Sprintf("%s:%d", "127.0.0.1", 6379)
	grd := redis.NewClient(&redis.Options{
		Addr:     add,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81", // 如果没有设置密码，就留空字符串
		DB:       10,                                 // 使用默认数据库
	})
	s := ` [{"time":1720767928,"device_uid":"1","data":[{"name":"Temperature","value":"23"},{"name":"Humidity",
"value":"30"},{"name":"A","value":"1"}],"nc":"1"}]`
	var dataRowList []common.DataRowList
	json.Unmarshal([]byte(s), &dataRowList)
	biz.Run(grd, "1", dataRowList)
}
func TestInitCache(t *testing.T) {

	script := `function main(jsonData) {
    var c = []
    for (var jsonDatum of jsonData) {
        var time = jsonDatum.Time;
        var arr = []
        var timeField = {
            "FieldName": "time",
            "Value": time
        }
		
        arr.push(timeField)
  		var idd = {
            "FieldName": "id",
            "Value": time
        }
        arr.push(idd)
        for (var e of jsonDatum.DataRows) {
            if (e.Name == "a") {
                var aField = {
                    "FieldName": "name",
                    "Value": e.Value
                }
                arr.push(aField)
            }
        }
        c.push(arr)
    }
    return c;
}
`

	v := ClickhouseTransmitCache{
		ID:       "clickhouse-" + strconv.Itoa(int(2)),
		Host:     "127.0.0.1",
		Port:     9000,
		Username: "default",
		Password: "",
		Database: "default",
		Table:    "NewTable",
		Script:   script,
	}

	add := fmt.Sprintf("%s:%d", "127.0.0.1", 6379)
	grd := redis.NewClient(&redis.Options{
		Addr:     add,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81", // 如果没有设置密码，就留空字符串
		DB:       10,                                 // 使用默认数据库
	})

	jsonData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	//println(string(jsonData))
	grd.LPush(context.Background(), "transmit:clickhouse:"+strconv.Itoa(1), jsonData)
	//grd.LRem(context.Background(), "transmit:clickhouse:"+strconv.Itoa(1), 1, jsonData)

}
