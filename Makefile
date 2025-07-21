.PHONY: build install clean test ui status help

# 默认目标
all: build

# 构建二进制文件
build:
	@echo "构建 tmux-touchpad-battery..."
	@mkdir -p bin
	@go build -o bin/tmux-touchpad-battery ./cmd/tmux-touchpad-battery
	@echo "构建完成: bin/tmux-touchpad-battery"

# 安装到系统路径
install: build
	@echo "安装 tmux-touchpad-battery 到 /usr/local/bin..."
	@sudo cp bin/tmux-touchpad-battery /usr/local/bin/
	@sudo chmod +x /usr/local/bin/tmux-touchpad-battery
	@echo "安装完成"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@echo "清理完成"

# 运行测试
test:
	@echo "运行测试..."
	@go test ./...

# 显示交互式 UI
ui: build
	@./bin/tmux-touchpad-battery -ui

# 显示电池状态
status: build
	@./bin/tmux-touchpad-battery -status

# 显示帮助
help: build
	@./bin/tmux-touchpad-battery -help

# 下载依赖
deps:
	@echo "下载依赖..."
	@go mod tidy
	@go mod download
	@echo "依赖下载完成"

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@echo "格式化完成"

# 检查代码
lint:
	@echo "检查代码..."
	@go vet ./...
	@echo "检查完成"

# 开发模式 - 构建并运行状态检查
dev: build status

# 显示项目信息
info:
	@echo "Tmux Touchpad Battery - Golang 版本"
	@echo "=================================="
	@echo "项目目录: $(PWD)"
	@echo "Go 版本: $(shell go version)"
	@echo "模块: $(shell head -1 go.mod | cut -d' ' -f2)"
	@echo ""
	@echo "可用命令:"
	@echo "  make build   - 构建项目"
	@echo "  make install - 安装到系统"
	@echo "  make ui      - 显示交互式界面"
	@echo "  make status  - 显示电池状态"
	@echo "  make test    - 运行测试"
	@echo "  make clean   - 清理构建文件"
	@echo "  make help    - 显示程序帮助"