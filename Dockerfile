# 使用多阶段构建
FROM golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache git openssh

# 添加认证文件
COPY .netrc /root/.netrc
RUN chmod 600 /root/.netrc

ENV GOPRIVATE=gitlab.ifreetalk.com/*
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o app ./servers/maze_main_server/main.go

  # 最小镜像运行阶段
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .
COPY servers/maze_main_server/conf.d/ ./conf.d/


CMD ["./app"]