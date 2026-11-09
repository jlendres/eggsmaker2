package main

import "testing"

func TestFindPartitionManagerPrefersGParted(t *testing.T) {
	available := map[string]bool{"gparted": true, "partitionmanager": true, "gnome-disks": true}
	got := findPartitionManager("KDE", func(name string) bool { return available[name] })
	if got != "gparted" {
		t.Fatalf("findPartitionManager() = %q; want gparted", got)
	}
}

func TestFindPartitionManagerUsesDesktopFallback(t *testing.T) {
	tests := []struct {
		desktop   string
		available map[string]bool
		want      string
	}{
		{"KDE", map[string]bool{"partitionmanager": true, "gnome-disks": true}, "partitionmanager"},
		{"GNOME:GNOME", map[string]bool{"partitionmanager": true, "gnome-disks": true}, "gnome-disks"},
		{"XFCE", map[string]bool{"gnome-disks": true}, "gnome-disks"},
		{"KDE", map[string]bool{}, ""},
	}
	for _, test := range tests {
		got := findPartitionManager(test.desktop, func(name string) bool {
			return test.available[name]
		})
		if got != test.want {
			t.Errorf("findPartitionManager(%q) = %q; want %q", test.desktop, got, test.want)
		}
	}
}
