#!/bin/bash

# 启动所有服务
echo "Starting all services..."

# 启动 api 服务
echo "Starting API service..."
bash start.sh --service api &

# 启动 user 服务
echo "Starting User service..."
bash start.sh --service user &

# 启动 comment 服务
echo "Starting Comment service..."
bash start.sh --service comment &

# 启动 relation 服务
echo "Starting Relation service..."
bash start.sh --service relation &

# 启动 message 服务
echo "Starting Message service..."
bash start.sh --service message &

# 启动 favorite 服务
echo "Starting Favorite service..."
bash start.sh --service favorite &

# 启动 favorite 服务
echo "Starting Video service..."
bash start.sh --service video &


echo "All services started."
