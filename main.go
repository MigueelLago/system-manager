package main

import (
	"fmt"
	"system-manager/internal/hardware"
	"system-manager/internal/system"
)

func main() {

	memoryInfo, _ := hardware.GetMemoryInfo()
	cpuInfo, _ := hardware.GetCPUInfo()
	osName := system.GetOS()
	osArch := system.GetArch()
	osVersion := system.GetSystemVersion()
	motherBoard := hardware.GetMotherBoard()
	biosInfo := hardware.GetBIOSInfo()

	fmt.Println("Sistema operacional:", osName)
	fmt.Println("Arquitetura:", osArch)
	fmt.Println("Versão do SO:", osVersion)
	fmt.Println("Informações da CPU:", "Modelo:", cpuInfo.ModelName, ", Núcleos físicos:", cpuInfo.Cores, ","+
		" "+
		"Núcleos lógicos:",
		cpuInfo.LogicsCores)
	fmt.Println("Placa-mãe:", motherBoard)
	fmt.Println("BIOS:", biosInfo)
	fmt.Println("Memória RAM total (GB):", memoryInfo.TotalGB)
	fmt.Println("Memória RAM disponível (GB):", memoryInfo.AvailableGB)
	fmt.Println("Memória RAM usada (GB):", memoryInfo.UsedGB)
}
