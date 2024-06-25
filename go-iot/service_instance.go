package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

// GetSizeLose 用户获取使用数量最小的节点
// 参数 pass_node_name: 忽略的节点名称
// 返回值 *NodeInfo: 返回节点信息中大小损失最小的节点指针，若不存在则返回nil
func GetSizeLose(passNodeName string) *NodeInfo {
	service, err := GetThisTypeService()
	if err != nil {
		return nil
	}

	if len(service) == 0 {
		return nil
	}

	var minNodeInfo *NodeInfo
	var minSize int64 = -1

	for _, v := range service {

		if v.Name == passNodeName {
			continue
		}

		i := globalRedisClient.SCard(context.Background(), "node_bind:"+v.Name).Val()

		if i < v.Size {

			if minNodeInfo == nil || v.Size < minSize {
				minNodeInfo = &v
				minSize = v.Size
			}
		} else {
			continue
		}
	}

	if minNodeInfo == nil {
		return nil
	}
	zap.S().Infof("选中节点 = %#v", minNodeInfo)

	return minNodeInfo
}

// GetThisTypeService 从Redis中获取指定类型的服务节点信息列表
//
// 参数：
// 无
//
// 返回值：
// []NodeInfo：返回包含所有服务节点信息的切片
// error：如果获取服务节点信息失败，则返回错误信息
func GetThisTypeService() ([]NodeInfo, error) {

	// 使用 context.Background() 作为请求的上下文
	all, err := globalRedisClient.HGetAll(context.Background(), "register:"+globalConfig.NodeInfo.Type).Result()
	if err != nil {
		return nil, err // 如果有错误，返回错误
	}

	var res []NodeInfo
	for _, v := range all {
		var nodeInfo NodeInfo
		err := json.Unmarshal([]byte(v), &nodeInfo)
		if err != nil {
			return nil, err
		}
		res = append(res, nodeInfo)
	}

	return res, nil

}

// GetNodeInfo 函数根据节点名称从Redis中获取节点信息
//
// 参数：
// name string - 节点名称
//
// 返回值：
// NodeInfo - 节点信息结构体
// error - 如果节点信息不存在或获取失败，则返回错误信息
func GetNodeInfo(name string) (NodeInfo, error) {
	ctx := context.Background()
	val, err := globalRedisClient.HGet(ctx, "register:"+globalConfig.NodeInfo.Type, name).Result()
	if errors.Is(err, redis.Nil) {
		return NodeInfo{}, errors.New("node info not found")
	}
	if err != nil {
		return NodeInfo{}, err
	}

	var nodeInfo NodeInfo
	err = json.Unmarshal([]byte(val), &nodeInfo)
	if err != nil {
		return NodeInfo{}, err
	}

	return nodeInfo, nil
}

// RemoveNodeInfo 从Redis中删除指定名称的节点信息
//
// 参数：
//
//	name string - 要删除的节点名称
//
// 返回值：
//
//	error - 如果删除操作失败，则返回错误信息；否则返回nil
func RemoveNodeInfo(name string) error {
	ctx := context.Background()
	_, err := globalRedisClient.HDel(ctx, "register:"+globalConfig.NodeInfo.Type, name).Result()

	if err != nil {
		return err
	}
	return nil
}

// BeatTask 函数用于创建一个定时任务，每隔一定时间执行一次 Register 函数，
// 将传入的 NodeInfo 结构体作为参数。
//
// 参数 f 是一个 NodeInfo 类型的结构体，包含了需要注册的节点信息。
//
// 该函数通过 time.NewTicker 创建一个定时器，每隔 1 秒触发一次，
// 在定时器的触发信号到来时，执行 Register 函数，完成节点信息的注册。
//
// 该函数没有返回值，会一直运行下去，直到程序被外部因素（如操作系统）终止。
func BeatTask(f NodeInfo) {

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		Register(f)
	}
}

// Register 函数用于将节点信息注册到Redis中
//
// 参数：
// f NodeInfo - 节点信息结构体
//
// 返回值：
// 无
//
// 函数会将传入的节点信息结构体序列化为JSON格式，并存储到Redis中。
// 存储时使用两个Redis key，分别为"beat:{Type}:{Name}"和"register:{Type}"。
// "beat:{Type}:{Name}"用于存储节点的名称，过期时间为1200毫秒。
// "register:{Type}"用于存储节点的JSON信息，以节点名称为field。
// 若序列化失败，则输出日志并终止程序。
func Register(f NodeInfo) {
	jsonData, err := json.Marshal(f)

	zap.S().Debugf("健康数据 data %v", string(jsonData))

	if err != nil {
		zap.S().Fatalf("Error marshalling to JSON: %v", err)
	}
	globalRedisClient.Set(context.Background(), "beat:"+f.Type+":"+f.Name, f.Name, 1200*time.Millisecond)
	globalRedisClient.HSet(context.Background(), "register:"+f.Type, f.Name, jsonData)
	zap.S().Debugf("健康心跳")
}
