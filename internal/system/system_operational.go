package system

import "runtime"

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

func GetOS() string {
	osName := runtime.GOOS

	for key, value := range systems {
		if key == osName {
			return value
		}
	}

	return "unknown"
}

func GetArch() string {
	arch := runtime.GOARCH

	for key, value := range architectures {
		if key == arch {
			return value
		}
	}

	return "unknown"
}

func GetSystemVersion() string {
	return readVersionOs()
}
