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

type MemoryInfo struct {
	TotalGB     float64 `json:"total_gb"`
	AvailableGB float64 `json:"available_gb"`
	UsedGB      float64 `json:"used_gb"`
}

type SystemInfo struct {
	OS      string     `json:"os"`
	Version string     `json:"version"`
	Arch    string     `json:"arch"`
	CPU     CPU        `json:"cpu"`
	Board   string     `json:"board"`
	Bios    string     `json:"bios"`
	Memory  MemoryInfo `json:"memory"`
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

func getBoardInfo() string {
	system := getOs()

	switch system {
	case "Linux":
		cmd := exec.Command("dmidecode", "-s", "baseboard-product-name")
		out, err := cmd.Output()
		if err == nil {
			boardName := strings.TrimSpace(string(out))
			if boardName != "" && boardName != "System Product Name" {
				return boardName
			}
		}

		data, err := os.ReadFile("/sys/class/dmi/id/board_name")
		if err == nil {
			boardName := strings.TrimSpace(string(data))
			if boardName != "" {
				return boardName
			}
		}

		vendor, err1 := os.ReadFile("/sys/class/dmi/id/board_vendor")
		name, err2 := os.ReadFile("/sys/class/dmi/id/board_name")
		if err1 == nil && err2 == nil {
			vendorStr := strings.TrimSpace(string(vendor))
			nameStr := strings.TrimSpace(string(name))
			if vendorStr != "" && nameStr != "" {
				return vendorStr + " " + nameStr
			}
		}

	case "Windows":
		cmd := exec.Command("wmic", "baseboard", "get", "Manufacturer,Product", "/format:csv")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.Contains(line, "Node,Manufacturer,Product") {
					fields := strings.Split(line, ",")
					if len(fields) >= 3 {
						manufacturer := strings.TrimSpace(fields[1])
						product := strings.TrimSpace(fields[2])
						if manufacturer != "" && product != "" {
							return manufacturer + " " + product
						}
					}
				}
			}
		}

	case "macOS":
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
	}

	return "unknown"
}

func getBIOSInfo() string {
	system := getOs()

	switch system {
	case "Linux":
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

	case "Windows":
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

	case "macOS":
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
	}

	return "unknown"
}

func getMemoryInfo() MemoryInfo {
	system := getOs()

	switch system {
	case "Linux":
		data, err := os.ReadFile("/proc/meminfo")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			var totalKB, freeKB, availableKB, buffersKB, cachedKB int64

			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					value, err := strconv.ParseInt(fields[1], 10, 64)
					if err != nil {
						continue
					}

					switch {
					case strings.HasPrefix(line, "MemTotal:"):
						totalKB = value
					case strings.HasPrefix(line, "MemFree:"):
						freeKB = value
					case strings.HasPrefix(line, "MemAvailable:"):
						availableKB = value
					case strings.HasPrefix(line, "Buffers:"):
						buffersKB = value
					case strings.HasPrefix(line, "Cached:"):
						cachedKB = value
					}
				}
			}

			totalGB := float64(totalKB) / 1024 / 1024

			var availableGB float64
			if availableKB > 0 {
				availableGB = float64(availableKB) / 1024 / 1024
			} else {
				availableGB = float64(freeKB+buffersKB+cachedKB) / 1024 / 1024
			}

			usedGB := totalGB - availableGB

			totalGB = float64(int(totalGB*100+0.5)) / 100
			availableGB = float64(int(availableGB*100+0.5)) / 100
			usedGB = float64(int(usedGB*100+0.5)) / 100

			return MemoryInfo{
				TotalGB:     totalGB,
				AvailableGB: availableGB,
				UsedGB:      usedGB,
			}
		}

	case "Windows":
		cmd := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory", "/format:csv")
		out, err := cmd.Output()
		var totalGB float64
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if !strings.Contains(line, "Node,TotalPhysicalMemory") && strings.TrimSpace(line) != "" {
					fields := strings.Split(line, ",")
					if len(fields) >= 2 {
						if totalBytes, err := strconv.ParseInt(strings.TrimSpace(fields[1]), 10, 64); err == nil {
							totalGB = float64(totalBytes) / 1024 / 1024 / 1024
							break
						}
					}
				}
			}
		}

		cmd = exec.Command("wmic", "OS", "get", "FreePhysicalMemory", "/format:csv")
		out, err = cmd.Output()
		var availableGB float64
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if !strings.Contains(line, "Node,FreePhysicalMemory") && strings.TrimSpace(line) != "" {
					fields := strings.Split(line, ",")
					if len(fields) >= 2 {
						if freeKB, err := strconv.ParseInt(strings.TrimSpace(fields[1]), 10, 64); err == nil {
							availableGB = float64(freeKB) / 1024 / 1024
							break
						}
					}
				}
			}
		}

		usedGB := totalGB - availableGB

		totalGB = float64(int(totalGB*100+0.5)) / 100
		availableGB = float64(int(availableGB*100+0.5)) / 100
		usedGB = float64(int(usedGB*100+0.5)) / 100

		return MemoryInfo{
			TotalGB:     totalGB,
			AvailableGB: availableGB,
			UsedGB:      usedGB,
		}

	case "macOS":
		cmd := exec.Command("sysctl", "-n", "hw.memsize")
		out, err := cmd.Output()
		var totalGB float64
		if err == nil {
			if totalBytes, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64); err == nil {
				totalGB = float64(totalBytes) / 1024 / 1024 / 1024
			}
		}

		cmd = exec.Command("vm_stat")
		out, err = cmd.Output()
		var freePages, inactivePages, speculativePages, wiredPages, activePages int64
		var pageSize int64 = 4096 // Default page size

		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)

				if strings.Contains(line, "page size of") {
					parts := strings.Split(line, " ")
					for i, part := range parts {
						if part == "of" && i+1 < len(parts) {
							if size, err := strconv.ParseInt(parts[i+1], 10, 64); err == nil {
								pageSize = size
							}
							break
						}
					}
					continue
				}

				if strings.Contains(line, ":") {
					parts := strings.Split(line, ":")
					if len(parts) == 2 {
						valueStr := strings.TrimSpace(strings.Replace(parts[1], ".", "", -1))
						if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
							switch {
							case strings.HasPrefix(line, "Pages free:"):
								freePages = value
							case strings.HasPrefix(line, "Pages inactive:"):
								inactivePages = value
							case strings.HasPrefix(line, "Pages speculative:"):
								speculativePages = value
							case strings.HasPrefix(line, "Pages wired down:"):
								wiredPages = value
							case strings.HasPrefix(line, "Pages active:"):
								activePages = value
							}
						}
					}
				}
			}
		}

		availableBytes := (freePages + inactivePages + speculativePages) * pageSize
		availableGB := float64(availableBytes) / 1024 / 1024 / 1024

		usedBytes := (wiredPages + activePages) * pageSize
		usedGB := float64(usedBytes) / 1024 / 1024 / 1024

		if availableGB+usedGB > totalGB*1.1 || availableGB < 0 {
			usedGB = totalGB - availableGB
		}

		totalGB = float64(int(totalGB*100+0.5)) / 100
		availableGB = float64(int(availableGB*100+0.5)) / 100
		usedGB = float64(int(usedGB*100+0.5)) / 100

		return MemoryInfo{
			TotalGB:     totalGB,
			AvailableGB: availableGB,
			UsedGB:      usedGB,
		}
	}

	return MemoryInfo{
		TotalGB:     0,
		AvailableGB: 0,
		UsedGB:      0,
	}
}

func GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:      getOs(),
		Version: getVersionOs(),
		Arch:    getArch(),
		CPU:     getCPUInfo(),
		Board:   getBoardInfo(),
		Bios:    getBIOSInfo(),
		Memory:  getMemoryInfo(),
	}
}
