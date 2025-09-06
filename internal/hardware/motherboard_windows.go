//go:build windows

package hardware

func readMotherBoard() string {
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
}
