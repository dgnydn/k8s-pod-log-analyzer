package main

import (
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// LoadNamespaces command to fetch Kubernetes namespaces
func LoadNamespaces() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("kubectl", "get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
		output, err := cmd.Output()
		if err != nil {
			return LoadNamespacesMsg{err: err}
		}

		namespaceStr := strings.TrimSpace(string(output))
		var namespaces []string
		if namespaceStr != "" {
			namespaces = strings.Fields(namespaceStr)
		}

		return LoadNamespacesMsg{namespaces: namespaces}
	}
}

// LoadPods command to fetch pods in a namespace
func LoadPods(namespace string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "custom-columns=NAME:.metadata.name,STATUS:.status.phase,READY:.status.conditions[?(@.type=='Ready')].status,RESTARTS:.status.containerStatuses[0].restartCount,AGE:.metadata.creationTimestamp", "--no-headers")
		output, err := cmd.Output()
		if err != nil {
			return LoadPodsMsg{err: err}
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		var pods []PodInfo

		for _, line := range lines {
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				status := fields[1]
				ready := "False"
				if len(fields) > 2 && fields[2] != "<none>" {
					ready = fields[2]
				}
				restarts := "0"
				if len(fields) > 3 && fields[3] != "<none>" {
					restarts = fields[3]
				}
				age := "Unknown"
				if len(fields) > 4 {
					age = CalculateAge(fields[4])
				}

				icon := GetStatusIcon(status, ready)

				pods = append(pods, PodInfo{
					Name:       fields[0],
					Status:     status,
					Ready:      ready,
					Restarts:   restarts,
					Age:        age,
					StatusIcon: icon,
				})
			}
		}

		return LoadPodsMsg{pods: pods}
	}
}

// LoadLogs command to fetch and analyze pod logs
func LoadLogs(namespace, pod, since string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("kubectl", "logs", "-n", namespace, pod, "--since="+since)
		output, err := cmd.Output()
		if err != nil {
			return LoadLogsMsg{pod: pod, err: err}
		}

		analysis := AnalyzeLogs(string(output))
		return LoadLogsMsg{pod: pod, analysis: analysis}
	}
}

// Tick command for periodic updates
func Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
