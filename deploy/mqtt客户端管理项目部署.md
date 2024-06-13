# mqtt客户端管理项目部署

## 部署步骤

### 1. 编译

编译项目以生成可执行文件 `go-iot`。

**操作步骤**：

- 进入项目目录：

  ```
  cd go-iot
  ```

- 下载依赖：

  ```
  go mod tidy
  ```

- 编译项目：

  ```
  go build -o go-iot
  ```

**结果**：编译完成后，在项目目录下会生成一个名为 `go-iot` 的可执行文件。

### 2. 修改配置文件

创建或修改配置文件 `app-node1.yml`，以配置 MQTT 客户端管理项目所需的 Redis 和消息队列信息。

**配置示例**：

```
node_info:
  host: 127.0.0.1
  port: 8081
  name: m1
  type: mqtt
  size: 1000
redis_config:
  host: 127.0.0.1
  port: 6379
  db: 0
  password: eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81

mq_config:
  host: 127.0.0.1
  port: 5672
  username: guest
  password: guest
```

**注意事项**：

- 请根据实际环境修改配置文件中的参数，如 Redis 密码等。
- 确保 Redis 和 RabbitMQ 服务已正确安装并运行。

### 3. 启动服务

使用以下命令启动 MQTT 客户端管理项目。

**启动命令**：

```
./go-iot -config app-node1.yml
```

**说明**：

- 确保在执行启动命令前，所有依赖服务均已启动并可访问。

## 注意事项

- 在部署过程中，确保防火墙规则允许相应的端口访问。
- 定期备份配置文件，以防数据丢失。
- 在生产环境中，建议使用更安全的密码，并限制对配置文件的访问权限。