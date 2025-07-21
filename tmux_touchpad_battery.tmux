#!/usr/bin/env bash

set -eu

readonly CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 检查 Go 二进制文件是否存在
BINARY_PATH="${CURRENT_DIR}/bin/tmux-touchpad-battery"

# 检查是否需要构建或重新构建
need_build=false

if [ ! -f "$BINARY_PATH" ]; then
    need_build=true
else
    # 检查源代码是否比二进制文件新
    if [ -d "${CURRENT_DIR}/cmd" ] || [ -d "${CURRENT_DIR}/internal" ]; then
        # 查找最新的 .go 文件
        newest_go_file=$(find "${CURRENT_DIR}" -name "*.go" -type f -exec stat -f "%m %N" {} \; 2>/dev/null | sort -nr | head -1 | cut -d' ' -f2-)
        if [ -n "$newest_go_file" ] && [ "$newest_go_file" -nt "$BINARY_PATH" ]; then
            need_build=true
        fi
    fi
fi

# 如果需要构建
if [ "$need_build" = true ]; then
    echo "Building tmux-touchpad-battery..." >&2
    mkdir -p "${CURRENT_DIR}/bin"
    
    # 检查是否有 go 命令
    if ! command -v go >/dev/null 2>&1; then
        echo "Error: Go is not installed or not in PATH" >&2
        exit 1
    fi
    
    # 切换到项目目录并构建
    cd "$CURRENT_DIR"
    if ! go build -o "$BINARY_PATH" ./cmd/tmux-touchpad-battery 2>/dev/null; then
        echo "Error: Failed to build tmux-touchpad-battery" >&2
        exit 1
    fi
    echo "Build completed successfully" >&2
fi

# 获取 tmux 选项的辅助函数
get_tmux_option() {
  local option="${1?option is required}"
  local default_value="${2:-}"

  local option_value="$(tmux show-option -gqv "$option")"
  if [ -z "$option_value" ]; then
    echo "$default_value"
  else
    echo "$option_value"
  fi
}

# 设置 tmux 选项的辅助函数
set_tmux_option() {
  local option="$1"
  local value="$2"
  tmux set-option -gq "$option" "$value"
}

# 占位符和命令
placeholders=(
  "\#{touchpad_battery}"
)

commands=(
  "#($BINARY_PATH)"
)

# 执行插值替换
do_interpolation() {
  local all_interpolated="$1"
  for ((i = 0; i < ${#commands[@]}; i++)); do
    all_interpolated=${all_interpolated//${placeholders[$i]}/${commands[$i]}}
  done
  echo "$all_interpolated"
}

# 更新 tmux 选项
update_tmux_option() {
  local option="$1"
  local option_value="$(get_tmux_option "$option")"
  local new_option_value="$(do_interpolation "$option_value")"
  set_tmux_option "$option" "$new_option_value"
}

# 主函数
main() {
  update_tmux_option "status-right"
  update_tmux_option "status-left"
}

main