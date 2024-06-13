package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"igp/glob"
	"igp/models"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type NodeBiz struct{}

var r = rand.New(rand.NewSource(13))

// 获取所有实例列表
func (biz *NodeBiz) randomNode() (*models.NodeInfo, error) {
	all, err := glob.GRedis.HGetAll(context.Background(), "register:mqtt").Result()
	if err != nil {
		return nil, err // 如果有错误，返回错误
	}

	var res []models.NodeInfo
	for _, v := range all {
		var nodeInfo models.NodeInfo
		err := json.Unmarshal([]byte(v), &nodeInfo)
		if err != nil {
			return nil, err
		}
		res = append(res, nodeInfo)
	}

	if len(res) == 0 {
		return nil, nil // 没有可用的节点
	}

	// 随机选择一个节点

	randomIndex := r.Intn(len(res)) // 生成一个随机索引

	return &res[randomIndex], nil // 返回随机选择的节点
}

// SendCreateParam 联通MQTT客户端管理集群创建MQTT客户端
func (biz *NodeBiz) SendCreateParam(mc *models.MqttClient) string {
	node, err2 := biz.randomNode()

	if err2 != nil {
		panic(err2)
		return ""
	}
	url := fmt.Sprintf("http://%s:%d/public_create_mqtt", node.Host, node.Port)

	mqttData := struct {
		Broker   string `json:"broker"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		SubTopic string `json:"sub_topic"`
		ClientID string `json:"client_id"`
	}{
		Broker:   mc.Host,
		Port:     mc.Port,
		Username: mc.Username,
		Password: mc.Password,
		SubTopic: mc.Subtopic,
		ClientID: mc.ClientId,
	}

	jsonData, err := json.Marshal(mqttData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	fmt.Println("Response Status Code:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	fmt.Println("Response Body:", string(body))
	return string(body)
}

func (biz *NodeBiz) SendNodeUsingStatus() string {
	node, err2 := biz.randomNode()

	if err2 != nil {
		panic(err2)
		return ""
	}
	url := fmt.Sprintf("http://%s:%d/node_using_status", node.Host, node.Port)
	req, err := http.NewRequest("get", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	fmt.Println("Response Status Code:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	fmt.Println("Response Body:", string(body))
	return string(body)
}

// SendStopParam 联通MQTT客户端管理集群关闭MQTT客户端
func (biz *NodeBiz) SendStopParam(mc *models.MqttClient) string {
	node, err2 := biz.randomNode()

	if err2 != nil {
		panic(err2)
		return ""
	}
	url := fmt.Sprintf("http://%s:%d/public_remove_mqtt_client?id=%s", node.Host, node.Port, mc.ClientId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	fmt.Println("Response Status Code:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	fmt.Println("Response Body:", string(body))
	return string(body)
}

// SendPushData 联通MQTT客户端管理集群关闭MQTT客户端
func (biz *NodeBiz) SendPushData(mqttClientId string, param []byte) string {
	node, err2 := biz.randomNode()

	if err2 != nil {
		panic(err2)
		return ""
	}
	url := fmt.Sprintf("http://%s:%d/public_push_data?id=%s", node.Host, node.Port, mqttClientId)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	fmt.Println("Response Status Code:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	fmt.Println("Response Body:", string(body))
	return string(body)
}
