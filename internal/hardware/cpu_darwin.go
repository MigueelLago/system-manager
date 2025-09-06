//go:build darwin

package hardware

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func readCPUInfo() (CPU, error) {
	cmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
	out, err := cmd.Output()
	modelName := "unknown"
	if err == nil {
		modelName = strings.TrimSpace(string(out))
	}

	// Get physical cores
	cmd = exec.Command("sysctl", "-n", "hw.physicalcpu")
	cores := runtime.NumCPU()
	if coresOut, err := cmd.Output(); err == nil {
		if coresInt, err := strconv.Atoi(strings.TrimSpace(string(coresOut))); err == nil {
			cores = coresInt
		}
	}

	// Get logical cores
	cmd = exec.Command("sysctl", "-n", "hw.logicalcpu")
	logicalCores := runtime.NumCPU()
	if logicalOut, err := cmd.Output(); err == nil {
		if logicalInt, err := strconv.Atoi(strings.TrimSpace(string(logicalOut))); err == nil {
			logicalCores = logicalInt
		}
	}

	return CPU{
		ModelName:   modelName,
		Cores:       cores,
		LogicsCores: logicalCores,
	}, nil
}
