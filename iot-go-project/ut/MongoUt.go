package ut

import (
	"context"
	"igp/glob"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// 计算集合名称
// 参数：
// prefix: 集合名称前缀
// id: 集合ID
// 返回值：
// 计算得到的集合名称
func CalcCollectionName(prefix string, id uint) string {
	return prefix + "_" + strconv.Itoa(int(id%100))
}

// 检查集合名称是否存在，如果不存在则创建集合
// 参数：
// prefix: 集合名称前缀
// collectionName: 集合名称
// 返回值：
// 无
func CheckCollectionAndCreate(prefix, collectionName string) {
	zap.S().Infof("prefix = %s, collectionName = %s", prefix, collectionName)
	db := glob.GMongoClient.Database(glob.GConfig.MongoConfig.Db)

	regex := primitive.Regex{Pattern: "^" + prefix, Options: "i"} // 'i' 表示不区分大小写

	filter := bson.M{"name": regex}

	collectionNames, err := db.ListCollectionNames(context.TODO(), filter)

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

	// 不存在则创建
	if !collectionExists {
		err := db.CreateCollection(context.TODO(), collectionName)
		if err != nil {
			zap.S().Fatal(err)
			return
		}

	}

}
