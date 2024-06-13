#!/bin/bash




docker rm -f go_iot
sleep 1
docker rmi -f go_iot:1.0
sleep 1
docker build -t go_iot:1.0 -f Dockerfile .
sleep 1
docker run  -d -p8919:8080  --restart=always --name go_iot   go_iot:1.0