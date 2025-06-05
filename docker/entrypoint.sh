#!/bin/sh
set -e

# 替换配置文件里的127.0.0.1为容器服务名
sed -i 's/127.0.0.1:6379/redis/g' /app/conf.d/config.ini
# 如果有 MySQL 同理替换
#sed -i 's/127.0.0.1/mysql/g' /app/conf.d/config.ini

# 执行主程序
exec "$@"
