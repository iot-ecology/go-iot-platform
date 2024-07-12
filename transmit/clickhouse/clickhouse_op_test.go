package clickhouse

import (
	"encoding/json"
	"fmt"
	"iot-transmit/common"
	"testing"
	"time"
)

var clickhouseOp = ClickhouseOp{}

func TestGet(t *testing.T) {

	house1, _ := GetClickHouse(1, []string{"127.0.0.1:9000"}, "default", "default", "")
	//house2, _ := GetClickHouse(1, []string{"127.0.0.1:9000"}, "default", "default", "")

	//if house1 == house2 {
	//	t.Logf("success ")
	//}
	//query, err := house1.Query(context.Background(), "SELECT * FROM \"default\".\"NewTable\" LIMIT 300 OFFSET 0;")
	//if err != nil {
	//	t.Error(err)
	//} else {
	//	for query.Next() {
	//		var name string
	//		err := query.Scan(&name)
	//		if err != nil {
	//			t.Error(err)
	//		}
	//		t.Logf(" name:%s", name)
	//	}
	//}

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
	print(script)
	clickhouseOp.HandleDataRowLists("NewTable", script, dataRowList, house1)

}
