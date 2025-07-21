package ui

import (
	"fmt"
	"time"

	"github.com/akayj/tmux-touchpad-battery/internal/battery"
	"github.com/akayj/tmux-touchpad-battery/internal/display"
	"github.com/akayj/tmux-touchpad-battery/internal/tmux"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model 表示 TUI 模型
type Model struct {
	batteryInfo *battery.BatteryInfo
	formatter   *display.Formatter
	config      *tmux.Config
	err         error
	quitting    bool
}

// tickMsg 定时更新消息
type tickMsg time.Time

// NewModel 创建新的 TUI 模型
func NewModel() *Model {
	config := tmux.GetConfig()
	formatter := display.NewFormatter(config)

	return &Model{
		config:    config,
		formatter: formatter,
	}
}

// Init 初始化模型
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.updateBattery(),
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
			return m, m.updateBattery()
		}

	case tickMsg:
		// 定时更新
		return m, tea.Batch(
			m.updateBattery(),
			m.tick(),
		)

	case *battery.BatteryInfo:
		m.batteryInfo = msg
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
		batteryDisplay := m.formatter.FormatBatteryWithStyle(m.batteryInfo)
		content += "Current Status: " + batteryDisplay + "\n\n"

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

// tick 定时器
func (m *Model) tick() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
