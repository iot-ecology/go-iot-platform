package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// false 标识通过 ，true 标识没通过

func CheckPushTime(protocol string, identificationCode string, deviceUid string, timeFromUnix int64) bool {
	preKey := "storage_time"
	key := fmt.Sprintf("%s:%s:%s:%s", preKey, protocol, deviceUid, identificationCode)
	background := context.Background()

	// 获取上次推送的时间
	lastPushTimeStr := globalRedisClient.Get(background, key).Val()
	if lastPushTimeStr != "" {
		lastPushTime, err := strconv.ParseInt(lastPushTimeStr, 10, 64)
		if err != nil {
			fmt.Println("Error parsing last push time:", err)
			return false
		}

		sk := fmt.Sprintf("share:device_info:%s:%s:%s", protocol, deviceUid, identificationCode)
		deviceInfoJsonString := globalRedisClient.Get(background, sk).Val()

		var deviceInfo DeviceInfo
		err = json.Unmarshal([]byte(deviceInfoJsonString), &deviceInfo)
		if err != nil {
			fmt.Println("Error unmarshalling device info:", err)
			return false
		}

		// 计算有效推送时间
		effectivePushTime := lastPushTime + int64(deviceInfo.PushInterval) - int64(deviceInfo.ErrorRate)

		// 检查当前时间是否在有效推送时间内
		return timeFromUnix <= effectivePushTime // 直接返回布尔值
	}

	return false // 表示时间无效或不存在
}

func SetPushTime(protocol string, identificationCode string, deviceUid string, timeFromUnix int64) {
	preKey := "storage_time"
	key := fmt.Sprintf("%s:%s:%s:%s", preKey, protocol, deviceUid, identificationCode)

	background := context.Background()
	err := globalRedisClient.Set(background, key, timeFromUnix, 0)
	if err != nil {
		fmt.Println("Error setting time:", err)
	}
}

