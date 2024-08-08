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
