import requests


def start(number, content):
    print(f"number:{number} content:{content}")
    url = 'https://xalert-ext.17paipai.cn/api/v1/alarms/call_phone/'
    data = {
        'call_phone': number,  # '18803244613',
        'message': content
    }

    headers = {
        'Content-Type': "application/json",
        'authorization': "Application 8d923df6-01dd-11e9-ad73-6805ca2f33bc",
    }
    res = requests.post(url=url, json=data, headers=headers)
    print(res.json())
