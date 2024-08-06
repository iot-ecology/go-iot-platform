---
publishDate: 2024-08-06T00:00:00Z
author: Zen HuiFer
title: 多协议支持
excerpt: 在 Go IoT 开发平台中支持使用WebSocket、MQTT、TCP/IP、COAP协议进行数据传输
image: "~/assets/images/communication-1439132_1280.jpg"

category: 教程
tags:
  - Go
  - 物联网
  - 开发平台
  - 应用
  - 实战
---

在Go IoT 开发平台中关于 WebSocket、MQTT、TCP/IP、COAP 端口默认使用情况如下



| 协议 | 端口 |
| ---- | ---- |
| WebSocket     | 13332 |
| MQTT     | 1883 8083 8084 8883 18083 |
| COAP     | 5683 |
| TCP/IP     | 3332 |



如果你需要配置Nginx可以考虑使用如下内容

-   TCP/IP

```
stream{
    upstream tcpserver {
        server 0.0.0.0:3332;
    }
    server {
        listen 22122;
        proxy_pass tcpserver;
    }
}
```

-   WebSocket Nginx配置

```

http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  65;    
   
    map $http_upgrade $connection_upgrade {
        default upgrade;
        '' close;
    }
    upstream sre_backend {
        server 127.0.0.1:13332;
    }

    server {
        listen       80;
        server_name  localhost;


        location / {
            proxy_pass              http://sre_backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "$connection_upgrade";
        }

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    }


  

}

```

-   COAP

```
stream{

    upstream coap_server {
        server 0.0.0.0:5683;
    }
    server {
        listen 15683 udp;
        proxy_pass coap_server;
    }
}


```

