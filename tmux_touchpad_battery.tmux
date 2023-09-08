#!/usr/bin/env bash

set -eu

readonly CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SCRITPS_DIR="${CURRENT_DIR}/scripts"

source "${SCRITPS_DIR}/helpers.sh"

placeholders=(
	"\#{touchpad_battery}"
)

commands=(
	"#($SCRITPS_DIR/battery.sh)"
)

do_interpolation() {
	local all_interpolated="$1"
	for ((i = 0; i < ${#commands[@]}; i++)); do
		all_interpolated=${all_interpolated//${placeholders[$i]}/${commands[$i]}}
	done
	echo "$all_interpolated"
}

update_tmux_option() {
	local option="$1"
	local option_value="$(get_tmux_option "$option")"
	local new_option_value="$(do_interpolation "$option_value")"
	set_tmux_option "$option" "$new_option_value"
}

main() {
	update_tmux_option "status-right"
	update_tmux_option "status-left"
}

main
