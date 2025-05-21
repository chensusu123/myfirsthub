#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import logging
import os
import requests
import jenkins_api
from flask import Flask, request, jsonify

app = Flask(__name__)


def handle_web_hook():
    logging.info("handle_web_hook")
    data = request.get_json()
    branch_name = read_branch_name(data)
    if len(branch_name) == 0:
        return 'branch_name is empty', 201
    user_name = read_user_name(data)
    if len(user_name) == 0:
        return 'user_name is empty', 201
    editor_tag, resource_build = read_editor_tag_and_env(branch_name)
    if len(editor_tag) == 0 and len(resource_build) == 0:
        return f'build param is empty editor_tag:{editor_tag} resource_build:{resource_build}', 201
    logging.info(f"user_name:{user_name} branch_name:{branch_name} editor_tag:{editor_tag} resource_build:{resource_build}")
    execute_build_atlas(branch_name, user_name, editor_tag, resource_build)
    return 'success', 200


def read_branch_name(data):
    if data and 'ref' in data:
        ref = data['ref']
        if 'refs/heads/' in ref:
            branch_name = ref.replace('refs/heads/', '').replace(" ", "")
            return branch_name
    return ""


def read_user_name(data):
    if data and 'user_name' in data:
        user_name = data['user_name']
        user_name.replace(" ", "")
        return user_name
    return ""


def read_editor_tag_and_env(branch_name):
    file_path = os.getcwd() + "/BuildProjectEditor.cs"
    if os.path.exists(file_path):
        os.remove(file_path)
    download_editor_tag("UnityProject/Assets/Editor/BuildProjectEditor/BuildProjectEditor.cs", file_path, branch_name)
    if not os.path.exists(file_path):
        return ""
    with open(file_path, 'r', encoding='utf-8') as f:
        res = f.read()  #读文件
        f.close()
    if '/***start***/"' in res and '";/***end***/' in res:
        tag_begin = res.find('/***start***/"')
        tag_end = res.find('";/***end***/')
        tag = res[tag_begin+14:tag_end]
    else:
        tag = ""
    if '/***start_resource_build***/"' in res and '";/***end_resource_build***/' in res:
        resource_build_begin = res.find('/***start_resource_build***/"')
        resource_build_end = res.find('";/***end_resource_build***/')
        resource_build = res[resource_build_begin + 29:resource_build_end]
    else:
        resource_build = ""
    return tag, resource_build


def download_editor_tag(file_path, save_as, branch_name):
    url = f"http://192.168.8.11/api/v4/projects/637/repository/files/{file_path.replace('/', '%2F')}/raw?ref={branch_name}"
    logging.info(url)
    headers = {"PRIVATE-TOKEN": "glpat-43biYN7KdEJC1iuGDzdN"}
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        with open(save_as, "wb") as file:
            file.write(response.content)
        logging.info(f"文件已保存为 {save_as}")
    else:
        logging.info(f"请求失败，状态码：{response.status_code}, 错误信息：{response.text}")


def execute_build_atlas(branch_name, user_name, editor_tag, resource_build):
    if editor_tag == "b55dffd24a16":
        jenkins_api.execute_command(f'build resource/online/build_resource -p pp_world_git={branch_name} -p build_type=build_atlas')
    elif resource_build == "publish":
        # 基于线上分支
        jenkins_api.execute_command(f'build resource/publish/build_resource -p pp_world_git={branch_name} -p build_type=build_atlas')
    else:
        # 基于人偶分支
        jenkins_api.execute_command(f'build resource/develop/build_resource -p pp_world_git={branch_name} -p build_type=build_atlas')


def all_url_rule():
    app.add_url_rule("/share_hook", "handle_web_hook", handle_web_hook, methods=["GET", "POST"])


if __name__ == '__main__':
    all_url_rule()
    logging.basicConfig(filename='run.log', format='%(asctime)s: %(message)s', level=logging.INFO, filemode='a')
    app.run(debug=False, host="192.168.100.99", port=54230)

