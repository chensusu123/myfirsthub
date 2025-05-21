#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import datetime
import os
import shutil
import subprocess
import sys
import time
from pathlib import Path

import requests

current_path = os.path.split(os.path.realpath(__file__))[0].replace('\\', '/')
java = os.path.join(current_path, "jdk-19.0.2/bin/java")
jenkins_jar = os.path.join(current_path, "jenkins-cli.jar")


def execute_command(command):
    cmd = f"{java} -jar {jenkins_jar} -s http://192.168.100.75:8080/ -auth zhangzhen:xiaozhen@123 {command}"
    print(f"execute_command >>> command:{command}")
    os.system(cmd)


def get_job_build_status(job_path):
    url = f"http://192.168.100.75:54139/{job_path}"
    print(f"get_job_build_status >>> url:{url}")
    res = requests.get(url)
    if res.status_code == 200:
        data = res.json()
        print(f"get_job_build_status >>> data:{data}")
        return data["status"]
    else:
        raise Exception(f"request {url} error")


def get_job_build_log(job_path, build_number):
    url = f"http://192.168.100.75:54141/get_log_content/{job_path}/{build_number}"
    print(f"get_job_build_log >>> url:{url}")
    res = requests.get(url)
    if res.status_code == 200:
        data = res.content
        return str(data, "utf-8")
    else:
        raise Exception(f"request {url} error")


def get_job_build_log_url(job_path, build_number):
    return f"http://192.168.100.75:54141/get_log_content/{job_path}/{build_number}"


def get_job_build_config(job_path):
    url = f"http://192.168.100.75:54140/{job_path}"
    print(f"get_job_build_config >>> url:{url}")
    res = requests.get(url)
    if res.status_code == 200:
        data = res.json()
        print(f"get_job_build_config >>> data:{data}")
        return data
    else:
        raise Exception(f"request {url} error")


def start_build_job(job_path, client_branch, unity_branch):
    print(f"start_build_job >>> job_path:{job_path}")
    execute_command(f"build {job_path} -p client_branch_name={client_branch} -p unity_branch_name={unity_branch}")



