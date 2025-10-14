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

package mysql

import (
	"iot-transmit/common"
	"testing"
)

func TestGet(t *testing.T) {
	connection, err := InitMySQLConnection("root", "127.0.0.1", "root123@", "iot", 3306, "1")
	if err != nil {
		t.Error(err)
	}
	sqlConnection, err := InitMySQLConnection("root", "127.0.0.1", "root123@", "iot", 3306, "1")

	if err != nil {
		t.Error(err)
	}
	if connection == sqlConnection {
		t.Logf("success ")
	}

}

func TestMock(t *testing.T) {
	var op = MysqlOp{}
	dataRowList := common.DataRowList{
		Time:               1650000000,
		DeviceUid:          "device123",
		IdentificationCode: "code456",
		DataRows: []common.DataRow{
			{
				Name:  "Temperature",
				Value: "22.5",
			},
			{
				Name:  "Humidity",
				Value: "45",
			},
		},
		Nc: "nc789",
	}

	var scrip = `
function main(jsonData) {
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

	var arr []common.DataRowList
	arr = append(arr, dataRowList)
	script := op.RunScript(arr, scrip)
	print(script)
}
