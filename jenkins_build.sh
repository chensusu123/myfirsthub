#!/bin/bash
#$1 jenkins工程工作目录
#$2 jenkins工程工作目录
#$3 tar 包输出名字 要存放在 ${1} 目录下
#$4 项目工程目录名称
#$5 项目可执行文件名称

# 设置环境变量
export GOPROXY=https://goproxy.cn,direct
export GOPRIVATE=gitlab.ifreetalk.com/*
# NOTE:
#   192.168.101.0 gitlab.ifreetalk.com

# 工作目录
WORKSPACE_DIR=${1}
# 打包文件名
RELEASE_NAME=${3}.tar.gz
# 项目工程名
SRC_DIR_NAME=${4}
# 项目可执行文件名称
EXE_NAME=${5}
# 出现错误及时停止脚本.

# set -e

if [ -d ${SRC_DIR_NAME} ]; then
    echo "begin to build" ${SRC_DIR_NAME}
else
    echo "Not ${SRC_DIR_NAME} dir"
    exit 1
fi

cd ${WORKSPACE_DIR}/${SRC_DIR_NAME}

# 获取版本号及tag号
VERSION=$(git log --date=iso --pretty=format:"%h" -1)
if [ $? -ne 0 ]; then
    VERSION="Not a git repo"
fi

TAG=$(git log -1 --pretty=format:"%s")

BUILDTIME=$(date +"%F %T %z")

# 生成版本文件
cat >version.go <<EOF
package main

import(
    "gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
)

func init() {
    fkserver.Version = "$VERSION"
    fkserver.Tag = "$TAG"
    fkserver.BuildTime = "$BUILDTIME"
}
EOF

# 生成程序
OS=$(go env GOOS)
echo $OS
if [ $OS != 'linux' ]; then
    # 如果不是 linxu系统. 先生成本地程序. 导出依赖配置.然后清理
    go build -o $EXE_NAME
    if [ $? -ne 0 ]; then
        rm -f version.go
        exit 1
    fi
    ./${SRC_DIR_NAME}/$EXE_NAME -e
    go clean
fi

# 生成linux程序
GOOS=linux GOARCH=amd64 go build -o $EXE_NAME
if [ $OS == 'linux' ]; then
    # linux 导出依赖配置
    #./$EXE_NAME -e
    echo "linux build"
fi

# 创建配置文件目录
if [ ! -d conf.d ]; then
    mkdir conf.d
fi

# 打包
tar -zcvf ${WORKSPACE_DIR}/${RELEASE_NAME} ./${EXE_NAME} ./extern.conf conf.d

# 删除临时文件
rm -f version.go
rm -f ./${EXE_NAME}
rm -f ./extern.conf
