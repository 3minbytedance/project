#!/bin/bash

# 要关闭的端口列表
port_list=(8080 4001 4002 4003 4004 4005 4006)

for port in "${port_list[@]}"; do
    # 查找端口上的进程
    pid=$(lsof -t -i :$port)

    if [ -n "$pid" ]; then
        # 终止进程
        echo "Terminating process on port $port (PID: $pid)"
        kill "$pid"
    else
        echo "No process found on port $port"
    fi
done
