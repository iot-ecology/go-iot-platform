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

package mongo

import (
	"encoding/json"
	"fmt"
	"iot-transmit/common"
	"testing"
	"time"
)

var mongoDbOp = MongoOp{}

func TestGet(t *testing.T) {
	// 示例：获取或创建MongoDB客户端
	Host := "localhost"
	Username := "admin"
	Password := "admin"
	Db := "iot"
	Port := 27017
	client, err := GetMongoDBClient(Host, Username, Password, Db, Port, 1)
	if err != nil {
		t.Fatal(err)
	}
	client2, err2 := GetMongoDBClient(Host, Username, Password, Db, Port, 1)
	if err2 != nil {
		t.Fatal(err)
	}
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
	mongoDbOp.HandleDataRowLists("iot", "sample", script, dataRowList, client)
}
