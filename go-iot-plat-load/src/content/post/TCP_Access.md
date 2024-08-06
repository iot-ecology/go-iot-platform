---
publishDate: 2024-08-05T00:00:00Z
author: Zen HuiFer
title: TCP/IP 方式接入设备
excerpt: 在Go IoT 开发平台中使用 TCP/IP 接入系统的细节
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

1.   假设你现在使用Linux系统，你可以尝试使用下面命令与TCP服务建立通讯

```
nc -v 127.0.0.1 3332
```

2.   认证方式为输入 `uid:${设备id}:{用户名}:{密码}`，测试案例如下。

```
Connection to 127.0.0.1 port 3332 [tcp/mcs-mailsvr] succeeded!
uid:1:admin:admin
成功识别设备编码.
```



>   注意用户名和密码不能出现`:`

3.   当链接完成后用户即可进行TCP/IP通讯，如果使用的是`nc`命令你只需要直接输入相关内容然后按下回车即可完成消息发送。下面是一个数据发送的案例

```
Connection to 127.0.0.1 port 3332 [tcp/mcs-mailsvr] succeeded!
uid:1:admin:admin
成功识别设备编码.
1
数据已处理.
2
数据已处理.
3
数据已处理.
4
数据已处理.
5
数据已处理.
```

