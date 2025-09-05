package hardware

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var systems = map[string]string{
	"linux":   "Linux",
	"windows": "Windows",
	"darwin":  "macOS",
}

var architectures = map[string]string{
	"amd64": "x86_64",
	"386":   "x86",
	"arm":   "ARM",
	"arm64": "ARM64",
}

type CPU struct {
	ModelName   string `json:"model_name"`
	Cores       int    `json:"cores"`
	LogicsCores int    `json:"logical_cores"`
}

type SystemInfo struct {
	OS      string `json:"os"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
	CPU     CPU    `json:"cpu"`
	//BIOS    string `json:"bios"`
	//Board   string `json:"board"`
}

func getOs() string {
	osName := runtime.GOOS

	for key, value := range systems {
		if key == osName {
			return value
		}
	}

	return "unknown"
}

func getArch() string {
	arch := runtime.GOARCH

	for key, value := range architectures {
		if key == arch {
			return value
		}
	}

	return "unknown"
}

func getVersionOs() string {

	// Check the OS
	osType := getOs()

	switch osType {
	case "Linux":
		cmd := exec.Command("uname", "-sr")
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output))
		}
	case "Windows":
		cmd := exec.Command("cmd", "ver")
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output))
		}
	case "macOS":
		cmd := exec.Command("sw_vers", "-ProductVersion")
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output))
		}
	default:
		cmd := exec.Command("uname", "-sr")
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output))
		}
	}

	return "unknown"
}

func getCPUInfo() CPU {
	// Check the OS
	system := getOs()

	switch system {
	case "Linux":
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
		}

	case "Windows":
		cmd := exec.Command("wmic", "cpu", "get", "Name,NumberOfCores,NumberOfLogicalProcessors", "/format:csv")
		out, err := cmd.Output()
		if err != nil {
			return CPU{
				ModelName:   "unknown",
				Cores:       runtime.NumCPU(),
				LogicsCores: runtime.NumCPU(),
			}
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
					}
				}
			}
		}

	case "macOS":
		// Get CPU model name
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
		}

	default:
		return CPU{
			ModelName:   "unknown",
			Cores:       runtime.NumCPU(),
			LogicsCores: runtime.NumCPU(),
		}
	}

	return CPU{
		ModelName:   "unknown",
		Cores:       runtime.NumCPU(),
		LogicsCores: runtime.NumCPU(),
	}
}

func GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:      getOs(),
		Version: getVersionOs(),
		Arch:    getArch(),
		CPU:     getCPUInfo(),
	}
}
