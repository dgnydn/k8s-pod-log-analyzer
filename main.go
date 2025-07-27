package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	namespace := ""
	since := "5m"
	language := LangEnglish // Default to English

	// Check for command line arguments
	if len(os.Args) > 1 {
		for i, arg := range os.Args[1:] {
			switch arg {
			case "-h", "--help":
				fmt.Println("Kubernetes Pod Log Analyzer")
				fmt.Println("Usage:")
				fmt.Println("  k8s-pod-log-analyzer [options]")
				fmt.Println("")
				fmt.Println("Options:")
				fmt.Println("  -n, --namespace <namespace>  Target namespace")
				fmt.Println("  -s, --since <duration>       Log duration (default: 5m)")
				fmt.Println("  --lang, --language <lang>    Language (en/tr, default: en)")
				fmt.Println("  -h, --help                   Show this help")
				fmt.Println("")
				fmt.Println("Examples:")
				fmt.Println("  k8s-pod-log-analyzer --lang tr")
				fmt.Println("  k8s-pod-log-analyzer -n kube-system --lang en")
				fmt.Println("  k8s-pod-log-analyzer -n default -s 10m --lang tr")
				os.Exit(0)
			case "-n", "--namespace":
				if i+2 < len(os.Args) {
					namespace = os.Args[i+2]
				}
			case "-s", "--since":
				if i+2 < len(os.Args) {
					since = os.Args[i+2]
				}
			case "--lang", "--language":
				if i+2 < len(os.Args) {
					langStr := os.Args[i+2]
					if langStr == "tr" || langStr == "turkish" {
						language = LangTurkish
					} else {
						language = LangEnglish
					}
				}
			}
		}
	}

	// If no namespace specified, start with namespace selection
	currentView := "pods"
	if namespace == "" {
		currentView = "namespaces"
	}

	// Get localization for selected language
	localization := GetLocalization(language)

	m := Model{
		namespace:    namespace,
		since:        since,
		logs:         make(map[string]LogAnalysis),
		currentView:  currentView,
		loading:      true,
		autoRefresh:  true,
		logOffset:    0,
		language:     language,
		localization: localization,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Hata: %v", err)
		os.Exit(1)
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	if m.currentView == "namespaces" {
		cmds = append(cmds, LoadNamespaces())
	} else {
		cmds = append(cmds, LoadPods(m.namespace))
	}

	// Start auto-refresh ticker
	if m.autoRefresh {
		cmds = append(cmds, Tick())
	}

	return tea.Batch(cmds...)
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case TickMsg:
		m.blinkState = !m.blinkState

		// Auto-refresh every 5 seconds
		if time.Since(time.Time(msg)).Seconds() >= 5 {
			var cmd tea.Cmd
			if m.currentView == "namespaces" {
				cmd = LoadNamespaces()
			} else if m.currentView == "pods" {
				cmd = LoadPods(m.namespace)
			}
			return m, tea.Batch(Tick(), cmd)
		}
		return m, Tick()

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case LoadNamespacesMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.namespaces = msg.namespaces
			m.err = nil
		}

	case LoadPodsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.pods = msg.pods
			m.err = nil
		}

	case LoadLogsMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.logs[msg.pod] = msg.analysis
			m.currentView = "analysis"
			m.err = nil
		}
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	if m.loading {
		return m.renderLoadingView()
	}

	if m.err != nil {
		return BorderStyle.Render(ErrorStyle.Render("❌ Error: ") + m.err.Error())
	}

	switch m.currentView {
	case "namespaces":
		return m.RenderNamespacesView()
	case "pods":
		return m.RenderPodsView()
	case "analysis":
		return m.RenderAnalysisView()
	default:
		return m.RenderNamespacesView()
	}
}

func (m Model) renderLoadingView() string {
	loadingText := m.localization.Loading
	if m.currentView == "namespaces" {
		loadingText = fmt.Sprintf("🔍 %s...", m.localization.Loading)
	} else if m.currentView == "pods" {
		loadingText = fmt.Sprintf("🔍 %s %s...", m.namespace, m.localization.Loading)
	}
	return BorderStyle.Render(loadingText)
}
