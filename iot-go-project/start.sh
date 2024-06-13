#!/bin/bash

# 使用nohup命令在后台运行igp程序，忽略挂断信号
nohup ./igp -config app-node1.yml > output.log 2>&1 &

echo "igp程序已启动，输出将被重定向到output.log"