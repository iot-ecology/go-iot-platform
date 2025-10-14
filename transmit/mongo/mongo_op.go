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
	"context"
	"fmt"
	"github.com/dop251/goja"
	"go.uber.org/zap"
	"iot-transmit/common"
	"net/url"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 全局map，用于存储id和MongoDB客户端的映射
var mongoClientMap = make(map[string]*mongo.Client)

// 用于同步的互斥锁
var mu sync.Mutex

// GetMongoDBClient 根据提供的参数获取或创建MongoDB客户端
func GetMongoDBClient(Host, Username, Password, Db string, Port int, id string) (*mongo.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	// 使用id作为map的key
	key := id

	client, exists := mongoClientMap[key]
	if exists {
		zap.S().Infof("Reusing existing MongoDB client for id: %d.", id)
		return client, nil
	}

	encodedUsername := url.QueryEscape(Username)
	encodedPassword := url.QueryEscape(Password)

	// 构建MongoDB连接字符串
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", encodedUsername, encodedPassword, Host, Port, Db)

	if Username == "admin" {
		connectionString = connectionString + "?authSource=admin"
	}
	// 创建MongoDB客户端
	clientOptions := options.Client().ApplyURI(connectionString)
	newClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// 检查连接是否成功
	err = newClient.Ping(context.TODO(), nil)
	if err != nil {
		zap.S().Fatalf("Failed to connect to MongoDB: %s", err)
		return nil, err
	}

	// 将新的MongoDB客户端存储到map中
	mongoClientMap[key] = newClient
	zap.S().Infof("New MongoDB client for database '%s' with id %d has been initiated.", Db, id)
	return newClient, nil
}

type MongoOp struct {
}

func (op *MongoOp) RunScript(param []common.DataRowList, script string) []map[string]interface{} {
	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Error("JS代码有问题！")
		return nil
	}
	var fn func(string2 []common.DataRowList) []map[string]interface{}
	get := vm.Get("main")
	if get != nil {

		err = vm.ExportTo(get, &fn)
		if err != nil {
			zap.S().Error("Js函数映射到 Go 函数失败！")
			return nil
		}
		a := fn(param)
		return a
	}
	return nil
}

func (op *MongoOp) HandleDataRowLists(

	Database, Collection, Script string, dataRowList []common.DataRowList, client *mongo.Client) error {
	res := op.RunScript(dataRowList, Script)
	if res != nil {

		var toInsert []interface{}

		for _, v := range res {
			toInsert = append(toInsert, v)
		}
		_, err := client.Database(Database).Collection(Collection).InsertMany(context.TODO(), toInsert)

		if err != nil {
			zap.S().Errorf("插入数据失败！%+v", err)
		}
	}
	return nil

}
