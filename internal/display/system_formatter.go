package display

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/akayj/tmux-touchpad-battery/internal/system"
	"github.com/akayj/tmux-touchpad-battery/internal/tmux"
)

// SystemFormatter 负责格式化系统信息（CPU/GPU 使用率）
type SystemFormatter struct {
	config *tmux.Config
}

// NewSystemFormatter 创建新的系统信息格式化器
func NewSystemFormatter(config *tmux.Config) *SystemFormatter {
	return &SystemFormatter{
		config: config,
	}
}

// FormatSystemInfo 格式化系统信息为 tmux 状态栏显示
func (f *SystemFormatter) FormatSystemInfo(info *system.SystemInfo) string {
	if !info.Available {
		return ""
	}

	// 检查是否启用了 CPU 和 GPU 信息显示
	if !f.config.ShowCPUInfo && !f.config.ShowGPUInfo {
		return ""
	}

	var parts []string

	// 添加前缀
	if f.config.SystemInfoPrefix != "" {
		parts = append(parts, f.config.SystemInfoPrefix)
	}

	// 添加 CPU 信息
	if f.config.ShowCPUInfo {
		cpuText := fmt.Sprintf("CPU:%.1f%%", info.CPUUsage)
		parts = append(parts, cpuText)
	}

	// 添加 GPU 信息
	if f.config.ShowGPUInfo {
		if info.GPUUsage == 0 {
			// GPU 使用率为 0 可能是因为权限问题
			parts = append(parts, "GPU:N/A")
		} else {
			gpuText := fmt.Sprintf("GPU:%.1f%%", info.GPUUsage)
			parts = append(parts, gpuText)
		}
	}

	// 添加后缀
	if f.config.SystemInfoSuffix != "" {
		parts = append(parts, f.config.SystemInfoSuffix)
	}

	if len(parts) == 0 {
		return ""
	}

	// 使用默认颜色
	color := "white"
	return fmt.Sprintf("#[fg=%s]%s", color, strings.Join(parts, " "))
}

// FormatSystemInfoWithStyle 使用 lipgloss 格式化系统信息（用于终端显示）
func (f *SystemFormatter) FormatSystemInfoWithStyle(info *system.SystemInfo) string {
	if !info.Available {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("System info not available")
	}

	// 检查是否启用了 CPU 和 GPU 信息显示
	if !f.config.ShowCPUInfo && !f.config.ShowGPUInfo {
		return ""
	}

	var parts []string

	// 添加前缀
	if f.config.SystemInfoPrefix != "" {
		parts = append(parts, f.config.SystemInfoPrefix)
	}

	// 添加 CPU 信息
	if f.config.ShowCPUInfo {
		cpuText := fmt.Sprintf("CPU: %.1f%%", info.CPUUsage)
		parts = append(parts, cpuText)
	}

	// 添加 GPU 信息
	if f.config.ShowGPUInfo {
		if info.GPUUsage == 0 {
			// GPU 使用率为 0 可能是因为权限问题
			parts = append(parts, "GPU: N/A")
		} else {
			gpuText := fmt.Sprintf("GPU: %.1f%%", info.GPUUsage)
			parts = append(parts, gpuText)
		}
	}

	// 添加后缀
	if f.config.SystemInfoSuffix != "" {
		parts = append(parts, f.config.SystemInfoSuffix)
	}

	if len(parts) == 0 {
		return ""
	}

	// 使用默认颜色
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	return style.Render(strings.Join(parts, " "))
}
