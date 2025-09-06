package hardware

type MemoryInfo struct {
	TotalGB     float64
	AvailableGB float64
	UsedGB      float64
}

func GetMemoryInfo() (MemoryInfo, error) {
	return readMemoryInfo()
}
