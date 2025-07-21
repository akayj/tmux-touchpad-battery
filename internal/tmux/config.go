package tmux

import (
	"os/exec"
	"strconv"
	"strings"
)

// Config 表示 tmux 配置
type Config struct {
	PercentPrefix     string
	PercentSuffix     string
	ColorCharging     string
	ColorHigh         string
	ColorMedium       string
	ColorStress       string
	StressThreshold   int
	MediumThreshold   int
	NotShowThreshold  int
	BlinkOnLowBattery bool
	ChargingIcon      string
	ShowChargingIcon  bool
}

// GetConfig 获取 tmux 配置
func GetConfig() *Config {
	return &Config{
		PercentPrefix:     getTmuxOption("@tpb_percent_prefix", "Touchpad:"),
		PercentSuffix:     getTmuxOption("@tpb_percent_suffix", "%"),
		ColorCharging:     getTmuxOption("@tpb_color_charging", "green"),
		ColorHigh:         getTmuxOption("@tpb_color_high", "white"),
		ColorMedium:       getTmuxOption("@tpb_color_medium", "yellow"),
		ColorStress:       getTmuxOption("@tpb_color_stress", "red"),
		StressThreshold:   getTmuxOptionInt("@tpb_stress_threshold", 30),
		MediumThreshold:   getTmuxOptionInt("@tpb_medium_threshold", 80),
		NotShowThreshold:  getTmuxOptionInt("@tpb_not_show_threshold", 100),
		BlinkOnLowBattery: getTmuxOptionBool("@tpb_blink_on_low_battery", false),
		ChargingIcon:      getTmuxOption("@tpb_charging_icon", "⚡"),
		ShowChargingIcon:  getTmuxOptionBool("@tpb_show_charging_icon", true),
	}
}

// getTmuxOption 获取 tmux 选项值
func getTmuxOption(option, defaultValue string) string {
	cmd := exec.Command("tmux", "show-option", "-gqv", option)
	output, err := cmd.Output()
	if err != nil {
		return defaultValue
	}

	value := strings.TrimSpace(string(output))
	if value == "" {
		return defaultValue
	}

	return value
}

// getTmuxOptionInt 获取 tmux 选项整数值
func getTmuxOptionInt(option string, defaultValue int) int {
	value := getTmuxOption(option, "")
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// getTmuxOptionBool 获取 tmux 选项布尔值
func getTmuxOptionBool(option string, defaultValue bool) bool {
	value := getTmuxOption(option, "")
	if value == "" {
		return defaultValue
	}

	// 支持多种布尔值表示方式
	switch strings.ToLower(value) {
	case "true", "1", "yes", "on", "enable", "enabled":
		return true
	case "false", "0", "no", "off", "disable", "disabled":
		return false
	default:
		return defaultValue
	}
}

// SetTmuxOption 设置 tmux 选项
func SetTmuxOption(option, value string) error {
	cmd := exec.Command("tmux", "set-option", "-gq", option, value)
	return cmd.Run()
}
