#!/bin/bash

# 要关闭的端口列表
port_list=(8080 4001 4002 4003 4004 4005 4006)

for port in "${port_list[@]}"; do
    # 查找端口上的进程
    pids=$(lsof -t -i :$port)

    if [ -n "$pids" ]; then
        # 终止进程
        echo "Terminating processes on port $port (PIDs: $pids)"
        kill -9 $pids
    else
        echo "No process found on port $port"
    fi
done
