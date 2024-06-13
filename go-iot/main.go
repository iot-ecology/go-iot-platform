package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"
)

var globalConfig ServerConfig

func main() {
	InitLog()

	var configPath string
	flag.StringVar(&configPath, "config", "app-node1.yml", "Path to the config file")
	flag.Parse()

	yfile, err := os.ReadFile(configPath)
	if err != nil {
		zap.S().Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(yfile, &globalConfig)
	if err != nil {
		zap.S().Fatalf("error: %v", err)
	}

	zap.S().Infof("node name = %v , host = %v , port = %v", globalConfig.NodeInfo.Name, globalConfig.NodeInfo.Host, globalConfig.NodeInfo.Port)
	InitRabbitCon()

	initGlobalRedisClient(globalConfig.RedisConfig)

	beforeStart()
	startHttp()

}

func beforeStart() {
	go removeOldData()

	go BeatTask(globalConfig.NodeInfo)
	go ListenerBeat()
	go CBeat()
	go timerNoHandlerConfig()
}
func removeOldData() {
	HandlerOffNode(globalConfig.NodeInfo.Name)
}

// CBeat 是一个无限循环函数，用于定时检查心跳并进行处理
func CBeat() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			lock := NewRedisDistLock(globalRedisClient, "c_beat")
			if lock.TryLock() {
				service, err := GetThisTypeService()
				if err == nil {
					if processHeartbeats(service) {
						return
					}
				}
				lock.Unlock()

			} else {

				zap.S().Error("没有获取到 c_beat 处理的锁")

			}
		}
	}

}

// processHeartbeats 函数用于处理心跳信息
//
// 参数：
// service []NodeInfo - 节点信息切片，包含待处理的心跳信息
//
// 返回值：
// bool - 若成功处理某个节点的心跳信息则返回true，否则返回false
func processHeartbeats(service []NodeInfo) bool {
	for _, info := range service {
		if !SendBeat(&info, "beat") {
			globalRedisClient.HDel(context.Background(), "register:"+globalConfig.NodeInfo.Type, info.Name)
			HandlerOffNode(info.Name)
		}

	}
	return false
}

// HandlerOffNode 函数用于处理节点下线的情况
//
// 参数：
//
//	node_name string - 节点名称
//
// 返回值：
//
//	bool - 如果处理成功返回false，否则返回true
func HandlerOffNode(nodeName string) {
	// 清除节点负载计数器

	// 获取节点对应的MQTT客户端ID
	mqttClientIds := GetBindClientId(nodeName)
	for _, ele := range mqttClientIds {
		// 获取MQTT客户端ID对应的MQTT配置
		cf := GetUseConfig(ele)
		if cf == "" {
			zap.S().Errorf("HandlerOffNode Error get mqtt config, client id = %s", ele)
			continue

		} else {

			var config MqttConfig
			bytes := []byte(cf)
			err := json.Unmarshal(bytes, &config)
			if err != nil {
				zap.S().Errorf("HandlerOffNode Error unmarshalling JSON: %s", err)
				continue
			}
			// 移除节点和MQTT客户端ID的关系
			RemoveBindNode(config.ClientId, nodeName)
		}

	}
}

func startHttp() {
	http.HandleFunc("/beat", HttpBeat)
	http.HandleFunc("/create_mqtt", CreateMqttClientHttp)
	http.HandleFunc("/node_list", NodeList)
	http.HandleFunc("/node_using_status", NodeUsingStatus)
	http.HandleFunc("/mqtt_config", GetUseMqttConfig)
	http.HandleFunc("/no_mqtt_config", GetNoUseMqttConfig)
	http.HandleFunc("/remove_mqtt_client", RemoveMqttClient)
	http.HandleFunc("/push_data", PushMqttData)
	//http.HandleFunc("/pprof-test", Handler)
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/public_create_mqtt", PubCreateMqttClientHttp)
	http.HandleFunc("/public_remove_mqtt_client", PubRemoveMqttClient)
	http.HandleFunc("/public_push_data", PubPushMqttData)

	if err := http.ListenAndServe(":"+strconv.Itoa(globalConfig.NodeInfo.Port), nil); err != nil {
		zap.S().Fatalf("Failed to start server: %s", err)
	}
}

func timerNoHandlerConfig() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			// 执行定时任务的逻辑

			noHandlerConfig()
		}
	}
}

// noHandlerConfig 函数用于处理没有MQTT客户端处理配置的情况
//
// 函数首先创建一个Redis分布式锁，锁的名字为 "no_handler_config_lock"
// 如果成功获取到锁，则打印日志 "获取处理 no_handler_config 的锁"
// 随后调用 GetNoUseConfig 函数获取未使用的MQTT客户端配置列表
// 遍历配置列表，对每一个配置调用 PubCreateMqttClientOp 函数创建MQTT客户端
// 如果创建失败，则继续处理下一个配置
// 最后释放锁资源
//
// 如果获取锁失败，则打印错误日志 "没有获取到 no_handler_config 处理的锁"
func noHandlerConfig() {
	lock := NewRedisDistLock(globalRedisClient, "no_handler_config_lock")
	if lock.TryLock() {

		zap.S().Debugf("获取处理 no_handler_config 的锁")
		configList := GetNoUseConfig()

		for _, conf := range configList {

			if PubCreateMqttClientOp(conf) == -1 {
				continue
			}
		}
		lock.Unlock()
	} else {
		zap.S().Error("没有获取到 no_handler_config 处理的锁")
	}

}

func PubCreateMqttClientOp(conf string) int {
	lose := GetSizeLose("")

	if lose != nil {

		if SendCreateMqttMessage(lose, conf) {
			return 1
		} else {
			zap.S().Errorf("发送 MQTT 客户端创建请求异常 %v", lose)
			return -1
		}

	} else {
		zap.S().Error("没有找到可用节点")
		return -1
	}
}
