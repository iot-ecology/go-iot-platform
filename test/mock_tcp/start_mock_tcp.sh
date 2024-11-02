#!/bin/zsh
docker-compose -f ./docker-compose.yml down
docker rmi mock_tcp:latest
docker-compose -f ./docker-compose.yml up -d
echo "tcp设备模拟启动"
