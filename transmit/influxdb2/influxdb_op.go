package influxdb2

import (
	"fmt"
	"github.com/dop251/goja"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"go.uber.org/zap"
	"iot-transmit/common"
	"sync"
	"time"
)

var influxClientMap = make(map[string]influxdb2.Client)
var mu sync.Mutex

func GetInfluxDb(host, token string, port int, id string) influxdb2.Client {
	mu.Lock()
	defer mu.Unlock()

	client, exists := influxClientMap[id]
	if exists {
		zap.S().Infof("Reusing existing InfluxDB client for id: %d.\n", id)
		return client
	}

	client = influxdb2.NewClient(fmt.Sprintf("http://%s:%d", host, port), token)

	influxClientMap[id] = client
	zap.S().Infof("New InfluxDB client for id %d has been initiated.\n", id)

	return client
}

type InfluxDbOp struct{}

func (op *InfluxDbOp) HandleDataRowLists(

	Bucket, Org, Measurement, Script string, dataRowList []common.DataRowList, client influxdb2.Client) error {
	writeAPI := client.WriteAPI(Org, Bucket)
	res := op.RunScript(dataRowList, Script)

	for _, re := range res {
		op.Save(re, writeAPI, Measurement)
	}
	writeAPI.Flush()
	return nil
}

func (op *InfluxDbOp) Save(dt common.DataRowList, api api.WriteAPI, measurement string) {
	timeFromUnix := time.Unix(dt.Time, 0)
	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("storage_time", time.Now().Unix()).
		AddField("push_time", dt.Time).
		SetTime(timeFromUnix)

	for _, row := range dt.DataRows {
		p.AddField(row.Name, row.Value)
	}
	api.WritePoint(p)

}

// 注意事项： 每个数据要会处理为数字类型的字符串。
func (op *InfluxDbOp) RunScript(dataRowList []common.DataRowList, script string) []common.DataRowList {

	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Errorf("JS代码有问题！")
		return nil
	}
	var fn func(string2 []common.DataRowList) []common.DataRowList
	get := vm.Get("main")
	if get != nil {

		err = vm.ExportTo(get, &fn)
		if err != nil {
			zap.S().Errorf("Js函数映射到 Go 函数失败！")
			return nil
		}
		a := fn(dataRowList)
		return a
	}
	return nil
}
