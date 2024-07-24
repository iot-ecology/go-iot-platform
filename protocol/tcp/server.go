package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net"
	"strings"
	"sync"
)

type Server struct {
	host        string
	port        string
	deviceIdMap map[string]*Client
	remoteIpMap map[string]string
	mu          sync.Mutex // 保护deviceIdMap的互斥锁
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
		host:        config.Host,
		port:        config.Port,
		deviceIdMap: make(map[string]*Client),
		remoteIpMap: make(map[string]string),
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

		// 更新或添加设备ID映射
		s := server.remoteIpMap[client.conn.RemoteAddr().String()]
		delete(server.deviceIdMap, s)
		delete(server.remoteIpMap, client.conn.RemoteAddr().String())
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

	_, ok := server.remoteIpMap[client.conn.RemoteAddr().String()]

	if ok {
		handlerData(server,message, client)
	} else {

		deviceId := handlerUid(message)
		if deviceId == "" {
			clientWrite(client, "请发送uid:xxx格式的消息进行设备ID映射。\n")

			return
		} else {
			server.mu.Lock()
			defer server.mu.Unlock()

			// 更新或添加设备ID映射
			server.deviceIdMap[deviceId] = client
			server.remoteIpMap[client.conn.RemoteAddr().String()] = deviceId
			clientWrite(client, "成功识别设备编码.\n")

			return

		}

	}

}

type TcpMessage struct {
	Uid string `json:"uid"`
	Message string `json:"message"`
}

func handlerData(server *Server, message string, client *Client) {
	println(message)

	zap.S().Debugf("处理消息: %s  客户端: %s\n", message, client.conn.RemoteAddr().String())

	s := server.remoteIpMap[client.conn.RemoteAddr().String()]

	// 创建 MQTTMessage 实例并序列化为 JSON
	mqttMsg := TcpMessage{
		Uid: s,
		Message:      message,
	}
	jsonData, err := json.Marshal(mqttMsg)
	if err != nil {
		zap.S().Errorf("Error marshalling MQTT message to JSON: %v", err)
		return
	}
	PushToQueue("pre_tcp_handler", jsonData)


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
