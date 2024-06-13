# Go IoT

单机 100 个客户端 CPU 火焰图

![](out.svg)

## MQTT 客户端管理方案

在物联网项目中关于 MQTT 客户端的管理是一个复杂问题，本仓库仅作为一个案例仓库用来说明相关设计思想。




## 核心设计

### 节点基础配置



```
type NodeInfo struct {
    Host string `json:"host,omitempty" yaml:"host,omitempty"`
    Port int    `json:"port,omitempty" yaml:"port,omitempty"`
    Name string `json:"name,omitempty" yaml:"name,omitempty"`
    Type string `json:"type,omitempty" yaml:"type,omitempty"`
    Size int64  `json:"size,omitempty" yaml:"size,omitempty"` // 最大处理数量
}
```





### Redis 键设计

| Redis 键             | 数据类型 | 明细                                                |
| -------------------- | -------- | --------------------------------------------------- |
| mqtt_config:no       | Hash     | HKEY = MQTT 客户端ID<br />HValue =  MQTT 客户端配置 |
| mqtt_config:use      | Hash     | HKEY = MQTT 客户端ID<br />HValue =  MQTT 客户端配置 |
| node_bind:{节点名称} | List     | 存储节点名称对应的 MQTT 客户端ID                    |
| register:mqtt        | Hash     | HKEY = 节点名称<br />HValue = 节点基础配置          |
| mqtt_size:{节点名称} | String   | 用来记录节点已经创建的 MQTT 客户端数量              |
| beat:mqtt:{节点名称} | string   | 用来设置心跳，过期时间为1200毫秒                    |



### 核心任务1： 创建 MQTT 客户端

1.   选择集群中的最小使用数量的节点，向他发送 HTTP 请求
2.   创建 MQTT 客户端





### 核心任务2： 处理没有使用的 MQTT 客户端配置

使用定时任务的方式从 `mqtt_config:no` 键中获取数据并创建





### 心跳检查

每隔 1000 毫秒将当前节点的名称写入到 redis 中，这个名称的过期时间 1000 毫秒. 




### 负载策略
所有节点每次创建 1 个 MQTT 客户端后会累加 1 , 每次都会选择创建数量最小的节点进行 MQTT 客户端创建

**限制**:

1. 通过配置文件可以约束单个节点最大多少个 MQTT 客户端



### MQTT 客户端被踢出

1.   计数器（`mqtt_size:{节点名称}`）减一
2.   将使用数据移除
     1.   `node_bind:{节点名称}` 移除
     2.   `mqtt_config:use` 移除
     3.   `mqtt_config:no` 加入

### 节点离线

发现节点离线的方式：

1.   基于 Redis 的过期监听
2.   基于程序的定时任务

发现节点离线后：

1.   将离线节点的数据进行移除操作
     1.   计数器（`mqtt_size:{节点名称}`）删除
     2.   `node_bind:{节点名称}` 移除
     3.   `mqtt_config:use` 移除
     4.   `mqtt_config:no` 加入
2.   等待  `mqtt_config:no` 进行任务分发


## 分析工具

```shell
go tool pprof -alloc_space -cum -svg http://127.0.0.1:8080/debug/pprof/heap > heap.svg
go tool pprof -alloc_objects -cum -svg http://127.0.0.1:8080/debug/pprof/heap > obj.svg
```