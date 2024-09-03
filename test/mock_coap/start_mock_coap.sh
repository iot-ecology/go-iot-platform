#!/bin/zsh
docker-compose -f ./docker-compose.yml down
docker rmi mock_coap:latest
docker-compose -f ./docker-compose.yml up -d
echo "coap设备模拟启动"
