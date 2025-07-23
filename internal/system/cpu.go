package system

import (
	"os/exec"
	"regexp"
	"strconv"
)

// GetCPUUsage 获取 CPU 使用率
func GetCPUUsage() (float64, error) {
	// 使用 top 命令获取 CPU 使用率
	cmd := exec.Command("top", "-l", "1", "-n", "0")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// 解析输出以获取 CPU 使用率
	// 查找类似 "CPU usage: 15.65% user, 18.39% sys, 65.94% idle" 的行
	re := regexp.MustCompile(`CPU usage: ([0-9.]+)% user, ([0-9.]+)% sys`)
	matches := re.FindStringSubmatch(string(output))

	if len(matches) < 3 {
		return 0, nil // 没有找到 CPU 使用率信息
	}

	// 计算总使用率 (用户 + 系统)
	user, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}

	sys, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return 0, err
	}

	return user + sys, nil
}

// GetSystemInfo 获取系统信息，包括 CPU 和 GPU 使用率
func GetSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{}

	// 获取 CPU 使用率
	cpuUsage, err := GetCPUUsage()
	if err != nil {
		return nil, err
	}

	info.CPUUsage = cpuUsage

	// 获取 GPU 使用率
	gpuUsage, err := GetGPUUsage()
	if err != nil {
		// 如果获取 GPU 使用率失败，只记录错误但不中断
		// GPU 使用率将保持为 0
		info.GPUUsage = 0
	} else {
		info.GPUUsage = gpuUsage
	}

	info.Available = true

	return info, nil
}
