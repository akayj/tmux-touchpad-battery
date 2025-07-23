package display

// Formatter 定义了信息格式化的通用接口
type Formatter interface {
	// Format 格式化信息为 tmux 状态栏显示
	Format() string

	// FormatWithStyle 使用样式格式化信息（用于终端显示）
	FormatWithStyle() string
}
