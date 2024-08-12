---
publishDate: 2024-07-27T00:00:00Z
author: Zen HuiFer
title: MQTT 实战
excerpt: Docker快速启动一个MQTT服务，并使用MQTT客户端可视化工具进行测试 
image: "~/assets/images/linux-151619_1280.png"
category: 教程
tags:
- Go
- 物联网
- 开发平台
- 设计
- 架构
---

# MQTT 实战

本章将进行MQTT的实际操作，主要包含如下内容：

1. 如何用Docker快速启动一个MQTT服务
2. EMQX的基础使用
3. MQTT客户端可视化工具使用
4. MQTT主题相关知识

## 搭建 MQTT 服务

在本节中，我们将探讨如何使用Docker和Docker Compose技术快速搭建MQTT服务。Docker提供了一种轻量级、可移植的容器化解决方案，而Docker
Compose则允许我们通过一个简单的配置文件管理多容器Docker应用程序。

### 使用Docker搭建MQTT服务

1. **获取Docker镜像** 要开始使用MQTT服务，首先需要获取MQTT的Docker镜像。Emqx是一个流行的开源MQTT代理程序，我们可以使用以下命令来获取其Docker镜像：

   ```
   docker pull emqx/emqx:5.6.1
   ```

2. **启动Docker容器** 获取镜像后，可以通过以下命令启动一个Emqx MQTT服务的Docker容器：

   ```
   docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx:5.6.1
   ```

这个命令将启动一个名为emqx的容器实例，并将Emqx MQTT服务的不同端口映射到宿主机的相应端口上。

### Docker Compose 启动

为了更方便地管理MQTT服务，特别是当需要部署集群时，可以使用Docker Compose。首先，需要编写一个docker-compose.yml文件，如下所示：

```Python
version: '3'

services:
  emqx1:
    image: emqx:5.4.1
    container_name: emqx1
    environment:
      - "EMQX_NODE_NAME=emqx@node1.emqx.io"
      - "EMQX_CLUSTER__DISCOVERY_STRATEGY=static"
      - "EMQX_CLUSTER__STATIC__SEEDS=[emqx@node1.emqx.io,emqx@node2.emqx.io]"
    healthcheck:
      test: ["CMD", "/opt/emqx/bin/emqx ctl", "status"]
      interval: 5s
      timeout: 25s
      retries: 5
    networks:
      emqx-bridge:
        aliases:
          - node1.emqx.io
    ports:
      - 1883:1883
      - 8083:8083
      - 8084:8084
      - 8883:8883
      - 18083:18083

  emqx2:
    image: emqx:5.4.1
    container_name: emqx2
    environment:
      - "EMQX_NODE_NAME=emqx@node2.emqx.io"
      - "EMQX_CLUSTER__DISCOVERY_STRATEGY=static"
      - "EMQX_CLUSTER__STATIC__SEEDS=[emqx@node1.emqx.io,emqx@node2.emqx.io]"
    healthcheck:
      test: ["CMD", "/opt/emqx/bin/emqx ctl", "status"]
      interval: 5s
      timeout: 25s
      retries: 5
    networks:
      emqx-bridge:
        aliases:
          - node2.emqx.io

networks:
  emqx-bridge:
    driver: bridge
```

这个配置文件定义了两个Emqx MQTT服务容器emqx1和emqx2，它们将通过自定义网络emqx-bridge进行通信，并设置了健康检查以确保服务的稳定性。

在编写完docker-compose.yml文件后，可以在文件所在目录执行以下命令来启动服务：

```
docker-compose up -d
```

这个命令将根据`docker-compose.yml`文件中的配置启动并运行MQTT服务。

### 登录及账号密码修改

在完成MQTT服务的部署后，通常需要进行一些基本的设置，包括登录、修改密码以及语言设置等。以下是对部署后的MQTT服务进行操作的步骤：

打开浏览器，输入 http://localhost:18083 访问MQTT服务的管理界面。此时，你将看到一个登录界面，使用默认账号 admin
和默认密码public登录。

![img](03-MQTT实战/1.jpg)

登录后，系统会引导你进入修改密码的界面。根据安全需要，修改默认密码以保证MQTT服务的安全性。

![img](03-MQTT实战/2.jpg)

**完成密码修改** 成功修改密码后，你将进入EMQX的首页，这里是管理MQTT服务的核心界面。

![img](03-MQTT实战/3.png)

### 简体中文设置

如果界面是英文状态，可以通过设置更改为简体中文。英文界面如图所示

![img](03-MQTT实战/4.png)

单击界面右上角的齿轮按钮（设置按钮）。具体位置如图所示

![img](03-MQTT实战/5.png)

在弹出的设置界面中，找到语言（Language）选项。

![img](03-MQTT实战/6.png)

点击Language下方的English，从下拉框中选择“简体中文”。选择后点击“Apply”应用更改。

![img](03-MQTT实战/7.png)

在下拉框中选择简体中文，然后鼠标单击Apply。此时界面切换简体中文结束，具体效果如图所示

![img](03-MQTT实战/8.png)

### 新建用户

在MQTT服务中，新建用户账号是一个基本的管理功能，尽管新创建的用户可能没有办法进行细粒度的权限控制，但新建用户账号仍然非常有用。以下是新建用户账号的步骤：

1. 鼠标左键单击菜单系统设置 -> 用户，此时界面如图所示

![img](03-MQTT实战/9.png)

1. 在这个用户界面中鼠标单击创建按钮，会弹出输入窗口界面如图所示

![img](03-MQTT实战/10.png)

1. 按照你的需要输入用户名、备注、密码三个数据后鼠标单击创建按钮完成创建。完成后用户界面如图所示

![img](03-MQTT实战/11.png)

**注意事项**

- **权限管理**：当前系统可能不支持对新用户进行详细的权限设置，新用户可能拥有与管理员相同的操作权限。这一点需要根据实际使用场景和安全需求来考虑。
- **安全性**：为了系统的安全性，建议定期更新用户密码，并限制对敏感操作的访问。
- **用户教育**：如果系统中有多个用户，确保每个用户都了解自己的权限范围和操作限制。

### 端口使用说明

默认情况下会使用4个端口，具体使用请查看下表

| 端口号  | 用途            | 安全性 |
|------|---------------|-----|
| 1883 | MQTT协议的非加密通信  | 非安全 |
| 8883 | MQTT协议的加密通信   | 安全  |
| 8083 | WebSocket通信   | 非安全 |
| 8084 | WebSocket安全通信 | 安全  |

更多 EMQX 使用说明请查阅：https://www.emqx.io/docs/zh/latest/ 文档

## 模拟消息发送和接收

为了模拟消息发送和接收，可以使用MQTTX，这是一个流行的MQTT客户端工具，它可以帮助我们快速建立连接、发布消息和订阅主题。

要使用MQTTX，请访问https://mqttx.app/zh下载并安装适用于您操作系统的版本。

### 创建MQTT客户端

启动MQTTX后，将看到一个用户友好的界面，具体如图所示

![img](03-MQTT实战/12.png)

接下来我们进行发布者创建，鼠标单击MQTTX软件左侧的加号按钮，此时界面如图所示

![img](03-MQTT实战/13.png)

在图中需要输入的内容如下

- **Name**: 客户端的名称，仅用于MQTTX界面显示。
- **Client ID**: 客户端ID，在一个MQTT服务上必须是全局唯一的。如果与现有客户端ID相同，新建的客户端将替换旧的客户端。
- **Host**: MQTT服务的IP地址或主机名。
- **Port**: MQTT服务的端口号，例如默认的1883（非加密）或8883（加密）。
- **Username**: MQTT服务的用户名，如果有的话。
- **Password**: MQTT服务的密码，如果有的话。

按照本地环境的配置输入相关信息后内容如图所示。

![img](03-MQTT实战/14.png)

### 发送消息和接收消息

打开发布者客户端，界面如图所示

![img](03-MQTT实战/15.png)

找到Connect按钮鼠标单击即可与MQTT服务建立链接，链接后界面如图所示

![img](03-MQTT实战/16.png)

鼠标单击 New Subscription 按钮添加一个订阅的主题，鼠标单击后效果如图所示。

![img](03-MQTT实战/17.png)

在图中需要修改Tpoic的数据内容，修改为test_topic/1，修改完成后点击Confirm按钮完成订阅行为。此时界面如图所示

![img](03-MQTT实战/18.png)

接下来进行消息发送的模拟，发送内容请参考下图

![img](03-MQTT实战/19.png)

按照个人需求向其中写入相关数据，本例发送主题为test_topic/1 ，发送内容为测试数据，填写后界面如图所示。

![img](03-MQTT实战/20.png)

点击右下角绿色发送按钮即可将消息发送。发送按钮点击以后界面如图所示。

![img](03-MQTT实战/21.png)

图中内容解释：

1. 右侧绿色框内是这个客户端发送出去的数据
2. 左侧灰色内容是这个客户端订阅的主题收到的数据

### 大小写敏感

在MQTT协议中，主题的命名是区分大小写的。这意味着，当你向 Test_topic/1 发送消息时，它与发送到 test_topic/1被视为两个独立的主题。如果客户端订阅了 test_topic/1，它将无法接收到发送至 Test_topic/1的消息。主题的这种大小写敏感特性要求在发布和订阅消息时必须严格匹配主题名称。具体测试案例如图所示

![img](03-MQTT实战/22.png)

### 通配符订阅

在MQTT（Message Queuing Telemetry
Transport）协议中，通配符订阅是一种非常有用的功能，它允许客户端订阅一个主题模式，而不是单个具体的主题。这样，客户端可以接收到匹配该模式的所有主题的消息。以下是对MQTT通配符订阅的详细说明：

**通配符订阅规则**

- **单层通配符 +**：这个通配符代表主题中的单个层级。当使用 `+` 通配符时，它将匹配任何单个层级的主题。
- **多层通配符 #**：这个通配符代表主题中从该位置起到末尾的所有层级。使用 `#` 通配符时，它将匹配任意长度的主题。

**通配符的使用限制**

- **仅用于订阅**：MQTT协议规定，通配符只能用于订阅操作，不能用于发布消息。

**示例场景**

假设我们有以下一系列主题：

```
t/1/a
t/1/b
t/2/a
t/2/b
```

如果客户端想要订阅所有以 `t/` 开头，紧接着一个任意单个层级（例如数字），然后是 `/a` 结尾的主题，可以使用单层通配符 `+`
来订阅。具体的语法表达如下：

- **订阅主题**：t/+/a

- 将匹配以下主题

    - t/1/a
    - t/2/a
    - 以及任何符合 t/任意数字/a 模式的其他主题

下面开始进行实际测试，测试案例如图所示

![img](03-MQTT实战/23.png)

如果想要订阅`t`下面的所有主题此时可以使用 # 通配符，具体的语法表达：t/#

测试案例如图所示

![img](03-MQTT实战/24.png)

### 系统主题

MQTT系统主题以 `$SYS/` 开头，它们在MQTT协议中用于提供关于MQTT服务器自身的信息，包括运行状态、统计数据以及客户端的上线和下线事件等。这些系统主题对于监控和管理MQTT服务非常有用。

系统主题可以被订阅以接收有关MQTT服务器状态的更新。例如，如果你想要监控EMQX集群中所有节点的状态，你可以订阅 $SYS/brokers主题。当任何节点的状态发生变化时，你将收到更新消息。

常见的MQTT 系统主题见表

| 主题                                   | 说明          |
|--------------------------------------|-------------|
| $SYS/brokers                         | EMQX 集群节点列表 |
| $SYS/brokers/emqx@127.0.0.1/version  | EMQX 版本     |
| $SYS/brokers/emqx@127.0.0.1/uptime   | EMQX 运行时间   |
| $SYS/brokers/emqx@127.0.0.1/datetime | EMQX 系统时间   |
| $SYS/brokers/emqx@127.0.0.1/sysdescr | EMQX 系统信息   |

### 主题的命名规则

在MQTT中，主题的命名遵循一些特定的规则和最佳实践，以确保主题的组织性和避免潜在的冲突：

1. **避免使用 `$SYS` 作为主题的开头**：$SYS 是保留的前缀，用于系统主题，这些主题提供了关于MQTT服务器自身状态的信息。
2. **使用 `/` 作为主题的开头**：通常建议主题以 / 开头，这有助于清晰地标识主题的开始，并与层级结构中的其他部分区分开来。
3. **避免使用数字作为主题的开头**：虽然MQTT协议没有明确禁止，但建议主题开头使用字母或有意义的标识符，以提高主题的可读性和管理性。
4. **避免使用 $share 作为主题开头**：$share 是一个特殊的前缀，用于实现共享订阅，允许多个客户端共享对同一主题的订阅。如果不需要共享订阅，应避免使用此前缀。

### 示例

以下是一些符合上述规则的主题命名示例：

- 正确：home/living_room/temperature
- 正确：sensor/ambient_light
- 错误：$SYS/our_temperature（应保留给系统主题）
- 错误：2nd_floor/temperature（以数字开头）
- 错误：$share/our_group_topic（除非用于共享订阅）
