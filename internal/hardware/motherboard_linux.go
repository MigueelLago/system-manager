//go:build linux

package hardware

func readMotherBoard() string {
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
}
