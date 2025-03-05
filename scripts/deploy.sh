#!/bin/bash

# 部署脚本
set -e

echo "开始部署 Govue3 项目..."

# 1. 构建项目
make build

# 2. 构建 Docker 镜像
docker build -t govue3:latest .

# 3. 运行容器
docker run -d \
  --name govue3 \
  -p 8080:8080 \
  -v $(pwd)/config:/app/config \
  govue3:latest

echo "部署完成！"