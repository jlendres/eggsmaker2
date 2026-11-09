package main

import (
	"os"
	"os/exec"
	"strings"
)

func findPartitionManager(desktop string, available func(string) bool) string {
	if available("gparted") {
		return "gparted"
	}

	desktop = strings.ToLower(desktop)
	candidates := []string{"partitionmanager", "gnome-disks"}
	if strings.Contains(desktop, "gnome") || strings.Contains(desktop, "unity") {
		candidates = []string{"gnome-disks", "partitionmanager"}
	}
	for _, candidate := range candidates {
		if available(candidate) {
			return candidate
		}
	}
	return ""
}

func installedPartitionManager() string {
	return findPartitionManager(os.Getenv("XDG_CURRENT_DESKTOP"), func(name string) bool {
		_, err := exec.LookPath(name)
		return err == nil
	})
}
