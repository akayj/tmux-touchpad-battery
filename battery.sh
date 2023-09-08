#!/usr/bin/env bash

set -eu

function show_touchpad_battery() {
	# 获取蓝牙触摸板的电量百分比
	local battery_percent=$(ioreg -l | grep BatteryPercent | awk '{print $NF}')

	if ((battery_percent < 50)); then
		if ((battery_percent <= 30)); then
			# 在status-right配置的末尾添加电量信息
			echo "#[fg=red]Touchpad:${battery_percent}%"
		else
			echo "#[fg=orange]Touchpad:${battery_percent}%"
		fi
	else
		echo "#[fg=blue]Touchpad:${battery_percent}%"
	fi
}

show_touchpad_battery
