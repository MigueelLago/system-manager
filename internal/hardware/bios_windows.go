//go:build windows

package hardware

import (
	"os/exec"
	"strings"
)

func readBiosInfo() string {
	cmd := exec.Command("wmic", "bios", "get", "SMBIOSBIOSVersion", "/format:csv")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.Contains(line, "Node,SMBIOSBIOSVersion") {
				fields := strings.Split(line, ",")
				if len(fields) >= 2 && strings.TrimSpace(fields[1]) != "" {
					return strings.TrimSpace(fields[1])
				}
			}
		}
	}
}
