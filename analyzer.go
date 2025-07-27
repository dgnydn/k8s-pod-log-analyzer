package main

import (
	"bufio"
	"regexp"
	"strings"
	"time"
)

// AnalyzeLogs analyzes pod logs and extracts errors, warnings, and info
func AnalyzeLogs(logs string) LogAnalysis {
	analysis := LogAnalysis{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
		Info:     make([]string, 0),
		RawLogs:  logs, // EN ÖNEMLİSİ: Raw log'ları sakla!
	}

	scanner := bufio.NewScanner(strings.NewReader(logs))

	// Enhanced regex patterns for better detection
	errorPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\b(error|err|exception|fatal|panic|crash|failed|failure)\b`),
		regexp.MustCompile(`(?i)\b(stack\s+trace|stacktrace)\b`),
		regexp.MustCompile(`(?i)\b(connection\s+(refused|failed|timeout))\b`),
		regexp.MustCompile(`(?i)\b(out\s+of\s+memory|oom)\b`),
		regexp.MustCompile(`(?i)\b(permission\s+denied|access\s+denied)\b`),
	}

	warningPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\b(warn|warning|deprecated|timeout|retry|retrying)\b`),
		regexp.MustCompile(`(?i)\b(slow\s+query|performance)\b`),
		regexp.MustCompile(`(?i)\b(connection\s+lost|reconnecting)\b`),
	}

	infoPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\b(info|starting|started|listening|ready|success|successful|completed)\b`),
		regexp.MustCompile(`(?i)\b(connected|initialized|loaded)\b`),
	}

	for scanner.Scan() {
		line := scanner.Text()
		analysis.TotalLines++

		isError := false
		isWarning := false

		// Check for errors
		for _, pattern := range errorPatterns {
			if pattern.MatchString(line) {
				analysis.ErrorCount++
				analysis.Errors = append(analysis.Errors, line)
				isError = true
				break
			}
		}

		// Check for warnings (only if not already an error)
		if !isError {
			for _, pattern := range warningPatterns {
				if pattern.MatchString(line) {
					analysis.WarningCount++
					analysis.Warnings = append(analysis.Warnings, line)
					isWarning = true
					break
				}
			}
		}

		// Check for info (only if not already an error or warning)
		if !isError && !isWarning {
			for _, pattern := range infoPatterns {
				if pattern.MatchString(line) {
					analysis.InfoCount++
					analysis.Info = append(analysis.Info, line)
					break
				}
			}
		}
	}

	// Analiz zamanını ekle
	analysis.AnalyzedAt = time.Now()

	return analysis
}
