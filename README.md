# Tmux Touchpad Battery (Golang 版本)

![normal](/screenshots/normal.jpg)
![charging](/screenshots/charging.jpg)

一个用 Golang 重构的 macOS 触摸板电池状态显示工具，使用 [Charm](https://charm.sh/) 的 Bubbles 库构建。

## 特性

- 🔋 显示触摸板电池电量百分比
- ⚡ 显示充电状态
- 🎨 可配置的颜色主题
- 🖥️ 交互式 TUI 界面
- 📊 实时状态监控
- 🔧 完全兼容原版 tmux 配置
- ⚠️ 低电量闪烁提醒功能

## 安装

### 通过 Tmux Plugin Manager (推荐)

在 `.tmux.conf` 文件中添加：

```bash
set -g @plugin 'akayj/tmux-touchpad-battery'
```

然后使用 `prefix + I` 安装插件。

### 手动安装

```bash
# 克隆仓库
git clone https://github.com/akayj/tmux-touchpad-battery.git
cd tmux-touchpad-battery

# 构建项目
make build

# 安装到系统路径（可选）
make install
```

## 使用方法

### 基本用法

在 tmux 状态栏中显示电池信息：

```bash
set -g status-right "#{touchpad_battery}"
```

### 命令行工具

```bash
# 显示当前电池状态
tmux-touchpad-battery -status

# 启动交互式 UI
tmux-touchpad-battery -ui

# 显示帮助信息
tmux-touchpad-battery -help

# 输出 tmux 格式（默认行为）
tmux-touchpad-battery
```

### 交互式 UI

运行 `make ui` 或 `tmux-touchpad-battery -ui` 启动交互式界面：

- 实时显示电池状态
- 显示配置信息
- 按 `r` 手动刷新
- 按 `q` 退出

## 配置选项

所有原版配置选项都得到支持：

| 选项                        | 默认值      | 说明                       |
| --------------------------- | ----------- | -------------------------- |
| `@tpb_percent_prefix`       | `Touchpad:` | 显示前缀                   |
| `@tpb_percent_suffix`       | `%`         | 显示后缀                   |
| `@tpb_color_charging`       | `green`     | 充电时颜色                 |
| `@tpb_color_high`           | `white`     | 高电量颜色                 |
| `@tpb_color_medium`         | `yellow`    | 中等电量颜色               |
| `@tpb_color_stress`         | `red`       | 低电量颜色                 |
| `@tpb_stress_threshold`     | `30`        | 低电量阈值                 |
| `@tpb_medium_threshold`     | `80`        | 中等电量阈值               |
| `@tpb_not_show_threshold`   | `100`       | 不显示阈值                 |
| `@tpb_blink_on_low_battery` | `off`       | 低电量时闪烁提醒（新功能） |

### 配置示例

```bash
# 自定义前缀和后缀
set -g @tpb_percent_prefix "🖱️ "
set -g @tpb_percent_suffix "%%"

# 自定义颜色
set -g @tpb_color_charging "#00ff00"
set -g @tpb_color_high "#ffffff"
set -g @tpb_color_medium "#ffff00"
set -g @tpb_color_stress "#ff0000"

# 自定义阈值
set -g @tpb_stress_threshold "20"
set -g @tpb_medium_threshold "70"
set -g @tpb_not_show_threshold "95"

# 启用低电量闪烁提醒
set -g @tpb_blink_on_low_battery "on"
```

## 开发

### 项目结构

```
.
├── cmd/tmux-touchpad-battery/    # 主程序入口
├── internal/
│   ├── battery/                  # 电池状态检测
│   ├── display/                  # 格式化和显示
│   ├── tmux/                     # tmux 配置读取
│   └── ui/                       # TUI 界面
├── scripts/                      # 原版 bash 脚本（保留）
├── screenshots/                  # 截图
├── Makefile                      # 构建脚本
├── go.mod                        # Go 模块定义
└── README_GO.md                  # 本文档
```

### 可用命令

```bash
make build      # 构建项目
make install    # 安装到系统
make ui         # 启动交互式界面
make status     # 显示电池状态
make test       # 运行测试
make clean      # 清理构建文件
make help       # 显示程序帮助
make info       # 显示项目信息
```

### 依赖

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI 框架
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI 组件
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - 样式库

## 技术特点

### 相比原版的改进

1. **性能提升**: Golang 编译后的二进制文件执行速度更快
2. **更好的错误处理**: 完善的错误处理和日志记录
3. **交互式界面**: 使用 Bubble Tea 构建的现代 TUI
4. **代码结构**: 模块化设计，易于维护和扩展
5. **类型安全**: Golang 的静态类型检查
6. **并发支持**: 为未来的功能扩展提供并发能力

### 兼容性

- ✅ 完全兼容原版 tmux 配置
- ✅ 支持所有原版功能
- ✅ 相同的输出格式
- ✅ macOS 10.12+ 支持

## 故障排除

### 常见问题

1. **触摸板未检测到**

   ```bash
   # 检查触摸板连接
   tmux-touchpad-battery -status
   ```

2. **权限问题**

   ```bash
   # 确保有执行权限
   chmod +x bin/tmux-touchpad-battery
   ```

3. **构建失败**
   ```bash
   # 更新依赖
   make deps
   ```

### 调试

启用详细输出：

```bash
# 查看详细状态
tmux-touchpad-battery -status

# 启动 UI 进行实时监控
tmux-touchpad-battery -ui
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 致谢

- 原版作者的优秀设计
- [Charm](https://charm.sh/) 团队的出色工具
- Tmux 社区的支持
