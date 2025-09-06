//go:build linux

package hardware

func readMemoryInfo() (MemoryInfo, error) {

	data, err := os.ReadFile("/proc/meminfo")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		var totalKB, freeKB, availableKB, buffersKB, cachedKB int64

		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				value, err := strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					continue
				}

				switch {
				case strings.HasPrefix(line, "MemTotal:"):
					totalKB = value
				case strings.HasPrefix(line, "MemFree:"):
					freeKB = value
				case strings.HasPrefix(line, "MemAvailable:"):
					availableKB = value
				case strings.HasPrefix(line, "Buffers:"):
					buffersKB = value
				case strings.HasPrefix(line, "Cached:"):
					cachedKB = value
				}
			}
		}

		totalGB := float64(totalKB) / 1024 / 1024

		var availableGB float64
		if availableKB > 0 {
			availableGB = float64(availableKB) / 1024 / 1024
		} else {
			availableGB = float64(freeKB+buffersKB+cachedKB) / 1024 / 1024
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
}
