package display

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/akayj/tmux-touchpad-battery/internal/battery"
	"github.com/akayj/tmux-touchpad-battery/internal/tmux"
)

// BatteryFormatter 负责格式化电池显示
type BatteryFormatter struct {
	config      *tmux.Config
	batteryInfo *battery.BatteryInfo
}

// NewBatteryFormatter 创建新的电池格式化器
func NewBatteryFormatter(config *tmux.Config) *BatteryFormatter {
	return &BatteryFormatter{
		config: config,
	}
}

// SetBatteryInfo 设置电池信息
func (f *BatteryFormatter) SetBatteryInfo(info *battery.BatteryInfo) {
	f.batteryInfo = info
}

// Format 格式化电池信息为 tmux 状态栏显示
func (f *BatteryFormatter) Format() string {
	if f.batteryInfo == nil || !f.batteryInfo.Available {
		return ""
	}

	// 如果电量达到不显示阈值，则不显示
	if f.batteryInfo.Percentage >= f.config.NotShowThreshold {
		return ""
	}

	color := f.getBatteryColor(f.batteryInfo)
	if color == "" {
		return ""
	}

	// 检查是否需要闪烁（低电量且未充电时）
	blinkAttr := ""
	if f.shouldBlink(f.batteryInfo) {
		blinkAttr = ",blink"
	}

	// 添加充电图标
	chargingIcon := ""
	if f.batteryInfo.IsCharging && f.config.ShowChargingIcon {
		chargingIcon = f.config.ChargingIcon
	}

	// 格式化为 tmux 颜色格式
	return fmt.Sprintf("#[fg=%s%s]%s%d%s%s",
		color,
		blinkAttr,
		f.config.PercentPrefix,
		f.batteryInfo.Percentage,
		f.config.PercentSuffix,
		chargingIcon,
	)
}

// FormatWithStyle 使用 lipgloss 格式化电池信息（用于终端显示）
func (f *BatteryFormatter) FormatWithStyle() string {
	if f.batteryInfo == nil || !f.batteryInfo.Available {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("Touchpad not connected")
	}

	color := f.getBatteryLipglossColor(f.batteryInfo)
	style := lipgloss.NewStyle().Foreground(color)

	text := fmt.Sprintf("%s%d%s",
		f.config.PercentPrefix,
		f.batteryInfo.Percentage,
		f.config.PercentSuffix,
	)

	if f.batteryInfo.IsCharging && f.config.ShowChargingIcon {
		text += " " + f.config.ChargingIcon
	}

	return style.Render(text)
}

// FormatBattery 格式化指定的电池信息为 tmux 状态栏显示（向后兼容）
func (f *BatteryFormatter) FormatBattery(info *battery.BatteryInfo) string {
	f.SetBatteryInfo(info)
	return f.Format()
}

// FormatBatteryWithStyle 使用 lipgloss 格式化指定的电池信息（向后兼容）
func (f *BatteryFormatter) FormatBatteryWithStyle(info *battery.BatteryInfo) string {
	f.SetBatteryInfo(info)
	return f.FormatWithStyle()
}

// getBatteryColor 获取电池颜色（tmux 格式）
func (f *BatteryFormatter) getBatteryColor(info *battery.BatteryInfo) string {
	if info.IsCharging {
		return f.config.ColorCharging
	}

	if info.Percentage < f.config.StressThreshold {
		return f.config.ColorStress
	} else if info.Percentage < f.config.MediumThreshold {
		return f.config.ColorMedium
	} else if info.Percentage < f.config.NotShowThreshold {
		return f.config.ColorHigh
	}

	return ""
}

// getBatteryLipglossColor 获取电池颜色（lipgloss 格式）
func (f *BatteryFormatter) getBatteryLipglossColor(info *battery.BatteryInfo) lipgloss.Color {
	if info.IsCharging {
		return f.tmuxColorToLipgloss(f.config.ColorCharging)
	}

	if info.Percentage < f.config.StressThreshold {
		return f.tmuxColorToLipgloss(f.config.ColorStress)
	} else if info.Percentage < f.config.MediumThreshold {
		return f.tmuxColorToLipgloss(f.config.ColorMedium)
	} else {
		return f.tmuxColorToLipgloss(f.config.ColorHigh)
	}
}

// tmuxColorToLipgloss 将 tmux 颜色转换为 lipgloss 颜色
func (f *BatteryFormatter) tmuxColorToLipgloss(tmuxColor string) lipgloss.Color {
	colorMap := map[string]string{
		"red":     "#FF0000",
		"green":   "#00FF00",
		"yellow":  "#FFFF00",
		"blue":    "#0000FF",
		"magenta": "#FF00FF",
		"cyan":    "#00FFFF",
		"white":   "#FFFFFF",
		"black":   "#000000",
	}

	if color, exists := colorMap[tmuxColor]; exists {
		return lipgloss.Color(color)
	}

	// 如果是十六进制颜色或数字，直接返回
	return lipgloss.Color(tmuxColor)
}

// shouldBlink 判断是否应该闪烁
func (f *BatteryFormatter) shouldBlink(info *battery.BatteryInfo) bool {
	// 只有在以下条件都满足时才闪烁：
	// 1. 启用了闪烁功能
	// 2. 电量低于压力阈值
	// 3. 没有在充电
	return f.config.BlinkOnLowBattery &&
		info.Percentage < f.config.StressThreshold &&
		!info.IsCharging
}
