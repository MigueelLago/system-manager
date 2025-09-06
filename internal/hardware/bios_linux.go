//go:build linux

package hardware

import (
	"os/exec"
	"strings"
)

func readBiosInfo() string {
	cmd := exec.Command("dmidecode", "-s", "bios-version")
	out, err := cmd.Output()
	if err == nil {
		biosVersion := strings.TrimSpace(string(out))
		if biosVersion != "" {
			return biosVersion
		}
	}

	data, err := os.ReadFile("/sys/class/dmi/id/bios_version")
	if err == nil {
		return strings.TrimSpace(string(data))
	}
}
