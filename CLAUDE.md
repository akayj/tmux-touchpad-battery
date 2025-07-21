# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码仓库中工作时提供指导。

## 项目概述

这是一个基于 Go 构建的 macOS 触摸板电池状态显示工具，从 bash 版本重写而来。它在 tmux 状态栏中显示触摸板电池百分比，并使用 Charm 的 Bubble Tea 框架提供交互式 TUI 界面。

## 架构

### 核心组件

- **cmd/tmux-touchpad-battery/main.go**: 入口点，处理不同模式的 CLI 标志 (-status, -ui, -help)
- **internal/battery/**: 使用 macOS `ioreg` 命令解析 BatteryPercent 和 BatteryStatusFlags 进行电池检测
- **internal/tmux/**: 配置读取器，执行 `tmux show-option` 命令获取用户设置
- **internal/display/**: 格式化器，基于电池电量和充电状态应用 tmux 颜色代码和样式
- **internal/ui/**: 使用 Bubble Tea 构建的交互式 TUI，用于实时监控

### 数据流

1. **配置加载**: `tmux.GetConfig()` → 通过 `tmux show-option -gqv` 命令读取用户设置
2. **电池检测**: `battery.GetTouchpadBatteryInfo()` → 使用正则表达式解析 `ioreg -l` 输出
3. **显示格式化**: `display.Formatter` → 基于电池电量和充电状态应用颜色和样式
4. **输出路由**: 主程序路由到 tmux 格式（静默）、状态显示（详细）或交互式 UI

### 模块依赖关系

```
main.go (入口点)
    ↓
display.Formatter (格式化层)
├─ battery.BatteryInfo (数据层)
└─ tmux.Config (配置层)
    ↑
ui.Model (UI 层)
```

## 开发命令

### 构建和运行
```bash
make build          # 构建二进制文件到 bin/tmux-touchpad-battery
make install        # 安装到 /usr/local/bin (需要 sudo)
make clean          # 删除 bin/ 目录
```

### 测试和开发
```bash
make test           # 使用 go test ./... 运行 Go 测试
make fmt            # 使用 go fmt ./... 格式化代码
make lint           # 使用 go vet ./... 检查代码
make deps           # 使用 go mod tidy 更新依赖

# 运行特定测试
go test ./internal/battery -v              # 测试电池检测
go test ./internal/tmux -v                 # 测试 tmux 配置
go test ./internal/display -v              # 测试显示格式化
```

### 应用模式
```bash
make status         # 在终端中显示电池状态
make ui             # 启动交互式 TUI
make help           # 显示应用帮助
make dev            # 构建并运行状态检查
```

## 配置系统

应用程序读取用户可以在 `.tmux.conf` 中设置的 tmux 配置选项：

- `@tpb_percent_prefix/suffix`: 百分比周围的显示文本
- `@tpb_color_*`: 不同电池状态的颜色（充电、高电量、中等电量、压力状态）
- `@tpb_*_threshold`: 颜色变化的电池电量阈值
- `@tpb_blink_on_low_battery`: 低电量时启用闪烁
- `@tpb_charging_icon`: 充电时显示的图标（默认：'⚡'）
- `@tpb_show_charging_icon`: 切换充电图标显示（默认：'on'）

所有选项在 internal/tmux/config.go:25-40 中都有合理的默认值。配置通过 `tmux show-option -gqv` 命令读取，具有优雅的默认值回退机制。

## 关键技术细节

### 电池检测
使用 macOS 特定的 `ioreg -l` 命令解析系统注册表中的触摸板电池信息。搜索具有 "BatteryPercent" 和 "BatteryStatusFlags" 条目的设备，使用正则表达式模式。充电状态由 BatteryStatusFlags 值 "3" 确定。

### Tmux 集成
执行 `tmux show-option -gqv` 读取用户配置，具有优雅的默认值回退机制。工具输出与 tmux 兼容的颜色代码用于状态栏集成。支持双输出格式：tmux 格式（带颜色代码）和终端格式（带 lipgloss 样式）。

### 错误处理
- **静默失败** 用于 tmux 输出模式（无错误消息以避免破坏状态栏）
- **详细错误报告** 用于 -status 和 -ui 模式
- **优雅降级** 当 tmux 选项缺失或触摸板不可用时
- **配置验证** 具有不同选项格式的类型转换

## 依赖项

- **Bubble Tea**: 交互模式的 TUI 框架
- **Lip Gloss**: 终端输出的样式库
- 标准 Go 库用于系统命令执行

项目与原始 bash 版本的 tmux 配置选项和输出格式保持完全向后兼容。