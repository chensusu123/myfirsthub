import subprocess
import os
from base import utils


def checkout(username, password, url, work_path):
    print(f"BEGIN (checkout): {url} {work_path}")
    command = ["svn", "checkout", url, work_path, "--username", username, "--password", password, "--non-interactive", "--trust-server-cert"]
    result = subprocess.run(command, capture_output=True, text=True)
    print("STDOUT (checkout):\n", result.stdout)
    print("STDERR (checkout):\n", result.stderr)
    print("ERROR_CODE (checkout):\n", result.returncode)


def revert(work_path):
    print(f"BEGIN (revert): {work_path}")
    command = ["svn", "revert", "-R", work_path]
    result = subprocess.run(command, capture_output=True, text=True)
    print("STDOUT (revert):\n", result.stdout)
    print("STDERR (revert):\n", result.stderr)
    print("ERROR_CODE (checkout):\n", result.returncode)


def remove_add(work_path):
    print(f"BEGIN (remove_add): {work_path}")
    command_status = ["svn", "status", work_path]
    result = subprocess.run(command_status, capture_output=True, text=True)
    print("STDOUT (remove_add):\n", result.stdout)
    print("STDERR (remove_add):\n", result.stderr)
    print("ERROR_CODE (checkout):\n", result.returncode)
    # 查找新增的文件（状态为 "A" 或 "?"）
    for line in result.stdout.splitlines():
        if line.startswith("A") or line.startswith("?"):
            # 提取文件路径
            file_path = os.path.join(work_path, line[8:].strip())
            print(f"Deleted : {file_path}")
            utils.delete_folder_or_file(file_path)


def add(work_path):
    print(f"BEGIN (add): {work_path}")
    command_add = ["svn", "add", "--force", work_path]
    result_add = subprocess.run(command_add, capture_output=True, text=True)
    print("STDOUT (add):\n", result_add.stdout)
    print("STDERR (add):\n", result_add.stderr)
    print("ERROR_CODE (checkout):\n", result_add.returncode)


def commit(username, password, work_path, commit_message):
    print(f"BEGIN (commit): {work_path}")
    command_commit = ["svn", "commit", work_path, "-m", commit_message, "--username", username, "--password", password]
    result_commit = subprocess.run(command_commit, capture_output=True, text=True)
    print("STDOUT (commit):\n", result_commit.stdout)
    print("STDERR (commit):\n", result_commit.stderr)
    print("ERROR_CODE (checkout):\n", result_commit.returncode)


def update(username, password, work_path):
    command_update = ["svn", "update", work_path, "--username", username, "--password", password]
    result_update = subprocess.run(command_update, capture_output=True, text=True)
    print("STDOUT (update):\n", result_update.stdout)
    print("STDERR (update):\n", result_update.stderr)
    print("ERROR_CODE (checkout):\n", result_update.returncode)


def clean_up(work_path):
    command_update = ["svn", "cleanup", work_path]
    result_update = subprocess.run(command_update, capture_output=True, text=True)
    print("STDOUT (update):\n", result_update.stdout)
    print("STDERR (update):\n", result_update.stderr)
    print("ERROR_CODE (checkout):\n", result_update.returncode)



