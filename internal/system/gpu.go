package system

import (
	"os/exec"
	"regexp"
	"strconv"
)

// GetGPUUsage 获取 GPU 使用率
// 注意：在 macOS 上获取 GPU 使用率需要 root 权限，因此这个函数可能会返回错误
func GetGPUUsage() (float64, error) {
	// 尝试使用 powermetrics 获取 GPU 使用率（需要 root 权限）
	// 如果没有权限，返回 0 和 nil
	cmd := exec.Command("powermetrics", "--samplers", "gpu_power", "--show-all", "--sample-count", "1")
	output, err := cmd.Output()
	if err != nil {
		// 如果没有权限或其他错误，返回 0
		return 0, nil
	}

	// 解析输出以获取 GPU 使用率
	// 查找类似 "GPU Power: 0.046123 W (100.0%)" 的行
	re := regexp.MustCompile(`GPU Power: [0-9.]+ W \(([0-9.]+)%\)`)
	matches := re.FindStringSubmatch(string(output))

	if len(matches) < 2 {
		return 0, nil // 没有找到 GPU 使用率信息
	}

	gpuUsage, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, nil
	}

	return gpuUsage, nil
}
