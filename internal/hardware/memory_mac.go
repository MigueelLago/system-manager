//go:build darwin

package hardware

import (
	"os/exec"
	"strconv"
	"strings"
)

func readMemoryInfo() (MemoryInfo, error) {

	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	out, err := cmd.Output()
	var totalGB float64
	if err == nil {
		if totalBytes, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64); err == nil {
			totalGB = float64(totalBytes) / 1024 / 1024 / 1024
		}
	}

	cmd = exec.Command("vm_stat")
	out, err = cmd.Output()
	var freePages, inactivePages, speculativePages, wiredPages, activePages int64
	var pageSize int64 = 4096 // Default page size

	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.Contains(line, "page size of") {
				parts := strings.Split(line, " ")
				for i, part := range parts {
					if part == "of" && i+1 < len(parts) {
						if size, err := strconv.ParseInt(parts[i+1], 10, 64); err == nil {
							pageSize = size
						}
						break
					}
				}
				continue
			}

			if strings.Contains(line, ":") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					valueStr := strings.TrimSpace(strings.Replace(parts[1], ".", "", -1))
					if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
						switch {
						case strings.HasPrefix(line, "Pages free:"):
							freePages = value
						case strings.HasPrefix(line, "Pages inactive:"):
							inactivePages = value
						case strings.HasPrefix(line, "Pages speculative:"):
							speculativePages = value
						case strings.HasPrefix(line, "Pages wired down:"):
							wiredPages = value
						case strings.HasPrefix(line, "Pages active:"):
							activePages = value
						}
					}
				}
			}
		}
	}

	availableBytes := (freePages + inactivePages + speculativePages) * pageSize
	availableGB := float64(availableBytes) / 1024 / 1024 / 1024

	usedBytes := (wiredPages + activePages) * pageSize
	usedGB := float64(usedBytes) / 1024 / 1024 / 1024

	if availableGB+usedGB > totalGB*1.1 || availableGB < 0 {
		usedGB = totalGB - availableGB
	}

	totalGB = float64(int(totalGB*100+0.5)) / 100
	availableGB = float64(int(availableGB*100+0.5)) / 100
	usedGB = float64(int(usedGB*100+0.5)) / 100

	return MemoryInfo{
		TotalGB:     totalGB,
		AvailableGB: availableGB,
		UsedGB:      usedGB,
	}, nil
}
