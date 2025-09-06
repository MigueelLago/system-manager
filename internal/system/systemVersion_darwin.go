//go:build darwin

package system

import (
	"os/exec"
	"strings"
)

func readVersionOs() string {
	cmd := exec.Command("sw_vers", "-ProductVersion")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	return "unknown"
}
