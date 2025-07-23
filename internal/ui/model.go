package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/akayj/tmux-touchpad-battery/internal/battery"
	"github.com/akayj/tmux-touchpad-battery/internal/display"
	"github.com/akayj/tmux-touchpad-battery/internal/system"
	"github.com/akayj/tmux-touchpad-battery/internal/tmux"
)

// Model 表示 TUI 模型
type Model struct {
	batteryInfo  *battery.BatteryInfo
	systemInfo   *system.SystemInfo
	formatter    display.Formatter
	sysFormatter display.Formatter
	config       *tmux.Config
	err          error
	quitting     bool
}

// tickMsg 定时更新消息
type tickMsg time.Time

// NewModel 创建新的 TUI 模型
func NewModel() *Model {
	config := tmux.GetConfig()
	batteryFormatter := display.NewBatteryFormatter(config)
	systemFormatter := display.NewSystemFormatter(config)

	return &Model{
		config:       config,
		formatter:    batteryFormatter,
		sysFormatter: systemFormatter,
	}
}

// Init 初始化模型
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.updateBattery(),
		m.updateSystemInfo(),
		m.tick(),
	)
}

// Update 更新模型
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			// 手动刷新
			return m, tea.Batch(
				m.updateBattery(),
				m.updateSystemInfo(),
			)
		}

	case tickMsg:
		// 定时更新
		return m, tea.Batch(
			m.updateBattery(),
			m.updateSystemInfo(),
			m.tick(),
		)

	case *battery.BatteryInfo:
		m.batteryInfo = msg
		m.err = nil

	case *system.SystemInfo:
		m.systemInfo = msg
		m.err = nil

	case error:
		m.err = msg
	}

	return m, nil
}

// View 渲染视图
func (m *Model) View() string {
	if m.quitting {
		return ""
	}

	var content string

	// 标题
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Tmux Touchpad Battery Monitor")

	content += title + "\n\n"

	// 错误信息
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
		content += errorStyle.Render("Error: "+m.err.Error()) + "\n\n"
	}

	// 电池信息
	if m.batteryInfo != nil {
		// 设置电池信息并格式化
		if bf, ok := m.formatter.(*display.BatteryFormatter); ok {
			bf.SetBatteryInfo(m.batteryInfo)
		}
		batteryDisplay := m.formatter.FormatWithStyle()
		content += "Battery Status: " + batteryDisplay + "\n\n"

		// 详细信息
		if m.batteryInfo.Available {
			detailStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888"))

			details := ""
			details += detailStyle.Render("Percentage: ") +
				lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%d%%", m.batteryInfo.Percentage)) + "\n"

			chargingText := "No"
			if m.batteryInfo.IsCharging {
				chargingText = "Yes ⚡"
			}
			details += detailStyle.Render("Charging: ") +
				lipgloss.NewStyle().Bold(true).Render(chargingText) + "\n"

			content += details + "\n"
		}
	} else {
		content += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("Loading battery information...") + "\n\n"
	}

	// 系统信息
	if m.systemInfo != nil {
		// 设置系统信息并格式化
		if sf, ok := m.sysFormatter.(*display.SystemInfoFormatter); ok {
			sf.SetSystemInfo(m.systemInfo)
		}
		systemDisplay := m.sysFormatter.FormatWithStyle()
		if systemDisplay != "" {
			content += "System Status: " + systemDisplay + "\n\n"

			// 详细信息
			if m.systemInfo.Available {
				detailStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#888888"))

				details := ""
				details += detailStyle.Render("CPU Usage: ") +
					lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%.1f%%", m.systemInfo.CPUUsage)) + "\n"

				gpuText := fmt.Sprintf("%.1f%%", m.systemInfo.GPUUsage)
				if m.systemInfo.GPUUsage == 0 {
					gpuText = "N/A"
				}
				details += detailStyle.Render("GPU Usage: ") +
					lipgloss.NewStyle().Bold(true).Render(gpuText) + "\n"

				content += details + "\n"
			}
		}
	} else {
		content += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("Loading system information...") + "\n\n"
	}

	// 配置信息
	configStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Border(lipgloss.RoundedBorder()).
		Padding(1)

	configContent := "Configuration:\n"
	configContent += "Prefix: " + m.config.PercentPrefix + "\n"
	configContent += "Suffix: " + m.config.PercentSuffix + "\n"
	configContent += fmt.Sprintf("Stress Threshold: %d%%\n", m.config.StressThreshold)
	configContent += fmt.Sprintf("Medium Threshold: %d%%", m.config.MediumThreshold)

	content += configStyle.Render(configContent) + "\n\n"

	// 帮助信息
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))
	content += helpStyle.Render("Press 'r' to refresh, 'q' to quit")

	return content
}

// updateBattery 更新电池信息
func (m *Model) updateBattery() tea.Cmd {
	return func() tea.Msg {
		info, err := battery.GetTouchpadBatteryInfo()
		if err != nil {
			return err
		}
		return info
	}
}

// updateSystemInfo 更新系统信息
func (m *Model) updateSystemInfo() tea.Cmd {
	return func() tea.Msg {
		info, err := system.GetSystemInfo()
		if err != nil {
			return err
		}
		return info
	}
}

// tick 定时器
func (m *Model) tick() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
