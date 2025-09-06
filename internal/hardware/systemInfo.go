package hardware

import (
	"os"
	"os/exec"
	"runtime"
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

type SystemInfo struct {
	OS      string `json:"os"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
	Board   string `json:"board"`
	Bios    string `json:"bios"`
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

func GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:      getOs(),
		Version: getVersionOs(),
		Arch:    getArch(),
		Board:   getBoardInfo(),
		Bios:    getBIOSInfo(),
	}
}
