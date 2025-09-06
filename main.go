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
	fmt.Println("Placa-mãe:", systemInfo.Board)
	fmt.Println("BIOS:", systemInfo.Bios)
	fmt.Println("Memória RAM total (GB):", systemInfo.Memory.TotalGB)
	fmt.Println("Memória RAM disponível (GB):", systemInfo.Memory.AvailableGB)
	fmt.Println("Memória RAM usada (GB):", systemInfo.Memory.UsedGB)
}
