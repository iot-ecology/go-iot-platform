#!/bin/bash

# 使用nohup命令在后台运行mqtt-mock程序，忽略挂断信号
nohup ./mqtt-mock  > output.log 2>&1 &

echo "mqtt-mock程序已启动，输出将被重定向到output.log"