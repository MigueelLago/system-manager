package hardware

type CPU struct {
	ModelName   string `json:"model_name"`
	Cores       int    `json:"cores"`
	LogicsCores int    `json:"logical_cores"`
}

func GetCPUInfo() (CPU, error) {

	return readCPUInfo()
}
