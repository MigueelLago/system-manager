//go:build windows

package hardware

func readMemoryInfo() (MemoryInfo, error) {
	cmd := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory", "/format:csv")
	out, err := cmd.Output()
	var totalGB float64
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if !strings.Contains(line, "Node,TotalPhysicalMemory") && strings.TrimSpace(line) != "" {
				fields := strings.Split(line, ",")
				if len(fields) >= 2 {
					if totalBytes, err := strconv.ParseInt(strings.TrimSpace(fields[1]), 10, 64); err == nil {
						totalGB = float64(totalBytes) / 1024 / 1024 / 1024
						break
					}
				}
			}
		}
	}

	cmd = exec.Command("wmic", "OS", "get", "FreePhysicalMemory", "/format:csv")
	out, err = cmd.Output()
	var availableGB float64
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if !strings.Contains(line, "Node,FreePhysicalMemory") && strings.TrimSpace(line) != "" {
				fields := strings.Split(line, ",")
				if len(fields) >= 2 {
					if freeKB, err := strconv.ParseInt(strings.TrimSpace(fields[1]), 10, 64); err == nil {
						availableGB = float64(freeKB) / 1024 / 1024
						break
					}
				}
			}
		}
	}

	usedGB := totalGB - availableGB

	totalGB = float64(int(totalGB*100+0.5)) / 100
	availableGB = float64(int(availableGB*100+0.5)) / 100
	usedGB = float64(int(usedGB*100+0.5)) / 100

	return MemoryInfo{
		TotalGB:     totalGB,
		AvailableGB: availableGB,
		UsedGB:      usedGB,
	}, nil
}
