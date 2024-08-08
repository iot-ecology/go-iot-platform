package ut

import (
	"context"
	"go.uber.org/zap"
	"igp/glob"
	"strconv"
)

// CalcBucketName 函数根据前缀、协议和id计算桶名
// prefix: 桶名前缀
// protocol: 使用的协议
// id: 桶的ID
// 返回值: 计算得到的桶名
func CalcBucketName(prefix, protocol string, id uint) string {
	return prefix + "_" + protocol + "_" + strconv.Itoa(int(id%100))
}

// CheckBucketNameAndCreate 检查桶名称是否存在，如果不存在则创建桶
// 参数：
//     bucket: 桶名称
// 返回值：
//     无
func CheckBucketNameAndCreate(bucket string) {
	name, err := glob.GInfluxdb.OrganizationsAPI().FindOrganizationByName(context.Background(), glob.GConfig.InfluxConfig.Org)
	if err != nil {
		zap.S().Errorf(err.Error())
	}

	if name == nil {
		return
	}
	id, err := glob.GInfluxdb.BucketsAPI().FindBucketsByOrgID(context.Background(), *name.Id)

	bucketExists := false

	for _, bucketID := range *id {
		if bucketID.Name == bucket {
			bucketExists = true
			break
		}
	}
	if !bucketExists {
		withName, err := glob.GInfluxdb.BucketsAPI().CreateBucketWithName(context.Background(), name, bucket)
		if err != nil {
			zap.S().Errorf(err.Error())
		} else {
			zap.S().Infof("Bucket %s created", withName)
		}

	}

}
