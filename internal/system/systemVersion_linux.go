//go:build linux

package system

import (
	"os/exec"
	"strings"
)

func readVersionOs() (string, error) {
	cmd := exec.Command("uname", "-sr")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	return "unknown", err
}
