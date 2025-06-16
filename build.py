#!/usr/bin/env python3
"""
build.py —— 在 Windows 环境下运行，同步生成：
  1. Windows/amd64 可执行 (myapp-windows.exe)
  2. Linux/amd64 可执行   (myapp-linux)

使用前提：
  - Windows 下已安装 Go，并且 go.exe 在 PATH 中可直接调用。
  - 项目根目录下没有任何 `import "C"`（即禁用 cgo）。
  - 将此脚本与 main.go 放在同一目录，并在 PowerShell/命令行中执行即可。
"""

import os
import sys
import subprocess
import argparse

def run_build(goos: str, goarch: str, cgo_enabled: str, output_name: str) -> bool:
    """
    设置环境变量并调用 go build。返回 True 表示成功，False 表示失败。
    """
    # 复制当前环境变量
    env = os.environ.copy()
    env["GOOS"] = goos
    env["GOARCH"] = goarch
    env["CGO_ENABLED"] = cgo_enabled

    cmd = ["go", "build", "-o", output_name]
    print(f"\n⏳ 正在执行: GOOS={goos} GOARCH={goarch} CGO_ENABLED={cgo_enabled} go build -o {output_name}")
    result = subprocess.run(cmd, env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)

    if result.returncode != 0:
        print(f"❌ 构建失败: {output_name}")
        print("------ stderr ------")
        print(result.stderr.strip())
        print("--------------------")
        return False

    print(f"✅ 已成功生成: {output_name}")
    return True

def main():
    parser = argparse.ArgumentParser(description="同时构建 Windows 和 Linux 可执行（禁用 cgo）")
    parser.add_argument(
        "-n", "--name", default="gin-go", 
        help="输出文件的基础名称（默认: gin-go）"
    )
    args = parser.parse_args()
    base_name = args.name

    print(f"🔨 开始构建，项目目录：{os.getcwd()}")

    # Windows/amd64
    win_output = f"{base_name}-windows.exe"
    success_win = run_build(goos="windows", goarch="amd64", cgo_enabled="0", output_name=win_output)
    if not success_win:
        sys.exit(1)

    # Linux/amd64
    linux_output = f"{base_name}-linux"
    success_linux = run_build(goos="linux", goarch="amd64", cgo_enabled="0", output_name=linux_output)
    if not success_linux:
        sys.exit(1)

    print("\n🎉 全部构建完成！")
    print(f"• Windows 可执行路径: {os.path.join(os.getcwd(), win_output)}")
    print(f"• Linux   可执行路径: {os.path.join(os.getcwd(), linux_output)}")

if __name__ == "__main__":
    main()
