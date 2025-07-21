package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/akayj/tmux-touchpad-battery/internal/battery"
	"github.com/akayj/tmux-touchpad-battery/internal/display"
	"github.com/akayj/tmux-touchpad-battery/internal/tmux"
	"github.com/akayj/tmux-touchpad-battery/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var (
		showStatus = flag.Bool("status", false, "显示电池状态")
		showUI     = flag.Bool("ui", false, "启动交互式 UI")
		showHelp   = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	if *showUI {
		runUI()
		return
	}

	if *showStatus {
		showBatteryStatus()
		return
	}

	// 默认行为：输出 tmux 格式
	outputTmuxFormat()
}

func printHelp() {
	fmt.Println("Tmux Touchpad Battery - Golang 版本")
	fmt.Println("用法:")
	fmt.Println("  tmux-touchpad-battery           输出 tmux 格式（默认）")
	fmt.Println("  tmux-touchpad-battery -status   显示电池状态")
	fmt.Println("  tmux-touchpad-battery -ui       启动交互式 UI")
	fmt.Println("  tmux-touchpad-battery -help     显示此帮助信息")
	fmt.Println()
	fmt.Println("配置选项:")
	fmt.Println("  @tpb_percent_prefix      显示前缀 (默认: 'Touchpad:')")
	fmt.Println("  @tpb_percent_suffix      显示后缀 (默认: '%')")
	fmt.Println("  @tpb_color_charging      充电时颜色 (默认: 'green')")
	fmt.Println("  @tpb_color_high          高电量颜色 (默认: 'white')")
	fmt.Println("  @tpb_color_medium        中等电量颜色 (默认: 'yellow')")
	fmt.Println("  @tpb_color_stress        低电量颜色 (默认: 'red')")
	fmt.Println("  @tpb_stress_threshold    低电量阈值 (默认: 30)")
	fmt.Println("  @tpb_medium_threshold    中等电量阈值 (默认: 80)")
	fmt.Println("  @tpb_not_show_threshold  不显示阈值 (默认: 100)")
	fmt.Println("  @tpb_blink_on_low_battery 低电量时闪烁 (默认: 'off')")
	fmt.Println("  @tpb_charging_icon       充电图标 (默认: '⚡')")
	fmt.Println("  @tpb_show_charging_icon  显示充电图标 (默认: 'on')")
}

func runUI() {
	model := ui.NewModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("启动 UI 失败: %v\n", err)
		os.Exit(1)
	}
}

func showBatteryStatus() {
	config := tmux.GetConfig()
	formatter := display.NewFormatter(config)

	info, err := battery.GetTouchpadBatteryInfo()
	if err != nil {
		fmt.Printf("获取电池信息失败: %v\n", err)
		os.Exit(1)
	}

	if !info.Available {
		fmt.Println("触摸板未连接或无电池信息")
		return
	}

	fmt.Printf("电池电量: %d%%\n", info.Percentage)
	fmt.Printf("充电状态: %v\n", info.IsCharging)
	fmt.Printf("格式化输出: %s\n", formatter.FormatBatteryWithStyle(info))
}

func outputTmuxFormat() {
	config := tmux.GetConfig()
	formatter := display.NewFormatter(config)

	info, err := battery.GetTouchpadBatteryInfo()
	if err != nil {
		// 静默失败，不输出任何内容
		return
	}

	output := formatter.FormatBattery(info)
	fmt.Print(output)
}
