#!/usr/bin/env bash

set -eu

LC_NUMERIC=C

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$CURRENT_DIR/helpers.sh"

tpb_percent_prefix=$(get_tmux_option "@tpb_percent_prefix" "Touchpad:")
tpb_percent_suffix=$(get_tmux_option "@tpb_percent_suffix" "%")

tpb_color_charging=$(get_tmux_option "@tpb_color_charging" "green")

tpb_color_high=$(get_tmux_option "@tpb_color_high" "white")
tpb_color_medium=$(get_tmux_option "@tpb_color_medium" "yellow")
tpb_color_stress=$(get_tmux_option "@tpb_color_stress" "red")

tpb_stress_threshold=$(get_tmux_option "@tpb_stress_threshold" "30")
tpb_medium_threshold=$(get_tmux_option "@tpb_medium_threshold" "80")
tpb_not_show_threshold=$(get_tmux_option "@tpb_not_show_threshold" "100")

get_trackpad_battery_status () {
  local status=$(ioreg -l 2>/dev/null | grep BatteryStatusFlags | awk '{print $NF}')

  if [ "$status" == "3" ]; then
    echo 1 # charging
  else
    echo 0 # not charging
  fi
}

function get_battery_color(){
  local battery_percent=$1

  if [ $(get_trackpad_battery_status) == "1" ]; then
    echo "$tpb_color_charging";
    return
  fi

  if fcomp "$battery_percent" "$tpb_stress_threshold"; then
    echo "$tpb_color_stress";
  elif fcomp "$battery_percent" "$tpb_medium_threshold"; then
    echo "$tpb_color_medium";
  elif fcomp "$battery_percent" "$tpb_not_show_threshold"; then
    echo "$tpb_color_high";
  fi
}

function show_touchpad_battery() {
  # 获取蓝牙触摸板的电量百分比
  local battery_percent=$(ioreg -l 2>/dev/null | grep BatteryPercent | awk '{print $NF}')

  if [ -z "${battery_percent}" ]; then
    return 1
  fi

  local battery_print_color=$(get_battery_color "$battery_percent")
  if [ -n "$battery_print_color" ]; then
    echo "#[fg=${battery_print_color}]${tpb_percent_prefix}${battery_percent}${tpb_percent_suffix}"
  fi
}

show_touchpad_battery
