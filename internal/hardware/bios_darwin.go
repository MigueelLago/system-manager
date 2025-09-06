//go:build darwin

package hardware

import (
	"os/exec"
	"strings"
)

func readBiosInfo() string {
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Boot ROM Version:") || strings.Contains(line, "Firmware Version:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return "unknown"
}
