package main

import (
	"fmt"
	"strconv"
	"strings"
)

// RenderNamespacesView renders the namespace selection view
func (m Model) RenderNamespacesView() string {
	title := TitleStyle.Render(m.localization.NamespaceSelectionTitle)

	var content strings.Builder
	content.WriteString(title + "\n\n")

	if len(m.namespaces) == 0 {
		content.WriteString(m.localization.NamespaceNotFound + "\n")
	} else {
		maxVisible := m.getMaxVisibleItems()
		start := m.pageOffset
		end := start + maxVisible
		if end > len(m.namespaces) {
			end = len(m.namespaces)
		}

		// Show pagination info if needed
		if len(m.namespaces) > maxVisible {
			content.WriteString(fmt.Sprintf("%d-%d / %d %s\n", start+1, end, len(m.namespaces), m.localization.Namespaces))
			if start > 0 {
				content.WriteString(m.localization.ScrollUp + "\n")
			}
			if end < len(m.namespaces) {
				content.WriteString(m.localization.ScrollDown + "\n")
			}
			content.WriteString("\n")
		}

		for i := start; i < end; i++ {
			ns := m.namespaces[i]
			prefix := "  "
			style := NormalStyle
			if i == m.selectedNS {
				prefix = "> "
				style = SelectedStyle
			}

			content.WriteString(fmt.Sprintf("%s%-30s\n", prefix, style.Render(ns)))
		}
	}

	content.WriteString(fmt.Sprintf("\n%s: %t\n", m.localization.AutoRefreshStatus, m.autoRefresh))
	content.WriteString("\n" + m.localization.Controls + ":\n")
	content.WriteString("  " + m.localization.Movement + "\n")
	content.WriteString("  " + m.localization.Select + "\n")
	content.WriteString("  " + m.localization.Refresh + "\n")
	content.WriteString("  " + m.localization.AutoRefresh + "\n")
	content.WriteString("  " + m.localization.Exit)

	return BorderStyle.Render(content.String())
}

// RenderPodsView renders the pod list view
func (m Model) RenderPodsView() string {
	title := TitleStyle.Render(fmt.Sprintf("%s: %s", m.localization.NamespaceTitle, m.namespace))

	var content strings.Builder
	content.WriteString(title + "\n\n")

	if len(m.pods) == 0 {
		content.WriteString(m.localization.PodNotFound + "\n")
	} else {
		// Calculate optimal pod dimensions based on terminal width
		podsPerRow := m.getPodsPerRow()

		// Reserve space for borders and margins
		availableWidth := m.width - 8 // BorderStyle padding and margins
		maxPodWidth := 60             // Maximum pod width to prevent overly wide boxes

		// Calculate actual pod width using remaining space
		totalSpacing := (podsPerRow - 1) * 2 // 2 = spacing between pods
		availableForPods := availableWidth - totalSpacing
		podWidth := min(maxPodWidth, availableForPods/podsPerRow)

		// Calculate how many rows we can display based on available height
		// Each pod row takes about 8 lines (pod box height + spacing)
		podRowHeight := 8
		availableHeight := m.height - 15 // Reserve space for title, controls, etc.
		maxVisibleRows := max(1, availableHeight/podRowHeight)

		// Calculate total rows and pagination
		totalRows := (len(m.pods) + podsPerRow - 1) / podsPerRow

		// Calculate scroll offset in terms of rows
		selectedRow := m.selectedPod / podsPerRow
		scrollOffset := 0

		if totalRows > maxVisibleRows {
			// Calculate scroll offset to keep selected pod visible
			if selectedRow >= maxVisibleRows {
				scrollOffset = selectedRow - maxVisibleRows + 1
			}
			if scrollOffset > totalRows-maxVisibleRows {
				scrollOffset = totalRows - maxVisibleRows
			}
		}

		// Show pagination info if needed
		if totalRows > maxVisibleRows {
			startRow := scrollOffset + 1
			endRow := min(totalRows, scrollOffset+maxVisibleRows)
			content.WriteString(fmt.Sprintf("Rows %d-%d / %d total (%d pods)\n",
				startRow, endRow, totalRows, len(m.pods)))
			if scrollOffset > 0 {
				content.WriteString(m.localization.ScrollUp + "\n")
			}
			if scrollOffset+maxVisibleRows < totalRows {
				content.WriteString(m.localization.ScrollDown + "\n")
			}
			content.WriteString("\n")
		}

		content.WriteString(m.localization.Pods + ":\n\n")

		// Render visible rows only
		startRowIndex := scrollOffset
		endRowIndex := min(totalRows, scrollOffset+maxVisibleRows)

		for rowIndex := startRowIndex; rowIndex < endRowIndex; rowIndex++ {
			startPodIndex := rowIndex * podsPerRow
			endPodIndex := min(len(m.pods), startPodIndex+podsPerRow)

			// Create a row of pods
			var rowBoxes []string

			for podIndex := startPodIndex; podIndex < endPodIndex; podIndex++ {
				pod := m.pods[podIndex]

				// Determine if this pod is selected
				isSelected := podIndex == m.selectedPod

				// Create individual pod box with calculated width
				podBox := m.renderPodBox(pod, isSelected, podWidth)
				rowBoxes = append(rowBoxes, podBox)
			}

			// Join boxes horizontally
			if len(rowBoxes) > 0 {
				content.WriteString(m.joinHorizontal(rowBoxes, podWidth) + "\n\n")
			}
		}
	}

	content.WriteString(fmt.Sprintf("\n%s: %t\n", m.localization.AutoRefreshStatus, m.autoRefresh))
	content.WriteString("\n" + m.localization.Controls + ":\n")
	content.WriteString("  " + m.localization.UpDown + ": Navigate rows\n")
	content.WriteString("  " + m.localization.LeftRight + ": Navigate columns\n")
	content.WriteString("  Page Up/Down, Ctrl+U/D: Fast scroll\n")
	content.WriteString("  Home/End, g/G: Go to first/last pod\n")
	content.WriteString("  " + m.localization.ViewLogs + "\n")
	content.WriteString("  Esc/Backspace: " + m.localization.NamespaceTitle + "\n")
	content.WriteString("  " + m.localization.Refresh + "\n")
	content.WriteString("  " + m.localization.AutoRefresh + "\n")
	content.WriteString("  " + m.localization.Exit)

	return BorderStyle.Render(content.String())
}

// RenderAnalysisView renders the log analysis view
func (m Model) RenderAnalysisView() string {
	selectedPod := m.pods[m.selectedPod].Name
	analysis, exists := m.logs[selectedPod]

	if !exists {
		return BorderStyle.Render(m.localization.LogNotFound + "\n\n" + m.localization.Loading)
	}

	title := TitleStyle.Render(fmt.Sprintf("%s: %s", m.localization.LogAnalysisTitle, selectedPod))

	var content strings.Builder
	content.WriteString(title + "\n\n")

	// Pod bilgileri (sadeleştirilmiş)
	selectedPodInfo := m.pods[m.selectedPod]
	content.WriteString(m.localization.PodDetails + ":\n")
	content.WriteString(fmt.Sprintf("  %s: %s\n", m.localization.Name, SelectedStyle.Render(selectedPodInfo.Name)))
	content.WriteString(fmt.Sprintf("  %s: %s\n", m.localization.Status, GetStatusStyle(selectedPodInfo.Status).Render(selectedPodInfo.Status)))
	content.WriteString(fmt.Sprintf("  %s: %s\n", m.localization.Ready, selectedPodInfo.Ready))
	content.WriteString(fmt.Sprintf("  %s: %s\n", m.localization.Restart, selectedPodInfo.Restarts))
	content.WriteString(fmt.Sprintf("  %s: %s\n", m.localization.Age, selectedPodInfo.Age))
	content.WriteString(fmt.Sprintf("  %s: %s\n", m.localization.Analysis, analysis.AnalyzedAt.Format("15:04:05")))
	content.WriteString("\n")

	// Log özeti (sadeleştirilmiş)
	content.WriteString(m.localization.LogSummary + ":\n")
	content.WriteString(fmt.Sprintf("  %s: %s\n", m.localization.TotalLines, InfoStyle.Render(strconv.Itoa(analysis.TotalLines))))
	if analysis.ErrorCount > 0 {
		content.WriteString(fmt.Sprintf("  %s: %s\n", m.localization.Errors, ErrorStyle.Render(strconv.Itoa(analysis.ErrorCount))))
	}
	if analysis.WarningCount > 0 {
		content.WriteString(fmt.Sprintf("  %s: %s\n", m.localization.Warnings, WarningStyle.Render(strconv.Itoa(analysis.WarningCount))))
	}
	content.WriteString("\n")

	// MAIN SECTION: RAW LOG LINES
	if analysis.RawLogs != "" {
		content.WriteString(m.localization.LogLines + ":\n")
		content.WriteString(strings.Repeat("-", min(80, m.width-10)) + "\n")

		lines := strings.Split(strings.TrimSpace(analysis.RawLogs), "\n")

		// Calculate visible lines based on terminal height
		maxVisibleLines := max(10, (m.height - 25)) // Reserve space for other UI elements
		totalLines := len(lines)

		// Apply scroll offset
		startIdx := max(0, totalLines-maxVisibleLines-m.logOffset)
		endIdx := min(totalLines, startIdx+maxVisibleLines)

		// Show pagination info for logs
		if totalLines > maxVisibleLines {
			content.WriteString(fmt.Sprintf("Showing lines %d-%d of %d total lines\n",
				startIdx+1, endIdx, totalLines))
			if m.logOffset > 0 {
				content.WriteString("↑ Scroll up for more logs\n")
			}
			if startIdx > 0 {
				content.WriteString("↓ Scroll down for more logs\n")
			}
			content.WriteString(strings.Repeat("-", min(80, m.width-10)) + "\n")
		}

		for i := startIdx; i < endIdx; i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			// Satır numarası ile birlikte göster
			lineNum := i + 1
			truncatedLine := m.truncateLogLine(line, m.width-15)

			// Log seviyesine göre renklendirme (icon olmadan)
			lowerLine := strings.ToLower(line)
			if strings.Contains(lowerLine, "error") || strings.Contains(lowerLine, "fail") ||
				strings.Contains(lowerLine, "exception") || strings.Contains(lowerLine, "panic") {
				content.WriteString(fmt.Sprintf("  %4d: %s\n", lineNum, ErrorStyle.Render(truncatedLine)))
			} else if strings.Contains(lowerLine, "warn") || strings.Contains(lowerLine, "warning") {
				content.WriteString(fmt.Sprintf("  %4d: %s\n", lineNum, WarningStyle.Render(truncatedLine)))
			} else if strings.Contains(lowerLine, "info") {
				content.WriteString(fmt.Sprintf("  %4d: %s\n", lineNum, InfoStyle.Render(truncatedLine)))
			} else {
				content.WriteString(fmt.Sprintf("  %4d: %s\n", lineNum, NormalStyle.Render(truncatedLine)))
			}
		}

		content.WriteString(strings.Repeat("-", min(80, m.width-10)) + "\n")
		content.WriteString(fmt.Sprintf("%s %d %s %d %s\n\n",
			m.localization.TotalFrom, totalLines, m.localization.LastLines, endIdx-startIdx, m.localization.ShowingLines))
	} else {
		content.WriteString(m.localization.LogEmpty + "\n\n")
	}

	// Basit durum özeti
	if analysis.ErrorCount > 0 {
		content.WriteString(ErrorStyle.Render(m.localization.StatusError) + "\n\n")
	} else if analysis.WarningCount > 0 {
		content.WriteString(WarningStyle.Render(m.localization.StatusWarning) + "\n\n")
	} else {
		content.WriteString(SuccessStyle.Render(m.localization.StatusNormal) + "\n\n")
	}

	content.WriteString(m.localization.Controls + ":\n")
	content.WriteString("  " + m.localization.UpDown + ": " + m.localization.ScrollUp + "/" + m.localization.ScrollDown + "\n")
	content.WriteString("  " + m.localization.GoBack + "\n")
	content.WriteString("  " + m.localization.RefreshLogs + "\n")
	content.WriteString("  " + m.localization.Exit)

	return BorderStyle.Render(content.String())
}
