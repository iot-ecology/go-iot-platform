package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net"
	"strings"
	"sync"
)

type Server struct {
	host string
	port string
	mu   sync.Mutex // 保护deviceIdMap的互斥锁
}

type Client struct {
	conn net.Conn
}

type Config struct {
	Host string
	Port string
}

func New(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go server.handleClient(client)
	}
}

func (server *Server) handleClient(client *Client) {
	defer func(conn net.Conn) {
		log.Printf("Client %s disconnected.\n", conn.RemoteAddr())

		server.mu.Lock()
		defer server.mu.Unlock()

		RemoveUid(client.conn.RemoteAddr().String())
		err := conn.Close()
		if err != nil {

		}
	}(client.conn)

	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return // 关闭连接并退出goroutine
		}

		server.handleMessage(client, message)
	}
}

func (server *Server) handleMessage(client *Client, message string) {

	// 判断这个客户端是否建立过uid映射，没有的话不处理数据

	ok := getUid(client.conn.RemoteAddr().String())

	if ok != "" {
		handlerData(message, client)
	} else {

		uid := handlerUid(message)
		if uid == "" {
			clientWrite(client, "请发送uid:xxx格式的消息进行设备ID映射。\n")

			return
		} else {
			server.mu.Lock()
			defer server.mu.Unlock()

			split := strings.Split(uid, ":")

			if len(split) != 3 {
				clientWrite(client, "uid格式错误.\n")
				return
			}
			var deviceId = split[0]
			var username = split[1]
			var password = split[2]

			usernameC, passwordC := FindDeviceMappingUP(deviceId)
			zap.S().Infof("device_id: %s", deviceId)
			zap.S().Infof("有效账号密码 username: %s, password: %s", usernameC, passwordC)
			if username != usernameC || password != passwordC {
				clientWrite(client, "账号密码不正确.\n")
				return
			}

			m := globalRedisClient.HRandField(context.Background(), "tcp_uid_f:"+globalConfig.NodeInfo.Name, -1).Val()

			if int64(len(m)) <= globalConfig.NodeInfo.Size {

				storageUid(deviceId, client.conn.RemoteAddr().String())

				TcpMap[client.conn.RemoteAddr().String()] = client

				clientWrite(client, "成功识别设备编码.\n")

				return
			} else {
				clientWrite(client, "当前服务器已满载.\n")
				client.conn.Close()
				return
			}
		}

	}

}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func FindDeviceMappingUP(deviceId string) (string, string) {
	//  从redis中根据deviceId获取用户名和密码
	val := globalRedisClient.HGet(context.Background(), "auth:tcp", deviceId).Val()
	var auth Auth
	err := json.Unmarshal([]byte(val), &auth)
	if err != nil {
		return "", ""
	}
	return auth.Username, auth.Password
}

func getUid(remoteAdd string) string {
	val := globalRedisClient.HGet(context.Background(), "tcp_uid_f:"+globalConfig.NodeInfo.Name, remoteAdd).Val()
	return val
}

func storageUid(uid, remoteAdd string) {
	globalRedisClient.HSet(context.Background(), "tcp_uid:"+globalConfig.NodeInfo.Name, uid, remoteAdd)
	globalRedisClient.HSet(context.Background(), "tcp_uid_f:"+globalConfig.NodeInfo.Name, remoteAdd, uid)

}
func RemoveUid(remoteAdd string) {
	val := globalRedisClient.HGet(context.Background(), "tcp_uid_f:"+globalConfig.NodeInfo.Name, remoteAdd).Val()
	if val == "" {
		return
	} else {
		globalRedisClient.HDel(context.Background(), "tcp_uid:"+globalConfig.NodeInfo.Name, val)
		globalRedisClient.HDel(context.Background(), "tcp_uid_f:"+globalConfig.NodeInfo.Name, remoteAdd)
	}
}

type TcpMessage struct {
	Uid     string `json:"uid"`
	Message string `json:"message"`
}

var TcpMap = make(map[string]*Client)

func handlerData(message string, client *Client) {

	zap.S().Debugf("处理消息: %s  客户端: %s\n", message, client.conn.RemoteAddr().String())

	s := getUid(client.conn.RemoteAddr().String())

	// 创建 TCPMessage 实例并序列化为 JSON
	tcpMsg := TcpMessage{
		Uid:     s,
		Message: message,
	}
	jsonData, err := json.Marshal(tcpMsg)
	if err != nil {
		zap.S().Errorf("Error marshalling TCP message to JSON: %v", err)
		return
	}
	SetLastOpTime(client.conn.RemoteAddr().String())
	s2 := PushToQueue("pre_tcp_handler", jsonData)
	if s2 != nil {
		zap.S().Errorf("Error pushing TCP message to queue: %v", s2)
		clientWrite(client, fmt.Sprintf("数据处理异常 %s.\n", s2))
	}

	clientWrite(client, "数据已处理.\n")

}

func clientWrite(client *Client, msg string) {
	_, err := client.conn.Write([]byte(msg))
	if err != nil {
		zap.S().Error(err)
		return
	}
}

// handlerUid 用于从消息中提取设备ID
func handlerUid(message string) string {
	if strings.HasPrefix(message, "uid:") {
		return strings.TrimSpace(message[4:])
	}
	return ""
}
