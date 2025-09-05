package main

import (
	"fmt"
	"system-manager/internal/hardware"
)

func main() {
	systemInfo := hardware.GetSystemInfo()

	fmt.Println("Sistema operacional:", systemInfo.OS)
	fmt.Println("Arquitetura:", systemInfo.Arch)
	fmt.Println("Versão do SO:", systemInfo.Version)
	fmt.Println("Informações da CPU:", "Modelo:", systemInfo.CPU.ModelName, ", Núcleos físicos:", systemInfo.CPU.Cores, ","+
		" "+
		"Núcleos lógicos:",
		systemInfo.CPU.LogicsCores)
}
