package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getMaxVisibleItems calculates how many items can be displayed on screen
func (m Model) getMaxVisibleItems() int {
	// Reserve space for title, instructions, etc.
	// Approximate 8 lines for UI elements (less padding now)
	availableHeight := m.height - 8
	if availableHeight < 10 {
		return 10 // Minimum visible items
	}
	if availableHeight > 30 {
		return 30 // Maximum for better performance
	}
	return availableHeight
}

// handleKeyMsg processes keyboard input
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		if m.currentView == "analysis" {
			m.currentView = "pods"
		} else if m.currentView == "pods" && m.namespace != "" {
			m.currentView = "namespaces"
			m.namespace = ""
		}
	}

	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.currentView == "namespaces" && m.selectedNS > 0 {
			m.selectedNS--
			// Check if we need to scroll up
			if m.selectedNS < m.pageOffset {
				m.pageOffset = m.selectedNS
			}
		} else if m.currentView == "analysis" {
			// Scroll up in log analysis
			if m.logOffset < 100 { // Limit scroll to prevent going too far
				m.logOffset += 5
			}
		} else if m.currentView == "pods" && len(m.pods) > 0 {
			// Calculate pods per row dynamically
			availableWidth := m.width - 8
			minPodWidth := 45
			spacing := 2
			podsPerRow := 1
			for testPodsPerRow := 1; testPodsPerRow <= 6; testPodsPerRow++ {
				requiredWidth := testPodsPerRow*minPodWidth + (testPodsPerRow-1)*spacing
				if requiredWidth <= availableWidth {
					podsPerRow = testPodsPerRow
				} else {
					break
				}
			}

			// Move up one row
			newIndex := m.selectedPod - podsPerRow
			if newIndex >= 0 {
				m.selectedPod = newIndex
			}
		}
	case "down", "j":
		if m.currentView == "namespaces" && m.selectedNS < len(m.namespaces)-1 {
			m.selectedNS++
			// Check if we need to scroll down
			maxVisible := m.getMaxVisibleItems()
			if m.selectedNS >= m.pageOffset+maxVisible {
				m.pageOffset = m.selectedNS - maxVisible + 1
			}
		} else if m.currentView == "analysis" {
			// Scroll down in log analysis
			if m.logOffset > 0 {
				m.logOffset -= 5
				if m.logOffset < 0 {
					m.logOffset = 0
				}
			}
		} else if m.currentView == "pods" && len(m.pods) > 0 {
			// Calculate pods per row dynamically
			availableWidth := m.width - 8
			minPodWidth := 45
			spacing := 2
			podsPerRow := 1
			for testPodsPerRow := 1; testPodsPerRow <= 6; testPodsPerRow++ {
				requiredWidth := testPodsPerRow*minPodWidth + (testPodsPerRow-1)*spacing
				if requiredWidth <= availableWidth {
					podsPerRow = testPodsPerRow
				} else {
					break
				}
			}

			// Move down one row
			newIndex := m.selectedPod + podsPerRow
			if newIndex < len(m.pods) {
				m.selectedPod = newIndex
			}
		}
	case "left", "h":
		if m.currentView == "pods" && m.selectedPod > 0 {
			// Calculate pods per row dynamically
			availableWidth := m.width - 8
			minPodWidth := 45
			spacing := 2
			podsPerRow := 1
			for testPodsPerRow := 1; testPodsPerRow <= 6; testPodsPerRow++ {
				requiredWidth := testPodsPerRow*minPodWidth + (testPodsPerRow-1)*spacing
				if requiredWidth <= availableWidth {
					podsPerRow = testPodsPerRow
				} else {
					break
				}
			}

			// Check if we're not at the beginning of a row
			currentCol := m.selectedPod % podsPerRow

			if currentCol > 0 {
				m.selectedPod--
			}
		}
	case "right", "l":
		if m.currentView == "pods" && m.selectedPod < len(m.pods)-1 {
			// Calculate pods per row dynamically
			availableWidth := m.width - 8
			minPodWidth := 45
			spacing := 2
			podsPerRow := 1
			for testPodsPerRow := 1; testPodsPerRow <= 6; testPodsPerRow++ {
				requiredWidth := testPodsPerRow*minPodWidth + (testPodsPerRow-1)*spacing
				if requiredWidth <= availableWidth {
					podsPerRow = testPodsPerRow
				} else {
					break
				}
			}

			// Check if we're not at the end of a row
			currentRow := m.selectedPod / podsPerRow
			currentCol := m.selectedPod % podsPerRow
			maxCol := min(podsPerRow-1, len(m.pods)-currentRow*podsPerRow-1)

			if currentCol < maxCol {
				m.selectedPod++
			}
		}
	case "enter":
		if m.currentView == "namespaces" && len(m.namespaces) > 0 {
			m.namespace = m.namespaces[m.selectedNS]
			m.currentView = "pods"
			m.loading = true
			return m, LoadPods(m.namespace)
		} else if m.currentView == "pods" && len(m.pods) > 0 {
			selectedPod := m.pods[m.selectedPod].Name
			m.logOffset = 0 // Reset scroll position when entering analysis
			return m, LoadLogs(m.namespace, selectedPod, m.since)
		}
	case "backspace":
		if m.currentView == "analysis" {
			m.currentView = "pods"
		} else if m.currentView == "pods" && m.namespace != "" {
			m.currentView = "namespaces"
			m.namespace = ""
		}
	case "r":
		// Refresh
		m.loading = true
		if m.currentView == "namespaces" {
			return m, LoadNamespaces()
		} else if m.currentView == "pods" {
			return m, LoadPods(m.namespace)
		} else if m.currentView == "analysis" && len(m.pods) > 0 {
			selectedPod := m.pods[m.selectedPod].Name
			return m, LoadLogs(m.namespace, selectedPod, m.since)
		}
	case "t":
		// Toggle auto-refresh
		m.autoRefresh = !m.autoRefresh
		if m.autoRefresh {
			return m, Tick()
		}
	}

	return m, nil
}

// CalculateAge calculates pod age from timestamp
func CalculateAge(timestamp string) string {
	if timestamp == "" || timestamp == "<none>" {
		return "Unknown"
	}

	// Parse Kubernetes timestamp format
	// Example: 2025-07-27T14:27:59Z
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		// Try alternative format without Z
		layout = "2006-01-02T15:04:05"
		t, err = time.Parse(layout, timestamp)
		if err != nil {
			return "Unknown"
		}
	}

	duration := time.Since(t)

	if duration.Hours() < 1 {
		return fmt.Sprintf("%.0fm", duration.Minutes())
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%.0fh", duration.Hours())
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}
}

// GetAutoRefreshIndicator returns the appropriate auto-refresh indicator
func GetAutoRefreshIndicator(autoRefresh, blinkState bool) string {
	if autoRefresh {
		if blinkState {
			return "🔄"
		} else {
			return "⏸️"
		}
	}
	return "⏸️"
}

// GetStatusIcon returns appropriate icon for pod status
func GetStatusIcon(status, ready string) string {
	switch status {
	case "Running":
		if ready == "True" {
			return "✅"
		}
		return "🟡"
	case "Pending":
		return "⏳"
	case "Failed", "Error":
		return "❌"
	case "Succeeded":
		return "✅"
	case "Terminating":
		return "🟠"
	case "CrashLoopBackOff":
		return "💥"
	case "ImagePullBackOff":
		return "📥"
	case "ContainerCreating":
		return "🔧"
	default:
		return "❔"
	}
}

// GetStatusStyle returns appropriate style for pod status
func GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "Running":
		return RunningStyle
	case "Pending", "ContainerCreating":
		return PendingStyle
	case "Failed", "Error", "CrashLoopBackOff", "ImagePullBackOff":
		return FailedStyle
	case "Terminating":
		return TerminatingStyle
	default:
		return UnknownStyle
	}
}

// formatSummaryLine formats a summary line with style
func formatSummaryLine(label, value string, count int) string {
	switch label {
	case "Hatalar":
		return fmt.Sprintf("  %s: %s\n", label, ErrorStyle.Render(strconv.Itoa(count)))
	case "Uyarılar":
		return fmt.Sprintf("  %s: %s\n", label, WarningStyle.Render(strconv.Itoa(count)))
	case "Bilgiler":
		return fmt.Sprintf("  %s: %s\n", label, SuccessStyle.Render(strconv.Itoa(count)))
	default:
		return fmt.Sprintf("  %s: %s\n", label, InfoStyle.Render(strconv.Itoa(count)))
	}
}

// renderPodBox creates a styled box for a single pod with dynamic width
func (m Model) renderPodBox(pod PodInfo, isSelected bool, width int) string {
	// Determine box style based on selection
	boxStyle := PodBoxStyle.Width(width)
	nameStyle := NormalStyle
	if isSelected {
		boxStyle = SelectedPodBoxStyle.Width(width)
		nameStyle = SelectedStyle
	}

	// Calculate name truncation based on box width
	maxNameLen := width - 8 // Account for padding and borders
	displayName := pod.Name
	if len(displayName) > maxNameLen {
		displayName = displayName[:maxNameLen-3] + "..."
	}

	// Format status with appropriate color and blinking
	statusText := fmt.Sprintf("%s %s", pod.StatusIcon, pod.Status)
	if pod.Status == "Failed" || pod.Status == "Error" || pod.Status == "CrashLoopBackOff" {
		if m.blinkState {
			statusText = ErrorStyle.Render(statusText)
		} else {
			statusText = FailedStyle.Render(statusText)
		}
	} else {
		statusStyle := GetStatusStyle(pod.Status)
		statusText = statusStyle.Render(statusText)
	}

	// Format ready status
	readyText := pod.Ready
	if pod.Ready == "True" {
		readyText = SuccessStyle.Render("✓ Ready")
	} else {
		readyText = WarningStyle.Render("✗ Not Ready")
	}

	// Format restart count with color
	restartText := pod.Restarts + " restarts"
	if restarts, err := strconv.Atoi(pod.Restarts); err == nil {
		if restarts > 10 {
			restartText = ErrorStyle.Render(fmt.Sprintf("%d restarts", restarts))
		} else if restarts > 3 {
			restartText = WarningStyle.Render(fmt.Sprintf("%d restarts", restarts))
		} else {
			restartText = SuccessStyle.Render(fmt.Sprintf("%d restarts", restarts))
		}
	}

	// Create box content
	content := fmt.Sprintf("%s\n%s\n%s\n%s\nAge: %s",
		nameStyle.Render(displayName),
		statusText,
		readyText,
		restartText,
		InfoStyle.Render(pod.Age),
	)

	return boxStyle.Render(content)
}

// joinHorizontal joins multiple strings horizontally with spacing
func (m Model) joinHorizontal(boxes []string, boxWidth int) string {
	if len(boxes) == 0 {
		return ""
	}
	if len(boxes) == 1 {
		return boxes[0]
	}

	// Split each box into lines
	boxLines := make([][]string, len(boxes))
	maxLines := 0

	for i, box := range boxes {
		lines := strings.Split(box, "\n")
		boxLines[i] = lines
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}

	// Join lines horizontally
	var result []string
	for lineNum := 0; lineNum < maxLines; lineNum++ {
		var lineParts []string
		for boxNum := 0; boxNum < len(boxLines); boxNum++ {
			if lineNum < len(boxLines[boxNum]) {
				lineParts = append(lineParts, boxLines[boxNum][lineNum])
			} else {
				// Add empty space for shorter boxes using actual box width
				lineParts = append(lineParts, strings.Repeat(" ", boxWidth))
			}
		}
		result = append(result, strings.Join(lineParts, "  ")) // 2 spaces between boxes
	}

	return strings.Join(result, "\n")
}

// truncateLogLine truncates a log line to fit within specified width
func (m Model) truncateLogLine(line string, maxWidth int) string {
	if len(line) <= maxWidth {
		return line
	}
	if maxWidth < 10 {
		return line[:maxWidth]
	}
	return line[:maxWidth-3] + "..."
}
