package hardware

import (
	"os/exec"
	"runtime"
	"strings"
)

var systems = map[string]string{
	"linux":     "Linux",
	"windows":   "Windows",
	"darwin":    "macOS",
	"freebsd":   "FreeBSD",
	"openbsd":   "OpenBSD",
	"netbsd":    "NetBSD",
	"dragonfly": "DragonFly BSD",
	"solaris":   "Solaris",
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
	Arch    string `son:"arch"`
	//CPU     string `json:"cpu"`
	//BIOS    string `json:"bios"`
	//Board   string `json:"board"`
}

func getOs() string {
	os := runtime.GOOS

	for key, value := range systems {
		if key == os {
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
	os := getOs()

	switch os {
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

func GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:      getOs(),
		Version: getVersionOs(),
		Arch:    getArch(),
	}
}
