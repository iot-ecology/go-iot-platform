package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	server := "192.168.3.101:3332"
	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	message := "uid:1:admin:admin"
	fmt.Println("发送:", message)
	writer.WriteString(message + "\n")
	writer.Flush()

	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取失败:", err)
			return
		}
		fmt.Println("接收:", response)

		// 检查响应是否包含 "成功识别设备编码"
		if strings.Contains(response, "成功识别设备编码")  || strings.Contains(response, "数据已处理"){
			fmt.Println("设备编码成功识别，继续发送消息")
			// 这里可以添加发送消息的逻辑
			// 例如：发送相同的消息或者发送新的消息
			writer.WriteString(  "data\n")
			writer.Flush()
		} else {
			fmt.Println("设备编码未识别，停止发送消息")
			break
		}
		time.Sleep(10*time.Second)
	}
}