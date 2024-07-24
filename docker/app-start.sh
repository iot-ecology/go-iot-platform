#!/bin/zsh
docker-compose -f ./app/docker-compose.yml down
docker rmi go-iot-project:latest go-iot-admin-vue:latest go-iot-mq:latest go-iot-mqtt:latest
docker-compose -f ./app/docker-compose.yml up -d
echo "项目启动中..."
