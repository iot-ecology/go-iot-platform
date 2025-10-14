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

package cassandra

import (
	"encoding/json"
	"fmt"
	"iot-transmit/common"
	"testing"
	"time"
)

var cassandraOp = CassandraOp{}

func TestGet(t *testing.T) {
	ips := []string{"127.0.0.1"}
	username := "admin"
	password := "admin"

	session1, err := GetCassandra(ips, username, password, 1)
	if err != nil {
		t.Fatal(err)
	}
	session2, err2 := GetCassandra(ips, username, password, 1)
	if err2 != nil {
		t.Fatal(err2)
	}
	if session1 == session2 {
		t.Logf("session1 == session2")
	}

	dataRowList := []common.DataRowList{
		{
			Time:      time.Now().Unix(),
			DeviceUid: "111",
			DataRows: []common.DataRow{
				{
					Name:  "a",
					Value: "1",
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
	//script = "function main(jsonData) {  return jsonData}"
	cassandraOp.HandleDataRowLists("sample", "sample", script, dataRowList, session1)

}
