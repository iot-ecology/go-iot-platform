#!/bin/bash

nohup ./go-iot -config app-local.yml > output.log 2>&1 &

echo "go-iot程序已启动，输出将被重定向到output.log"