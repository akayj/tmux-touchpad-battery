package battery

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// BatteryInfo 表示电池信息
type BatteryInfo struct {
	Percentage int
	IsCharging bool
	Available  bool
}

// GetTouchpadBatteryInfo 获取触摸板电池信息
func GetTouchpadBatteryInfo() (*BatteryInfo, error) {
	info := &BatteryInfo{}

	// 获取电池百分比
	percentage, err := getBatteryPercentage()
	if err != nil {
		return nil, err
	}

	if percentage == -1 {
		info.Available = false
		return info, nil
	}

	info.Percentage = percentage
	info.Available = true

	// 获取充电状态
	isCharging, err := getChargingStatus()
	if err != nil {
		return nil, err
	}

	info.IsCharging = isCharging
	return info, nil
}

// getBatteryPercentage 获取电池百分比
func getBatteryPercentage() (int, error) {
	cmd := exec.Command("ioreg", "-l")
	output, err := cmd.Output()
	if err != nil {
		return -1, err
	}

	// 查找 BatteryPercent
	re := regexp.MustCompile(`"BatteryPercent"\s*=\s*(\d+)`)
	matches := re.FindStringSubmatch(string(output))

	if len(matches) < 2 {
		return -1, nil // 没有找到电池信息
	}

	percentage, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, err
	}

	return percentage, nil
}

// getChargingStatus 获取充电状态
func getChargingStatus() (bool, error) {
	cmd := exec.Command("ioreg", "-l")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	// 查找 BatteryStatusFlags
	re := regexp.MustCompile(`"BatteryStatusFlags"\s*=\s*(\d+)`)
	matches := re.FindStringSubmatch(string(output))

	if len(matches) < 2 {
		return false, nil
	}

	status := strings.TrimSpace(matches[1])
	return status == "3", nil
}
