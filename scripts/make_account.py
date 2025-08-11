#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import random
import string
import redis
import sys

# 定义不同环境的配置
ENV_CONFIG = {
    "test": {
        "redis_host": "10.101.110.239",
        "redis_port": 65001,
        "http_url": "https://test-reg.midudutech.com/user/register/mail",
        "gm_url_template": "http://test-gm.midudutech.com/s5/%s/generateUser?AuthId=%s"
    },
    "play": {
        "redis_host": "10.101.110.239",
        "redis_port": 65002,
        "http_url": "https://play-reg.midudutech.com/user/register/mail",
        "gm_url_template": "http://play-gm.midudutech.com/s4/%s/generateUser?AuthId=%s"
    }
}

def call_gm_api(auth_id, env):
    gm_url = ENV_CONFIG[env]["gm_url_template"] % (auth_id, auth_id)
    try:
        response = requests.get(gm_url)
        response.raise_for_status()
        result = response.json()
        error_code = result.get("errorCode")
        error_msg = result.get("errorMsg")
        user_id = result.get("userId")

        if error_code == 0:
            return True, user_id, error_msg
        else:
            return False, user_id, error_msg
    except requests.RequestException as e:
        return False, 0, f"调用 GM 接口失败: {e}"
    except ValueError:
        return False, 0, "无法解析 GM 接口返回的 JSON 数据"

def create_email_accounts(email_prefix, count, env):
    if env not in ENV_CONFIG:
        print(f"不支持的环境: {env}，请使用 test 或 play")
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

    # 输出表头
    print("邮箱账号,密码,AuthID,角色ID")
    import hashlib  # 导入 hashlib 库

    # todo 测试count是否预处理读取
    i = 0 
    failCount = 0
    while i < count:
        if failCount > 5:
            print("连续5次失败，退出")
            break
        # 从 user:id:pool 队列中获取最右侧的数字
        tmpId = r.get('s:0:account:id:pool')
        if tmpId is None:
            print("账号ID分配失败，无法创建邮箱账号。")
            break

        last_id = int(tmpId) + 1

        # 取最后 7 位数，不足 7 位左侧补 0
        email_num_part = str(last_id)[-7:].zfill(7)
        email = f"{email_prefix}{email_num_part}@v8.com"
        # 生成首位不为 0 的 6 位数字密码
        first_digit = random.choice(string.digits[1:])
        rest_digits = ''.join(random.choices(string.digits, k=5))
        password = first_digit + rest_digits

        # 计算密码的 MD5 值
        md5_password = hashlib.md5(password.encode('utf-8')).hexdigest()

        data = {
            'mail': email,
            'password': md5_password,  # 使用 MD5 密码
            'op_type': 1,
        }

        try:
            response = requests.post(url, headers=headers, data=data)
            response.raise_for_status()  # 检查请求是否成功
            result = response.json()
            if result.get("status") == 200:
                auth_id = result.get("data", {}).get("auth_info", {}).get("auth_id")
                success, user_id, error_msg = call_gm_api(auth_id, env)
                if success:
                    print(f"\"{email}\",\"{password}\",\"{auth_id}\",\"{user_id}\"")
                    failCount = 0
                else:
                    print(f"创建邮箱账号 {email} 后，调用 GM 接口失败: {error_msg}") 
            elif result.get("status") == 500:
                count += 1
                r.incrby('s:0:account:id:pool', 1)
                failCount += 1
            else:
                print(f"创建邮箱账号 {email} 失败，接口返回非 200 状态码: {result.get('desc', '未知错误')}")
        except requests.RequestException as e:
            print(f"创建邮箱账号 {email} 失败: {e}")
        except ValueError:
            print(f"创建邮箱账号 {email} 失败，无法解析接口返回的 JSON 数据")
        
        i+=1

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("用法: python make_account.py <邮箱前缀> <数量> <test/play）>")
        sys.exit(1)

    email_prefix = sys.argv[1]
    count = int(sys.argv[2])
    env = sys.argv[3]

    create_email_accounts(email_prefix, count, env)