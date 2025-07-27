package main

import (
	"os/exec"
	"strings"
)

func getNamespaces() ([]string, error) {
	cmd := exec.Command("kubectl", "get", "namespaces", "-o", "custom-columns=:metadata.name")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

func getPods(namespace string) ([]PodInfo, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "custom-columns=NAME:.metadata.name,STATUS:.status.phase,READY:.status.conditions[?(@.type==\"Ready\")].status,AGE:.metadata.creationTimestamp")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var pods []PodInfo

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			pod := PodInfo{
				Name:   parts[0],
				Status: parts[1],
				Ready:  "Unknown",
				Age:    "Unknown",
			}
			if len(parts) >= 4 {
				pod.Ready = parts[2]
				if len(parts) >= 5 {
					pod.Age = parts[3]
				}
			}
			pods = append(pods, pod)
		}
	}

	return pods, nil
}

func getLogs(namespace, pod, since string) (string, error) {
	cmd := exec.Command("kubectl", "logs", pod, "-n", namespace, "--since="+since)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
