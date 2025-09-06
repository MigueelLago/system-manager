# System Manager

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey?style=for-the-badge)](https://github.com/miguellago/system-manager)

## 📋 Descrição

**System Manager** é uma ferramenta multiplataforma desenvolvida em Go que permite coletar e exibir informações detalhadas de hardware e sistema operacional. O projeto oferece compatibilidade total com **Windows**, **macOS** e **Linux**, fornecendo uma interface unificada para acessar dados críticos do sistema.

> 🎯 **Objetivo Educacional**: Este projeto foi desenvolvido como parte de uma iniciativa de aprendizado para aprofundar conhecimentos na linguagem Go e explorar sua biblioteca padrão (standard library), especialmente no que se refere à interação com o sistema operacional e coleta de informações de hardware.

## ✨ Funcionalidades Atuais

### 🖥️ Informações de Sistema
- **Sistema Operacional**: Detecção automática (Windows/macOS/Linux)
- **Arquitetura**: x86_64, x86, ARM, ARM64
- **Versão do SO**: Versão completa do sistema operacional

### 🔧 Informações de Hardware
- **CPU**:
  - Modelo e fabricante do processador
  - Número de núcleos físicos
  - Número de núcleos lógicos (threads)
- **Memória RAM**:
  - Memória total (GB)
  - Memória disponível (GB)
  - Memória em uso (GB)
- **Placa-mãe**: Fabricante e modelo
- **BIOS/UEFI**: Versão do firmware

## 🚀 Funcionalidades Planejadas

### 📡 Rede
- [ ] Interfaces de rede ativas
- [ ] Endereços IP (IPv4/IPv6)
- [ ] Status de conectividade
- [ ] Estatísticas de tráfego

### 💾 Armazenamento
- [ ] Discos e partições
- [ ] Espaço total, usado e disponível
- [ ] Tipo de sistema de arquivos
- [ ] Estado de saúde dos discos (S.M.A.R.T.)

### 🔋 Gerenciamento de Energia
- [ ] Status da bateria
- [ ] Nível de carga atual
- [ ] **Controle de carregamento**: Limitar carregamento em porcentagem específica
- [ ] Modo de economia de energia

## 🛠️ Instalação

### Pré-requisitos
- Go 1.25 ou superior
- Privilégios administrativos (para algumas funcionalidades de hardware)

### Clonando o Repositório
```bash
git clone https://github.com/miguellago/system-manager.git
cd system-manager
```

### Executando
```bash
go run main.go
```

### Compilando
```bash
# Para o sistema atual
go build -o system-manager

# Para Windows
GOOS=windows GOARCH=amd64 go build -o system-manager.exe

# Para Linux
GOOS=linux GOARCH=amd64 go build -o system-manager-linux

# Para macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o system-manager-macos-intel

# Para macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o system-manager-macos-arm
```

## 📖 Exemplo de Uso

```bash
$ go run main.go

Sistema operacional: macOS
Arquitetura: ARM64
Versão do SO: 15.6.1
Informações da CPU: Modelo: Apple M2 , Núcleos físicos: 8 , Núcleos lógicos: 8
Placa-mãe: MacBook Air
BIOS: 11881.140.96
Memória RAM total (GB): 8
Memória RAM disponível (GB): 1.32
Memória RAM usada (GB): 2.77
```

## 🏗️ Arquitetura do Projeto

```
system-manager/
├── main.go                    # Ponto de entrada da aplicação
├── go.mod                     # Definições do módulo Go
├── internal/
│   ├── hardware/              # Módulo de hardware
│   │   ├── bios.go           # Interface BIOS
│   │   ├── bios_*.go         # Implementações específicas por OS
│   │   ├── cpu.go            # Interface CPU
│   │   ├── cpu_*.go          # Implementações específicas por OS
│   │   ├── memory.go         # Interface Memória
│   │   ├── memory_*.go       # Implementações específicas por OS
│   │   ├── motherboard.go    # Interface Placa-mãe
│   │   ├── motherboard_*.go  # Implementações específicas por OS
│   │   └── systemInfo.go     # Agregador de informa��ões
│   └── system/               # Módulo de sistema
│       ├── system_operational.go     # Informações do SO
│       └── systemVersion_*.go        # Versões específicas por OS
└── README.md
```

## 🌐 Compatibilidade

| Sistema Operacional | Status | Versões Testadas |
|-------------------|--------|------------------|
| **macOS** | ✅ Suportado | 10.15+ (Intel & Apple Silicon) |
| **Linux** | ✅ Suportado | Ubuntu 18.04+, CentOS 7+, Fedora 30+ |
| **Windows** | ✅ Suportado | Windows 10, Windows 11 |

## 🔒 Permissões Necessárias

### Linux
```bash
# Para acesso completo às informações de hardware
sudo ./system-manager
```

### macOS
```bash
# Algumas informações podem requerer permissões administrativas
sudo ./system-manager
```

### Windows
```cmd
# Execute como Administrador para acesso completo
system-manager.exe
```

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor, siga estas diretrizes:

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Desenvolvimento

```bash
# Executar testes
go test ./...

# Verificar formatação
go fmt ./...

# Verificar linting
golangci-lint run
```

## 📝 Roadmap

### Versão 1.0 (Atual)
- [x] Informações básicas de sistema
- [x] Informações de CPU
- [x] Informações de memória
- [x] Informações de placa-mãe e BIOS

### Versão 1.1 (Próxima)
- [ ] Informações de rede
- [ ] Informações de discos e armazenamento

### Versão 2.0 (Futuro)
- [ ] Gerenciamento de bateria
- [ ] Interface web
- [ ] API REST
- [ ] Monitoramento em tempo real

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👨‍💻 Autor

**Miguel Lago**
- GitHub: [@miguellago](https://github.com/miguellago)

## 🙏 Agradecimentos

- Comunidade Go pela excelente documentação
- Contribuidores open-source que inspiraram este projeto

---

<div align="center">
  <p>⭐ Se este projeto foi útil para você, considere dar uma estrela!</p>
</div>
