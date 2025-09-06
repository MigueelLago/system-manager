//go:build windows

package hardware

func readCPUInfo() (CPU, error) {
	cmd := exec.Command("wmic", "cpu", "get", "Name,NumberOfCores,NumberOfLogicalProcessors", "/format:csv")
	out, err := cmd.Output()
	if err != nil {
		return CPU{
			ModelName:   "unknown",
			Cores:       runtime.NumCPU(),
			LogicsCores: runtime.NumCPU(),
		}, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "Node,Name,NumberOfCores") {
			fields := strings.Split(line, ",")
			if len(fields) >= 4 {
				modelName := strings.TrimSpace(fields[1])
				coresStr := strings.TrimSpace(fields[2])
				logicalCoresStr := strings.TrimSpace(fields[3])

				cores, err1 := strconv.Atoi(coresStr)
				logicalCores, err2 := strconv.Atoi(logicalCoresStr)

				if err1 != nil || err2 != nil {
					cores = runtime.NumCPU()
					logicalCores = runtime.NumCPU()
				}

				return CPU{
					ModelName:   modelName,
					Cores:       cores,
					LogicsCores: logicalCores,
				}, nil
			}
		}
	}
}
