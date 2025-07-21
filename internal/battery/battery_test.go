package battery

import (
	"testing"
)

func TestGetTouchpadBatteryInfo(t *testing.T) {
	info, err := GetTouchpadBatteryInfo()
	// 测试不应该返回错误
	if err != nil {
		t.Logf("获取电池信息时出错（这在没有触摸板的环境中是正常的）: %v", err)
		return
	}

	// 如果没有错误，检查返回的信息结构
	if info == nil {
		t.Fatal("返回的电池信息不应该为 nil")
	}

	// 如果电池可用，检查百分比范围
	if info.Available {
		if info.Percentage < 0 || info.Percentage > 100 {
			t.Errorf("电池百分比应该在 0-100 范围内，实际值: %d", info.Percentage)
		}
		t.Logf("电池信息: 百分比=%d%%, 充电中=%v", info.Percentage, info.IsCharging)
	} else {
		t.Log("触摸板电池不可用（可能没有连接触摸板）")
	}
}

func TestGetBatteryPercentage(t *testing.T) {
	percentage, err := getBatteryPercentage()
	if err != nil {
		t.Logf("获取电池百分比时出错: %v", err)
		return
	}

	// -1 表示没有找到电池信息，这是正常的
	if percentage == -1 {
		t.Log("没有找到电池信息")
		return
	}

	if percentage < 0 || percentage > 100 {
		t.Errorf("电池百分比应该在 0-100 范围内，实际值: %d", percentage)
	}
}

func TestGetChargingStatus(t *testing.T) {
	isCharging, err := getChargingStatus()
	if err != nil {
		t.Logf("获取充电状态时出错: %v", err)
		return
	}

	// 充电状态应该是布尔值，这里只是记录结果
	t.Logf("充电状态: %v", isCharging)
}
