#!/usr/bin/env python3
"""
Git 自动提交并递增版本号打 Tag 推送到远程。

功能：
1. 提交所有本地更改到 Git。
2. 获取最新的版本 Tag（格式 va.b.c 或 path/va.b.c）。
3. 根据参数递增版本号（1=c+1, 2=b+1, 3=a+1）。
4. 打上新的 Tag 并推送到远程仓库。
"""

import argparse
import subprocess
import sys
from datetime import datetime


def run_command(command):
    """执行 Git 命令并返回结果"""
    try:
        result = subprocess.run(command, shell=True, capture_output=True, text=True, check=True)
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        return None


def get_latest_tag(tag_pattern):
    """获取匹配模式的最新 tag"""
    command = f'git describe --tags --abbrev=0 --match "{tag_pattern}"'
    return run_command(command)


def commit_changes():
    """提交所有本地更改"""
    print("提交所有本地更改...")
    run_command("git add .")
    
    commit_message = f"Auto commit at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}"
    result = run_command(f'git commit -m "{commit_message}"')
    
    if result is None:
        print("无更改需要提交")
    else:
        print(f"提交信息: {commit_message}")


def increment_version(version_parts, level):
    """根据递增级别递增版本号"""
    a, b, c = version_parts
    
    if level == 1:
        c += 1
        print("递增修订号 (c)")
    elif level == 2:
        b += 1
        c = 0
        print("递增次版本号 (b)")
    elif level == 3:
        a += 1
        b = 0
        c = 0
        print("递增主版本号 (a)")
    
    return a, b, c


def main():
    parser = argparse.ArgumentParser(description="Git 自动提交并递增版本号打 Tag 推送到远程")
    parser.add_argument("-lv", "--level", type=int, choices=[1, 2, 3], default=1,
                        help="版本号递增级别: 1=修订号(c)+1(默认), 2=次版本号(b)+1, 3=主版本号(a)+1")
    parser.add_argument("-sp", "--subPath", type=str, default="",
                        help='子包路径，如 "cmd/xf"，如果不提供则为主包')
    
    args = parser.parse_args()
    
    # 1. 提交所有本地更改
    commit_changes()
    
    # 2. 构造 Tag 前缀
    tag_prefix = f"{args.subPath}/" if args.subPath else ""
    
    # 3. 获取最新的版本 Tag
    tag_pattern = f"{tag_prefix}v*.*.*" if tag_prefix else "v*.*.*"
    latest_tag = get_latest_tag(tag_pattern)
    
    print(f"最新 {tag_pattern} 版本 Tag: {latest_tag}")
    
    if not latest_tag:
        # 如果没有找到匹配的 Tag，尝试从主版本继承或使用默认值
        if tag_prefix:
            # 对于子包，尝试获取主包的最新版本
            main_latest_tag = get_latest_tag("v*.*.*")
            if main_latest_tag:
                # 使用主包版本作为起点
                version_part = main_latest_tag.replace('v', '')
                latest_tag = f"{tag_prefix}{version_part}"
                print(f"未找到子包 {args.subPath} 的版本 Tag，将从主包版本 {main_latest_tag} 开始")
            else:
                latest_tag = f"{tag_prefix}0.0.0"
                print(f"未找到任何版本 Tag，将从 {latest_tag} 开始")
        else:
            latest_tag = "v0.0.0"  # 如果没有 Tag，默认从 v0.0.0 开始
            print(f"未找到版本 Tag，将从 {latest_tag} 开始")
    
    # 4. 解析版本号并递增
    if tag_prefix:
        version_string = latest_tag[len(tag_prefix):].lstrip('v')
    else:
        version_string = latest_tag.lstrip('v')
    
    version_parts = list(map(int, version_string.split('.')))
    a, b, c = version_parts
    
    # 递增版本号
    a, b, c = increment_version([a, b, c], args.level)
    
    # 构造新 tag
    new_tag = f"{tag_prefix}v{a}.{b}.{c}" if tag_prefix else f"v{a}.{b}.{c}"
    
    print(f"当前版本: {latest_tag} → 新版本: {new_tag}")
    
    # 5. 打 Tag 并推送到远程
    run_command(f"git tag {new_tag}")
    run_command("git push origin --tags")
    run_command("git push origin")  # 推送提交
    
    print(f"已提交并推送 Tag: {new_tag}")


if __name__ == "__main__":
    main()
