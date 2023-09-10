#!/usr/bin/env bash

set -eu

pb_prefix='Touchpad:'
pb_suffix='%'

pb_low_power_fg="red"
pb_middle_power_fg="orange"
pb_high_power_fg="green"
pb_no_show=100

function show_touchpad_battery() {
  # 获取蓝牙触摸板的电量百分比
  local battery_percent=$(ioreg -l 2>/dev/null | grep BatteryPercent | awk '{print $NF}')

  if [ -z "${battery_percent}" ]; then
    return 1
  fi

  if ((battery_percent < ${pb_no_show})); then
    if ((battery_percent < 50)); then
      if ((battery_percent <= 30)); then
        # 在status-right配置的末尾添加电量信息
        echo "#[fg=${pb_low_power_fg}]${pb_prefix}${battery_percent}${pb_suffix}"
      else
        echo "#[fg=${pb_middle_power_fg}]${pb_prefix}${battery_percent}${pb_suffix}"
      fi
    else
      echo "#[fg=${pb_high_power_fg}]${pb_prefix}${battery_percent}${pb_suffix}"
    fi
  fi
}

show_touchpad_battery
