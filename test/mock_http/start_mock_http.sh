#!/bin/zsh
docker-compose -f ./docker-compose.yml down
docker rmi mock_http:latest
docker-compose -f ./docker-compose.yml up -d
echo "http设备模拟启动"
