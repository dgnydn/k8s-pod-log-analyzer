package main

import (
	"time"
)

// PodInfo holds pod information including status
type PodInfo struct {
	Name       string
	Status     string
	Ready      string
	Age        string
	Restarts   string
	StatusIcon string
}

// LogAnalysis holds the analysis results for a pod
type LogAnalysis struct {
	TotalLines   int
	ErrorCount   int
	WarningCount int
	InfoCount    int
	Errors       []string
	Warnings     []string
	Info         []string
	RawLogs      string
	AnalyzedAt   time.Time
}

// Model represents the application state
type Model struct {
	namespace    string
	since        string
	namespaces   []string
	pods         []PodInfo
	selectedPod  int
	selectedNS   int
	logs         map[string]LogAnalysis
	currentView  string // "namespaces", "pods", "analysis"
	loading      bool
	err          error
	width        int
	height       int
	autoRefresh  bool
	blinkState   bool // For blinking error indicator
	pageOffset   int  // For pagination
	logOffset    int  // For log scrolling
	language     Language
	localization Localization
}

// Messages
type LoadNamespacesMsg struct {
	namespaces []string
	err        error
}

type LoadPodsMsg struct {
	pods []PodInfo
	err  error
}

type LoadLogsMsg struct {
	pod      string
	analysis LogAnalysis
	err      error
}

type TickMsg time.Time
