---
publishDate: 2024-08-05T00:00:00Z
author: Zen HuiFer
title: COAP方式接入设备
excerpt: 在Go IoT 开发平台中使用 COAP 接入系统的细节
image: "~/assets/images/communication-1439132_1280.jpg"

category: 教程
tags:
  - Go
  - 物联网
  - 开发平台
  - 应用
  - 实战
---



## 接入流程

演示程序使用Go语言编写，单机环境下COAP服务端口为5683。完整客户端程序代码如下。



```go
package main

import (
	"encoding/json"
	"github.com/dustin/go-coap"
	"log"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
}

func main() {
	c, err := coap.Dial("udp", "localhost:5683")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	auth(c)
	data(c)
}

func data(c*coap.Conn) {
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 12345,
		Payload:   []byte("test"),
	}

	path := "/data"

	req.SetOption(coap.ETag, "weetag")
	req.SetOption(coap.MaxAge, 3)
	req.SetPathString(path)



	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		log.Printf("Response payload: %s", rv.Payload)
	}
}

func auth( c *coap.Conn) {
	auth := Auth{
		Username: "admin",
		Password: "admin",
		DeviceId: "1234567890",
	}
	marshal, _ := json.Marshal(auth)
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 12345,
		Payload:   marshal,
	}

	path := "/auth"

	req.SetOption(coap.ETag, "weetag")
	req.SetOption(coap.MaxAge, 3)
	req.SetPathString(path)

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		log.Printf("Response payload: %s", rv.Payload)
	}
}

```



>   注意事项：一定要使用相同的链接完成 `/auth` 和 `/data` 的数据请求。服务端会对链接信息做校验，两个链接分别完成 `/auth` 和 `/data` 是无法正常交互的。

上述程序执行后效果如下

```
2024/08/06 10:16:32 Response payload: 安全认证成功
2024/08/06 10:16:32 Response payload: 数据处理成功
```







