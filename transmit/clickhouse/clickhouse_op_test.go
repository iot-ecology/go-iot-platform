/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

	house1, _ := GetClickHouse("1", []string{"127.0.0.1:9000"}, "default", "default", "default")
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
