#!/bin/zsh
docker-compose -f ./docker-compose.yml down
docker rmi mock_mqtt:latest
docker-compose -f ./docker-compose.yml up -d
echo "mqtt设备模拟启动"
