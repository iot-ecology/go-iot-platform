# rabbit 消费者部署

## 部署步骤

### 1. 编译

编译项目以生成可执行文件 `gim`。

**操作步骤**：

- 进入项目目录：

  ```
  cd go-iot-mq
  ```

- 下载依赖：

  ```
  go mod tidy
  ```

- 编译项目：

  ```
  go build -o gim
  ```

**结果**：编译完成后，在项目目录下会生成一个名为 `gim` 的可执行文件。

### 2. 修改配置文件

创建或修改配置文件，以配置 MQTT 客户端管理项目所需的 Redis、消息队列、MongoDB和 InfluxDB 信息。

**配置示例**：

```
node_info:
  host: 127.0.0.1
  port: 29001
  name: mq1
  type: calc_queue # pre_handler、 waring_handler、 calc_queue、waring_delay_handler
redis_config:
  host: 127.0.0.1
  port: 6379
  db: 10
  password: eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81

mq_config:
  host: 127.0.0.1
  port: 5672
  username: guest
  password: guest
influx_config:
  host: 127.0.0.1
  port: 8086
  token: i6XHSnNXeUoU3GoFXMm4qqrrgt69JKvQLqm0FCtnYG-rjb-nkDcry0pdwv4fpcXsSwi-mTGmAUTygkJtR-6CWA==
  org: myorg
  bucket: buc
mongo_config:
  host: 127.0.0.1
  port: 27017
  username: admin
  password: admin
  db: iot
  collection: calc
  waring_collection: waring
  script_waring_collection: script_waring
```

**注意事项**：

- 请根据实际环境修改配置文件中的参数，如数据库密码、Redis密码等。
- 确保所有服务（如Redis、RabbitMQ、InfluxDB、MongoDB等）已正确安装并运行。

### 3. 启动服务

使用以下命令启动 MQTT 客户端管理项目。

**启动命令**：

```
./gim -config app-local-calc.yml
./gim -config app-local-pre_handler.yml
./gim -config app-local-waring_handler.yml
./gim -config app-local-wd.yml
```

**说明**：

- 确保在执行启动命令前，所有依赖服务均已启动并可访问。

## 注意事项

- 在部署过程中，确保防火墙规则允许相应的端口访问。
- 定期备份配置文件，以防数据丢失。
- 在生产环境中，建议使用更安全的密码，并限制对配置文件的访问权限。