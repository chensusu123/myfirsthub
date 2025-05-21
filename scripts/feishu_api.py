import time
import requests


def send_msg(msg):
    count = 0
    url = "http://xalert.ext.17paipai.cn:8099/api/v1/alarms/general/"
    payload = msg
    headers = {
        'Content-Type': "application/json",
        'authorization': "Application 8d923df6-01dd-11e9-ad73-6805ca2f33bc",
    }
    while True:
        try:
            response = requests.request("POST", url, json=payload, headers=headers)
            print(response.status_code, 'status_code')
            if int(response.status_code) == 200:
                break
            else:
                count += 1
            if count >= 10:
                break
        except Exception as e:
            count += 1
            if count >= 10:
                break
            else:
                pass


def start(title, content, receiver):
    msg = {
        'app': 'hot_related_alarms',
        'type': 'hot_related_alarms',
        'host_name': f'\n{title}\n',
        'content': f'{content}',
        'broadcast': '电话语音播报测试大于十五条',  # 语音播报内容, 如不需要删除字段后留空. 最长十五个字.
        'severity': '3',
        # 填谁艾特谁
        'receiver': f'{receiver}'
    }
    send_msg(msg)
