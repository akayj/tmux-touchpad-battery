package system

// SystemInfo 表示系统信息，包括 CPU 和 GPU 使用率
type SystemInfo struct {
	CPUUsage  float64
	GPUUsage  float64
	Available bool
}
