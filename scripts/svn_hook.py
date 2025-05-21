import json
import os
import logging
import jenkins_api
from flask import Flask, request, jsonify, send_file

app = Flask(__name__)


def on_changed(svn_address):
    data = request.get_json()
    check_execute_task(svn_address, data)
    return f"success", 200


def check_execute_task(svn_address, request_data):
    user = request_data["user"]
    commit = request_data["commit"]
    file_array = request_data["files"]
    logging.info(f"check_execute_task >>> svn_address:{svn_address} user:{user} commit:{commit} file_array:{file_array}")
    if "192.168.8.69:3692" not in svn_address:
        logging.info("check_execute_task >>> svn_address is error")
    elif not check_gpu_particle_bake(file_array):
        logging.info("check_execute_task >>> file is not changed")
    else:
        execute_gpu_particle_bake()


def check_gpu_particle_bake(file_array):
    if file_array is None or len(file_array) == 0:
        logging.info(f"check_gpu_particle_bake >>> file array is empty")
        return False
    for item in file_array:
        if "Assets/Res/RemoteRes/Prefab/Effects/EscortSkill/Prefab" in item:
            logging.info(f"check_gpu_particle_bake >>> already found need build ")
            return True
    logging.info(f"check_gpu_particle_bake >>> not found file to build")
    return False


def execute_gpu_particle_bake():
    logging.info("execute_gpu_particle_bake")
    jenkins_api.execute_command( 'build resource/develop/gpu_particle_bake')


def all_url_rule():
    app.add_url_rule("/svn/<string:svn_address>", "on_changed", on_changed, methods=["GET", "POST"])


if __name__ == '__main__':
    all_url_rule()
    logging.basicConfig(filename='run.log', format='%(asctime)s: %(message)s', level=logging.INFO, filemode='a')
    app.run(debug=False, host="192.168.100.99", port=54226)

