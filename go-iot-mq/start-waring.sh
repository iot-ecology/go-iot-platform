#!/bin/bash
rm -rf output.log
nohup ./gim -config app-local-waring_handler.yml > output.log 2>&1 &

echo "gim程序已启动，输出将被重定向到output.log"