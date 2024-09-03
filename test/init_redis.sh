#!/bin/bash

# Redis服务器地址
REDIS_HOST="localhost"
# Redis服务器端口
REDIS_PORT="6379"
# Redis密码
REDIS_PASSWORD="eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"

# 连接到Redis服务器并进行身份验证
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD <<EOF
select 11
HSET auth:coap 1234567890 "{\"username\":\"admin\",\"password\":\"admin\",\"device_id\":\"1234567890\"}"
HSET auth:http aaa "{\"username\":\"admin\",\"password\":\"admin\",\"device_id\":\"aaa\"}"
HSET auth:tcp 1 "{\"username\":\"admin\",\"password\":\"admin\",\"device_id\":\"1\"}"
HSET auth:ws 321 "{\"username\":\"admin\",\"password\":\"admin\",\"device_id\":\"321\"}"

EOF

#