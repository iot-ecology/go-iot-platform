#!/bin/zsh
docker-compose -f ./docker-compose.yml down
docker rmi mock_ws:latest
docker-compose -f ./docker-compose.yml up -d
echo "ws设备模拟启动"
