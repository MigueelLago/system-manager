//go:build darwin

package hardware

import (
	"os/exec"
	"strings"
)

func readMotherBoard() string {
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "Model Name:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	cmd = exec.Command("ioreg", "-l", "-p", "IODeviceTree")
	out, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "model") && strings.Contains(line, "=") {
				parts := strings.Split(line, "=")
				if len(parts) > 1 {
					model := strings.Trim(strings.TrimSpace(parts[1]), "\"<>")
					if model != "" {
						return model
					}
				}
			}
		}
	}

	return "unknown"
}
