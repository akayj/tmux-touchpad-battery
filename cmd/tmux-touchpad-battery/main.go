package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/akayj/tmux-touchpad-battery/internal/battery"
	"github.com/akayj/tmux-touchpad-battery/internal/display"
	"github.com/akayj/tmux-touchpad-battery/internal/system"
	"github.com/akayj/tmux-touchpad-battery/internal/tmux"
	"github.com/akayj/tmux-touchpad-battery/internal/ui"
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
	fmt.Println("  @tpb_show_cpu_info       显示 CPU 信息 (默认: 'on')")
	fmt.Println("  @tpb_show_gpu_info       显示 GPU 信息 (默认: 'on')")
	fmt.Println("  @tpb_system_info_prefix  系统信息前缀 (默认: '')")
	fmt.Println("  @tpb_system_info_suffix  系统信息后缀 (默认: '')")
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
	batteryFormatter := display.NewFormatter(config)
	systemFormatter := display.NewSystemFormatter(config)

	// 获取电池信息
	batteryInfo, err := battery.GetTouchpadBatteryInfo()
	if err != nil {
		fmt.Printf("获取电池信息失败: %v\n", err)
		os.Exit(1)
	}

	// 获取系统信息
	systemInfo, err := system.GetSystemInfo()
	if err != nil {
		fmt.Printf("获取系统信息失败: %v\n", err)
		os.Exit(1)
	}

	if batteryInfo.Available {
		fmt.Printf("电池电量: %d%%\n", batteryInfo.Percentage)
		fmt.Printf("充电状态: %v\n", batteryInfo.IsCharging)
		fmt.Printf("电池格式化输出: %s\n", batteryFormatter.FormatBatteryWithStyle(batteryInfo))
	}

	if systemInfo.Available {
		fmt.Printf("CPU 使用率: %.1f%%\n", systemInfo.CPUUsage)
		fmt.Printf("GPU 使用率: %.1f%%\n", systemInfo.GPUUsage)
		fmt.Printf("系统格式化输出: %s\n", systemFormatter.FormatSystemInfoWithStyle(systemInfo))
	}
}

func outputTmuxFormat() {
	config := tmux.GetConfig()
	batteryFormatter := display.NewFormatter(config)
	systemFormatter := display.NewSystemFormatter(config)

	// 获取电池信息
	batteryInfo, err := battery.GetTouchpadBatteryInfo()
	if err != nil {
		// 静默失败，不输出任何内容
		return
	}

	// 获取系统信息
	systemInfo, err := system.GetSystemInfo()
	if err != nil {
		// 静默失败，不输出任何内容
		return
	}

	// 格式化输出
	batteryOutput := batteryFormatter.FormatBattery(batteryInfo)
	systemOutput := systemFormatter.FormatSystemInfo(systemInfo)

	// 如果两个输出都为空，则不输出任何内容
	if batteryOutput == "" && systemOutput == "" {
		return
	}

	// 输出结果，用空格分隔
	var outputs []string
	if batteryOutput != "" {
		outputs = append(outputs, batteryOutput)
	}
	if systemOutput != "" {
		outputs = append(outputs, systemOutput)
	}

	if len(outputs) > 0 {
		fmt.Print(outputs[0])
		for _, output := range outputs[1:] {
			fmt.Print(" ", output)
		}
	}
}
