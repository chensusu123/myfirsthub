#!/bin/bash

# 用法示例：
# ./sql_init_linux.sh root mypass

if [ $# -lt 2 ]; then
    echo "用法: $0 用户名 密码"
    echo "示例: $0 root mypass"
    exit 1
fi

MYSQL_USER="$1"
MYSQL_PWD="$2"
MYSQL_DB="test"
MYSQL_HOST="localhost"
SQL_FILE="init.sql"

if ! command -v mysql &> /dev/null; then
    echo "错误：mysql 命令未找到，请确认已安装 MySQL 客户端。"
    exit 1
fi

mysql -u "$MYSQL_USER" -p"$MYSQL_PWD" -h "$MYSQL_HOST" "$MYSQL_DB" < "$SQL_FILE"

if [ $? -eq 0 ]; then
    echo "执行成功。"
else
    echo "执行失败。"
    exit 1
fi
