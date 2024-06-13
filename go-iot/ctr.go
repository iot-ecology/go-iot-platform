package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func SendCreateMqttMessage(node *NodeInfo, param string) bool {

	url := fmt.Sprintf("http://%s:%d/create_mqtt", node.Host, node.Port)
	data := []byte(param)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		zap.S().Fatalf("Error creating request: %s", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	//bodyBytes, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//bodyString := string(bodyBytes)

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	bodyString := buf.String()

	var m map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &m)
	if err != nil {
		zap.S().Error("Error: %s", err)
	}

	zap.S().Infof("Response Status: %v , body = %v", resp.Status, bodyString)
	status := m["status"].(float64)
	return status == 200

}

func SendBeat(node *NodeInfo, param string) bool {

	url := fmt.Sprintf("http://%s:%d/beat", node.Host, node.Port)
	data := []byte(param)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		zap.S().Fatalf("Error creating request: %s", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//zap.S().Fatalf("发送请求异常:", err)
		return false
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return resp.Status == "200 OK"

}

func HttpBeat(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return
	}
	// 向客户端发送响应消息
	_, _ = fmt.Fprintf(w, "ok")
}

func CreateMqttClientHttp(w http.ResponseWriter, r *http.Request) {
	// 确保请求方法是POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return
	}
	// 读取请求体
	//body, err := ioutil.ReadAll(r.Body)

	//if err != nil {
	//	http.Error(w, "Error reading request body", http.StatusInternalServerError)
	//	return
	//}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
		return
	}
	var body = buf.Bytes()
	var config MqttConfig
	err = json.Unmarshal(body, &config)
	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if CheckHasConfig(config) {
		err := json.NewEncoder(w).Encode(map[string]any{"status": 400, "message": "已经存在客户端id"})
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
		return

	} else {

		usz := CreateMqttClient(config)

		if usz == -1 {
			err := json.NewEncoder(w).Encode(map[string]any{"status": 400, "message": "达到最大客户端数量"})
			if err != nil {
				zap.S().Errorf("Error: %+v", err)
			}
			return

		}
		if usz == -2 {
			err := json.NewEncoder(w).Encode(map[string]any{"status": 400, "message": "MQTT客户端配置异常"})
			if err != nil {
				zap.S().Errorf("Error: %+v", err)
			}
			return
		} else {
			AddNoUseConfig(config, body)
			BindNode(config, globalConfig.NodeInfo.Name)
			err := json.NewEncoder(w).Encode(map[string]any{"status": 200, "message": "创建成功", "size": usz})
			if err != nil {
				zap.S().Errorf("Error: %+v", err)
			}
			return

		}
	}

}

func PubCreateMqttClientHttp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return
	}
	// 确保请求方法是POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}
	var str = buf.String()

	w.Header().Set("Content-Type", "application/json")

	if PubCreateMqttClientOp(str) == 1 {
		err := json.NewEncoder(w).Encode(map[string]any{"status": 200, "message": "创建成功"})
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
		return
	} else {
		err := json.NewEncoder(w).Encode(map[string]any{"status": 200, "message": "创建失败"})
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
		return

	}
}

func PubRemoveMqttClient(w http.ResponseWriter, r *http.Request) {
	if cros(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()

	// 获取"id"查询参数的值
	id := query.Get("id")
	nodeName := FindMqttClientId(id)
	if nodeName == "" {
		err := json.NewEncoder(w).Encode(map[string]any{"status": 200, "message": "节点未找到"})
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
		return
	}
	info, err := GetNodeInfo(nodeName)
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}

	sendRemoveMqttClient(id, info)

	// fixme: 找到节点并发送请求
	err = json.NewEncoder(w).Encode(map[string]any{"status": 200, "message": "移除成功"})
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}

}

func sendRemoveMqttClient(id string, nodeinfo NodeInfo) {
	baseUrl := fmt.Sprintf("http://%s:%d/remove_mqtt_client?id=%s", nodeinfo.Host, nodeinfo.Port, id)

	// 发送 GET 请求
	resp, err := http.Get(baseUrl)
	if err != nil {
		zap.S().Error("Error sending GET request: %s", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
	}(resp.Body)

	// 打印响应状态码
	zap.S().Infof("Response Status Code: %d", resp.StatusCode)

}

func PubPushMqttData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return
	}
	// 确保请求方法是POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}
	var str = buf.String()

	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()

	// 获取"id"查询参数的值
	id := query.Get("id")
	nodeName := FindMqttClientId(id)
	if nodeName == "" {
		err := json.NewEncoder(w).Encode(map[string]any{"status": 200, "message": "节点未找到"})
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
		return
	}
	info, err := GetNodeInfo(nodeName)
	if err != nil {
		zap.S().Error("Error getting node info: %s", err)
	}

	sendPushMqttData(info, str)

	err = json.NewEncoder(w).Encode(map[string]any{"status": 200, "message": "消息发送成功"})
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}

}

func sendPushMqttData(node NodeInfo, param string) bool {

	url := fmt.Sprintf("http://%s:%d/push_data", node.Host, node.Port)
	data := []byte(param)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		zap.S().Fatalf("Error creating request: %s", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
	}(resp.Body)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}
	var bodyString = buf.String()

	zap.S().Infof("Response Status: %v , body = %v", resp.Status, bodyString)
	return resp.Status == "200 OK"

}

func NodeList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	service, err := GetThisTypeService()
	if err != nil {
		zap.S().Errorf("节点列表获取失败: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 设置响应的Content-Type为application/json
	w.Header().Set("Content-Type", "application/json")

	// 编码并发送JSON响应
	err = json.NewEncoder(w).Encode(map[string]any{
		"status":  200,
		"message": "创建成功",
		"data":    service,
	})
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}
}

func NodeUsingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return
	}
	//if r.Method != http.MethodGet {
	//	http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
	//	return
	//}

	w.Header().Set("Content-Type", "application/json")

	service, err := GetThisTypeService()
	if err != nil {
		zap.S().Errorf("节点列表获取失败: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 定义一个结构体用于存放节点名称和大小
	type NodeInfo struct {
		Name        string       `json:"name"`
		Size        int64        `json:"size"`
		ClientIds   []string     `json:"client_ids"`
		ClientInfos []MqttConfig `json:"client_infos"`
		MaxSize     int64        `json:"max_size"`
	}

	// 初始化一个NodeInfo类型的切片
	var nodeInfos []NodeInfo

	// 遍历service中的每个info
	for _, info := range service {

		size := globalRedisClient.SCard(context.Background(), "node_bind:"+info.Name).Val()

		var mc []MqttConfig

		for _, el := range GetBindClientId(info.Name) {
			// 假设GetUseConfig函数返回配置的JSON字符串和错误
			configJSON := GetUseConfig(el)

			var config MqttConfig
			b := []byte(configJSON)
			err := json.Unmarshal(b, &config)
			if err != nil {
				zap.S().Fatalf("HandlerOffNode Error unmarshalling JSON: %s", err)
			}
			mc = append(mc, config)

		}

		// 创建NodeInfo实例并添加到切片中
		nodeInfos = append(nodeInfos, NodeInfo{
			Name:        info.Name,
			Size:        size,
			MaxSize:     info.Size,
			ClientIds:   GetBindClientId(info.Name),
			ClientInfos: mc,
		})

	}
	err = json.NewEncoder(w).Encode(map[string]any{
		"status":  200,
		"message": "成功",
		"data":    nodeInfos,
	})
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}
}

func GetUseMqttConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// 从请求中解析查询参数
	query := r.URL.Query()

	// 获取"id"查询参数的值
	id := query.Get("id")

	// 检查是否获取到了id参数
	if id == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	// 假设GetUseConfig函数返回配置的JSON字符串和错误
	configJSON := GetUseConfig(id)

	var config MqttConfig
	b := []byte(configJSON)
	err := json.Unmarshal(b, &config)
	if err != nil {
		zap.S().Fatalf("HandlerOffNode Error unmarshalling JSON: %s", err)
		err := json.NewEncoder(w).Encode(map[string]any{
			"status":  200,
			"message": "Success",
			"data":    nil,
		})
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
		return
	}

	// 将配置信息编码为JSON并发送给客户端
	err = json.NewEncoder(w).Encode(map[string]any{
		"status":  200,
		"message": "Success",
		"data":    config,
	})
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}
}

func GetNoUseMqttConfig(w http.ResponseWriter, r *http.Request) {
	if cros(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// 从请求中解析查询参数
	query := r.URL.Query()

	// 获取"id"查询参数的值
	id := query.Get("id")

	// 检查是否获取到了id参数
	if id == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	// 假设GetUseConfig函数返回配置的JSON字符串和错误
	configJSON := GetNoUseConfigById(id)

	var config MqttConfig
	b := []byte(configJSON)
	err := json.Unmarshal(b, &config)
	if err != nil {
		zap.S().Fatalf("HandlerOffNode Error unmarshalling JSON: %s", err)
		err := json.NewEncoder(w).Encode(map[string]any{
			"status":  200,
			"message": "Success",
			"data":    nil,
		})
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(map[string]any{
		"status":  200,
		"message": "Success",
		"data":    config,
	})
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}
}
func RemoveMqttClient(w http.ResponseWriter, r *http.Request) {
	if cros(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// 从请求中解析查询参数
	query := r.URL.Query()

	// 获取"id"查询参数的值
	id := query.Get("id")

	// 检查是否获取到了id参数
	if id == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	// 假设GetUseConfig函数返回配置的JSON字符串和错误
	StopMqttClient(id)

	// 将配置信息编码为JSON并发送给客户端
	err := json.NewEncoder(w).Encode(map[string]any{
		"status":  200,
		"message": "Success",
		"data":    "已停止",
	})
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}
}
func PushMqttData(w http.ResponseWriter, r *http.Request) {
	if cros(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	type paramStruct struct {
		ClientID string `json:"client_id"`
		Topic    string `json:"topic"`
		QOS      byte   `json:"qos"`
		Retained bool   `json:"retained"`
		Payload  string `json:"payload"`
	}

	var params paramStruct

	// 解析请求体到 params 结构体
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	PushMqttMsg(params.ClientID, params.Topic, params.QOS, params.Retained, params.Payload)

	// 将配置信息编码为JSON并发送给客户端
	err = json.NewEncoder(w).Encode(map[string]any{
		"status":  200,
		"message": "Success",
		"data":    "已发送",
	})
	if err != nil {
		zap.S().Errorf("Error: %+v", err)
	}
}

func cros(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 允许的请求头部
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	// 非简单请求时，浏览器会先发送一个预检请求(OPTIONS)，这里处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 200 OK 也可以
		return true
	}
	return false
}
