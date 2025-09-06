//go:build linux

package hardware

func readCPUInfo() (CPU, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return CPU{
			ModelName:   "unknown",
			Cores:       runtime.NumCPU(),
			LogicsCores: runtime.NumCPU(),
		}
	}
	cpuInfo := string(data)
	lines := strings.Split(cpuInfo, "\n")

	var modelName string
	physicalCores := 0
	processors := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "model name") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				modelName = strings.TrimSpace(parts[1])
			}
		}
		if strings.HasPrefix(line, "physical id") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				physicalID := strings.TrimSpace(parts[1])
				processors[physicalID] = true
			}
		}
	}

	physicalCores = len(processors)
	if physicalCores == 0 {
		physicalCores = runtime.NumCPU() / 2
	}

	if modelName == "" {
		modelName = "unknown"
	}

	return CPU{
		ModelName:   modelName,
		Cores:       physicalCores,
		LogicsCores: runtime.NumCPU(),
	}, nil
}
