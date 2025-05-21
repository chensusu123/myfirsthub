#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import json
import os
import subprocess

import git
from base import utils

TAG = "base.git_util.py"


def git_update(path, branch):
    utils.log(TAG, "git_update", f"path:{path} branch:{branch}")
    repo = git.Repo(path)
    repo.git.reset(".")
    repo.git.checkout(".")
    repo.git.clean('-df')
    try:
        remote = repo.remotes.origin
        remote.fetch()
    except Exception as e:
        utils.log(TAG, "git_update", f" remote.fetch e:{e}")
    cur_branch = repo.active_branch
    if cur_branch.name != branch:
        repo.git.checkout(f"{branch}")
        repo.git.branch("-D", cur_branch)
    repo.remotes.origin.pull()


def git_push(path, branch):
    utils.log(TAG, "git_push", f"path:{path} branch:{branch}")
    repo = git.Repo.init(path)
    # 如果冲突 目前需要手动给解决 本次打包失败
    repo.remotes.origin.pull()
    repo.remotes.origin.push(branch)
    return None


def git_tag_max_version(path):
    repo = git.Repo.init(path)
    max_version = 0
    max_version_desc = ""
    for item in repo.tags:
        if '0.0.0' not in item.name:
            continue
        if "-" in item.name:
            version_str = item.name.split("-")[0]
            version_desc = item.name.split("-")[1]
        else:
            version_str = item.name
            version_desc = ""
        version = str.split(version_str, '.').pop()
        version = int(version)
        if max_version <= version:
            max_version = version
            max_version_desc = version_desc
    if len(max_version_desc) > 0:
        return f"0.0.0.{max_version}-{max_version_desc}"
    else:
        return f"0.0.0.{max_version}"


def get_current_branch_name(repo_dir):
    try:
        repo = git.Repo(repo_dir)
        return repo.active_branch.name
    except Exception as e:
        commands = "cd " + repo_dir + " && git name-rev HEAD"
        print("get_git_branch_name cmd-->" + commands)
        return str(subprocess.check_output(commands, shell=True)).replace("HEAD", "").replace("remotes/", "") \
            .replace("origin/", "").replace("b'", "").replace("'", "").replace("\n", "").replace("\\n", "").replace(
            " ", "")


def git_clone(repo_url, clone_dir):
    git.Repo.clone_from(repo_url, clone_dir)


def git_status(repo_dir):
    repo = git.Repo(repo_dir)
    repo.git.status()


def git_pull(repo_dir):
    repo = git.Repo(repo_dir)
    origin = repo.remotes.origin
    origin.pull()


def get_last_commit(repo_dir):
    repo = git.Repo(repo_dir)
    last_commit = repo.head.commit
    return f"Author:{last_commit.author}, Date:{last_commit.committed_datetime}, Message:{last_commit.message}".replace("\n", "").replace("\t", "")


def get_last_commit_id(repo_dir):
    repo = git.Repo(repo_dir)
    latest_commit = repo.head.commit
    return latest_commit.hexsha


def get_diff_list(repo_dir, commit_id_1, commit_id_2):
    try:
        repo = git.Repo(repo_dir)
        commit1 = repo.commit(commit_id_1)
        commit2 = repo.commit(commit_id_2)
        diff = commit1.diff(commit2)
        diff_array = []
        for item in diff:
            item_json = {
                "change_type": item.change_type,
                "a_path": item.a_path,
                "b_path": item.b_path
            }
            diff_array.append(item_json)  # 将字典转为 JSON 字符串
        return diff_array
    except Exception as e:
        utils.log(TAG, "get_diff_list", f"Error: {e}")
        return None


def get_commit_history(repo_dir, commit_id_1, commit_id_2, max_commit_count, max_file_count):
    try:
        repo = git.Repo(repo_dir)
        commit1 = repo.commit(commit_id_1)
        commit2 = repo.commit(commit_id_2)
        commits = list(repo.iter_commits(f'{commit1.hexsha}..{commit2.hexsha}', max_count=max_commit_count))
        commit_history = []
        for commit in commits:
            file_array = get_diff_file_list(commit, commit.parents, max_file_count)
            commit_info = {
                "commit_id": commit.hexsha,
                "author": commit.author.name,
                "date": commit.committed_datetime.strftime("%Y-%m-%d %H:%M:%S"),
                "message": commit.message.strip(),
                "file_array": file_array
            }
            commit_history.append(commit_info)
        return commit_history
    except Exception as e:
        print(f"An error occurred: {e}")
        return None


def get_last_commit_history(repo_dir, branch_name, max_commit_count, max_file_count):
    try:
        repo = git.Repo(repo_dir)
        commits = list(repo.iter_commits(f'{branch_name}', max_count=max_commit_count))
        commit_history = []
        for commit in commits:
            file_array = get_diff_file_list(commit, commit.parents, max_file_count)
            commit_info = {
                "commit_id": commit.hexsha,
                "author": commit.author.name,
                "date": commit.committed_datetime.strftime("%Y-%m-%d %H:%M:%S"),
                "message": commit.message.strip(),
                "file_array": file_array
            }
            commit_history.append(commit_info)
        return commit_history
    except Exception as e:
        print(f"An error occurred: {e}")
        return None


def get_diff_file_list(commit1, commit2, max_count):
    if commit1 is None or commit2 is None:
        return []
    try:
        diff_list = commit1.diff(commit2)
        if diff_list is None:
            return []
        file_array = []
        for diff_item in diff_list:
            if len(file_array) >= max_count:
                break
            file_array.append(diff_item.b_path)
        return file_array
    except Exception as e:
        print(f"An error occurred: {e}")
        return []


def config_pull_rebase(path):
    if not os.path.exists(path):
        raise Exception(f"config_pull_rebase {path} is not exist")
    cur_path = os.getcwd()
    os.chdir(path)
    cmd = "git config pull.rebase true"
    os.system(cmd)
    os.chdir(cur_path)
