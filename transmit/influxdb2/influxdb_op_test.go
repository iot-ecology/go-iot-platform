package influxdb2

import (
	"encoding/json"
	"fmt"
	"iot-transmit/common"
	"testing"
	"time"
)

var influxdbOp = InfluxDbOp{}

func TestGet(t *testing.T) {
	host := "127.0.0.1"
	token := "mytoken"
	port := 8086
	client := GetInfluxDb(host, token, port, 1)
	client2 := GetInfluxDb(host, token, port, 1)
	if client == client2 {
		t.Logf("success")
	}

	dataRowList := []common.DataRowList{
		{
			Time:      time.Now().Unix(),
			DeviceUid: "111",
			DataRows: []common.DataRow{
				{
					Name:  "a",
					Value: "测试",
				},
			},
			Nc: "111",
		},
	}

	jsonData, err := json.MarshalIndent(dataRowList, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	fmt.Println(string(jsonData))

	script := `function main(jsonData) {
    
    return jsonData;
}
`
	influxdbOp.HandleDataRowLists("buc", "94addf3e3af15b2f", "sample", script, dataRowList, client)

}
