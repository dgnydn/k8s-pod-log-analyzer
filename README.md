# 🚀 Kubernetes Pod Log Analyzer

A modern, interactive Terminal User Interface (TUI) for analyzing Kubernetes pod logs with intelligent pattern recognition and multilingual support.

![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey.svg)

## ✨ Features

- 🎨 **Beautiful TUI**: Modern interface built with Bubble Tea framework
- 🔍 **Smart Log Analysis**: Advanced regex patterns to categorize errors, warnings, and info messages
- 🌍 **Multilingual Support**: English and Turkish language support
- 📊 **Detailed Statistics**: Total lines, errors, warnings count with analysis timestamps
- ⚡ **Real-time Updates**: Live pod listing with auto-refresh capabilities
- 📱 **Responsive Design**: Adapts to terminal size with intelligent layout
- 🎯 **Namespace Support**: Analyze pods from any namespace
- 🔄 **Auto-refresh**: Automatic updates with visual indicators
- 📖 **Raw Log Display**: View actual log lines with syntax highlighting

## 🎬 Demo

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Kubernetes Namespace Selection                                              │
│                                                                             │
│ ▶ default                                                                   │
│   kube-system                                                               │
│   monitoring                                                                │
│   ingress-nginx                                                             │
│                                                                             │
│ Auto-refresh: true                                                          │
│                                                                             │
│ Controls:                                                                   │
│   k/j or arrow keys: Move                                                   │
│   Enter: Select namespace                                                   │
│   r: Refresh                                                                │
│   t: Auto-refresh                                                           │
│   q: Exit                                                                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - [Install Go](https://golang.org/doc/install)
- **kubectl** - [Install kubectl](https://kubernetes.io/docs/tasks/tools/)
- **Active Kubernetes cluster connection**

### Installation

```bash
# Clone the repository
git clone https://github.com/dgnydn/k8s-pod-log-analyzer.git
cd k8s-pod-log-analyzer

# Build the application
go build -o k8s-log-analyzer .

# Run with default settings
./k8s-log-analyzer
```

### Usage Examples

```bash
# Basic usage (English, default namespace, last 5 minutes)
./k8s-log-analyzer

# Specify namespace and time range
./k8s-log-analyzer --namespace production --since 10m

# Turkish interface with custom settings
./k8s-log-analyzer --lang tr --namespace kube-system --since 1h

# Show help
./k8s-log-analyzer --help
```

## 🎮 Controls

### Namespace Selection

| Key            | Action                 |
| -------------- | ---------------------- |
| `↑/↓` or `k/j` | Navigate namespaces    |
| `Enter`        | Select namespace       |
| `r`            | Refresh namespace list |
| `t`            | Toggle auto-refresh    |
| `q`            | Exit application       |

### Pod Grid View

| Key                    | Action                        |
| ---------------------- | ----------------------------- |
| `↑/↓/←/→` or `k/j/h/l` | Navigate pods                 |
| `Enter`                | View pod logs                 |
| `Esc/Backspace`        | Return to namespace selection |
| `r`                    | Refresh pod list              |
| `t`                    | Toggle auto-refresh           |
| `q`                    | Exit application              |

### Log Analysis View

| Key             | Action              |
| --------------- | ------------------- |
| `↑/↓` or `k/j`  | Scroll through logs |
| `Esc/Backspace` | Return to pod grid  |
| `r`             | Refresh logs        |
| `q`             | Exit application    |

## 📊 Log Analysis Features

The application intelligently categorizes log entries:

### 🚨 Error Detection

Automatically identifies error patterns:

- `error`, `err`, `failed`, `failure`
- `fatal`, `exception`, `panic`
- `timeout`, `refused`, `denied`

### ⚠️ Warning Detection

Catches warning indicators:

- `warn`, `warning`
- `deprecated`, `deprecation`
- `retry`, `retrying`

### ℹ️ Information Logs

All other log entries are categorized as informational with proper highlighting.

## 🌍 Multilingual Support

The application supports multiple languages through the `--lang` parameter:

| Language          | Code | Example                        |
| ----------------- | ---- | ------------------------------ |
| English (default) | `en` | `./k8s-log-analyzer --lang en` |
| Turkish           | `tr` | `./k8s-log-analyzer --lang tr` |

## 🏗️ Architecture

The project follows a clean, modular architecture:

```
├── main.go          # Application entry point and CLI parsing
├── types.go         # Data structures and type definitions
├── localization.go  # Multilingual text management
├── views.go         # TUI rendering and layouts
├── commands.go      # Kubernetes API interactions
├── analyzer.go      # Log analysis and pattern matching
├── helpers.go       # Utility functions
└── styles.go        # Terminal styling and themes
```

## 🛠️ Development

### Tech Stack

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Modern TUI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **Go 1.21+** - Backend language
- **kubectl** - Kubernetes API access

### Building from Source

```bash
# Clone and enter directory
git clone https://github.com/dgnydn/k8s-pod-log-analyzer.git
cd k8s-pod-log-analyzer

# Install dependencies
go mod download

# Run in development mode
go run .

# Build for production
go build -ldflags="-s -w" -o k8s-log-analyzer .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines

- Follow Go coding standards and conventions
- Add tests for new features
- Update documentation for API changes
- Ensure multilingual support for new text
- Test on multiple terminal environments

## 🐛 Troubleshooting

### Common Issues

**kubectl not found**

```bash
# Ensure kubectl is installed and in PATH
kubectl version --client
```

**Permission denied**

```bash
# Check kubectl permissions
kubectl auth can-i get pods --all-namespaces
```

**No pods found**

```bash
# Verify namespace exists and contains pods
kubectl get pods -n <namespace>
```

## 📈 Roadmap

- [ ] **Export Options**: JSON, CSV, and HTML export formats
- [ ] **Advanced Filtering**: Custom regex patterns and filters
- [ ] **Log Streaming**: Real-time log tailing capability
- [ ] **Plugin System**: Extensible analysis plugins
- [ ] **Cluster Metrics**: Resource usage and health monitoring
- [ ] **Theme Support**: Customizable color schemes
- [ ] **Configuration Files**: Persistent settings and preferences

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - For the amazing TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - For beautiful terminal styling
- [Kubernetes Community](https://kubernetes.io/) - For the incredible platform

## 📞 Support

- 🐛 **Bug Reports**: [Create an issue](https://github.com/dgnydn/k8s-pod-log-analyzer/issues)
- 💡 **Feature Requests**: [Start a discussion](https://github.com/dgnydn/k8s-pod-log-analyzer/discussions)
- 📖 **Documentation**: [Wiki](https://github.com/dgnydn/k8s-pod-log-analyzer/wiki)

---

<div align="center">
  <b>⭐ Star this repository if you find it helpful! ⭐</b>
</div>
