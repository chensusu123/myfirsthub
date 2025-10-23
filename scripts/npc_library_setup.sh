#!/bin/bash

# 图鉴系统初始化脚本
# 用于设置数据库、Redis缓存和系统配置

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    # 检查MySQL客户端
    if ! command -v mysql &> /dev/null; then
        log_error "MySQL客户端未安装"
        exit 1
    fi
    
    # 检查Redis客户端
    if ! command -v redis-cli &> /dev/null; then
        log_error "Redis客户端未安装"
        exit 1
    fi
    
    log_info "依赖检查完成"
}

# 数据库配置
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-"root"}
DB_PASSWORD=${DB_PASSWORD:-""}
DB_NAME=${DB_NAME:-"maze_game"}

# Redis配置
REDIS_HOST=${REDIS_HOST:-"localhost"}
REDIS_PORT=${REDIS_PORT:-6379}
REDIS_PASSWORD=${REDIS_PASSWORD:-""}

# 初始化数据库
init_database() {
    log_info "初始化数据库..."
    
    # 创建数据库（如果不存在）
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    
    # 创建分片表
    for i in {0..63}; do
        log_info "创建分片表 npc_library_$i"
        
        # 替换SQL模板中的占位符
        sed "s/%d/$i/g" scripts/npc_library_init.sql | mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME"
    done
    
    log_info "数据库初始化完成"
}

# 初始化Redis
init_redis() {
    log_info "初始化Redis缓存..."
    
    # 测试Redis连接
    if [ -n "$REDIS_PASSWORD" ]; then
        redis-cli -h"$REDIS_HOST" -p"$REDIS_PORT" -a"$REDIS_PASSWORD" ping
    else
        redis-cli -h"$REDIS_HOST" -p"$REDIS_PORT" ping
    fi
    
    # 设置Redis配置
    if [ -n "$REDIS_PASSWORD" ]; then
        redis-cli -h"$REDIS_HOST" -p"$REDIS_PORT" -a"$REDIS_PASSWORD" config set maxmemory-policy allkeys-lru
    else
        redis-cli -h"$REDIS_HOST" -p"$REDIS_PORT" config set maxmemory-policy allkeys-lru
    fi
    
    log_info "Redis初始化完成"
}

# 创建配置文件
create_config() {
    log_info "创建配置文件..."
    
    # 复制配置文件到项目根目录
    if [ ! -f "config/npc_library_config.yaml" ]; then
        log_error "配置文件不存在: config/npc_library_config.yaml"
        exit 1
    fi
    
    # 创建日志目录
    mkdir -p logs
    
    log_info "配置文件创建完成"
}

# 验证安装
verify_installation() {
    log_info "验证安装..."
    
    # 检查数据库表
    table_count=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES LIKE 'npc_library_%';" | wc -l)
    if [ $table_count -lt 64 ]; then
        log_error "数据库表创建不完整，期望64个表，实际$table_count个"
        exit 1
    fi
    
    # 检查Redis连接
    if [ -n "$REDIS_PASSWORD" ]; then
        redis-cli -h"$REDIS_HOST" -p"$REDIS_PORT" -a"$REDIS_PASSWORD" ping > /dev/null
    else
        redis-cli -h"$REDIS_HOST" -p"$REDIS_PORT" ping > /dev/null
    fi
    
    if [ $? -ne 0 ]; then
        log_error "Redis连接失败"
        exit 1
    fi
    
    log_info "安装验证完成"
}

# 显示使用说明
show_usage() {
    log_info "图鉴系统初始化完成！"
    echo ""
    echo "使用说明："
    echo "1. 数据库已创建64个分片表：npc_library_0 到 npc_library_63"
    echo "2. Redis缓存已配置"
    echo "3. 配置文件已创建：config/npc_library_config.yaml"
    echo "4. 日志目录已创建：logs/"
    echo ""
    echo "环境变量："
    echo "  DB_HOST: 数据库主机 (默认: localhost)"
    echo "  DB_PORT: 数据库端口 (默认: 3306)"
    echo "  DB_USER: 数据库用户 (默认: root)"
    echo "  DB_PASSWORD: 数据库密码 (默认: 空)"
    echo "  DB_NAME: 数据库名称 (默认: maze_game)"
    echo "  REDIS_HOST: Redis主机 (默认: localhost)"
    echo "  REDIS_PORT: Redis端口 (默认: 6379)"
    echo "  REDIS_PASSWORD: Redis密码 (默认: 空)"
    echo ""
    echo "启动服务："
    echo "  go run servers/maze_main_server/main.go"
}

# 主函数
main() {
    log_info "开始初始化图鉴系统..."
    
    check_dependencies
    init_database
    init_redis
    create_config
    verify_installation
    show_usage
    
    log_info "图鉴系统初始化完成！"
}

# 错误处理
trap 'log_error "初始化失败，请检查错误信息"; exit 1' ERR

# 执行主函数
main "$@"


