#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import random
import string
import redis
import sys

# 定义不同环境的配置
ENV_CONFIG = {
    "dev": {
        "redis_host": "10.101.110.231",
        "redis_port": 9004,
        "http_url": "https://test-reg.midudutech.com/user/register/mail"
    },
    "test": {
        "redis_host": "your_test_redis_host",  # 请替换为实际的测试环境 Redis 主机地址
        "redis_port": 6379,  # 请替换为实际的测试环境 Redis 端口
        "http_url": "https://your_test_http_url"  # 请替换为实际的测试环境 HTTP 地址
    }
}

def create_email_accounts(email_prefix, count, env):
    if env == "test":
        print("暂时不支持 test 环境，请使用 dev 环境。")
        return
    if env not in ENV_CONFIG:
        print(f"不支持的环境: {env}，请使用 dev 或 test")
        return

    config = ENV_CONFIG[env]
    # 连接 Redis
    r = redis.Redis(host=config["redis_host"], port=config["redis_port"], decode_responses=True)

    url = config["http_url"]
    headers = {
        'Accept': '*/*',
        'Accept-Encoding': 'gzip, deflate, br',
        'Cache-Control': 'no-cache',
        'Connection': 'keep-alive',
        'User-Agent': 'Make-Maze-Account/1.1.0'
    }

    for _ in range(count):
        # 从 user:id:pool 队列中获取最右侧的数字
        last_id = r.lindex('user:id:pool', -1)
        if last_id is None:
            print("Redis 队列 user:id:pool 为空，无法创建邮箱账号。")
            break

        # 取最后 7 位数，不足 7 位左侧补 0
        email_num_part = str(last_id)[-7:].zfill(7)
        email = f"{email_prefix}{email_num_part}@v8.com"
        password = ''.join(random.choices(string.digits, k=6))
        data = {
            'mail': email,
            'password': password,
            'op_type': 1,
        }

        try:
            response = requests.post(url, headers=headers, data=data)
            response.raise_for_status()  # 检查请求是否成功
            result = response.json()
            if result.get("status") == 200:
                auth_id = result.get("data", {}).get("auth_info", {}).get("auth_id")
                print(f"邮箱账号: {email}, 密码: {password}, auth_id: {auth_id}")
            else:
                print(f"创建邮箱账号 {email} 失败，接口返回非 200 状态码: {result.get('desc', '未知错误')}")
        except requests.RequestException as e:
            print(f"创建邮箱账号 {email} 失败: {e}")
        except ValueError:
            print(f"创建邮箱账号 {email} 失败，无法解析接口返回的 JSON 数据")

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("用法: python make_account.py <邮箱前缀> <数量> <环境（dev/test）>")
        sys.exit(1)

    email_prefix = sys.argv[1]
    count = int(sys.argv[2])
    env = sys.argv[3]

    create_email_accounts(email_prefix, count, env)