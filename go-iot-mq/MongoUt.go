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

package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"strconv"
)

func CalcCollectionName(prefix string, id uint) string {
	return prefix + "_" + strconv.Itoa(int(id%100))
}


func CheckCollectionAndCreate(prefix , collectionName string ) {
	db	 := GMongoClient.Database(globalConfig.MongoConfig.Db)

	regex := primitive.Regex{Pattern: "^" + prefix, Options: "i"} // 'i' 表示不区分大小写

	filter := bson.M{"name": regex}

	collectionNames, err := db.ListCollectionNames(context.TODO(),filter)

	if err != nil {

		zap.S().Fatal(err)

	}

	collectionExists := false
	for _, name := range collectionNames {
		if name == collectionName {
			collectionExists = true
			break
		}
	}
	zap.S().Infof("collection %s exists: %v", collectionName, collectionExists)

	if collectionExists {
		err := db.CreateCollection(context.TODO(), collectionName)
		if err != nil {
			zap.S().Fatal(err)
			return
		}

	}

}
