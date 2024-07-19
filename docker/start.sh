#!/bin/zsh
 docker-compose -f env/base-env-docker-compose.yml up  -f app/docker-compose.yml -d
 echo "执行完毕,项目启动中..."
