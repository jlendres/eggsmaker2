package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFormatBytes(t *testing.T) {
	cases := map[int64]string{0: "0 B", 1023: "1023 B", 1024: "1.0 KiB", 1048576: "1.0 MiB"}
	for input, want := range cases {
		if got := formatBytes(input); got != want {
			t.Fatalf("formatBytes(%d)=%q, want %q", input, got, want)
		}
	}
}

func TestLatestISOSelectsNewest(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "old.iso")
	newer := filepath.Join(dir, "new.iso")
	if err := os.WriteFile(old, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newer, []byte("new"), 0600); err != nil {
		t.Fatal(err)
	}
	oldTime := time.Now().Add(-time.Hour)
	newTime := time.Now()
	if err := os.Chtimes(old, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(newer, newTime, newTime); err != nil {
		t.Fatal(err)
	}
	got, err := latestISO(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != newer {
		t.Fatalf("latestISO()=%q, want %q", got, newer)
	}
}

func TestLatestISORequiresISO(t *testing.T) {
	if _, err := latestISO(t.TempDir()); err == nil {
		t.Fatal("latestISO should fail when no ISO exists")
	}
}

func TestPackageSelectionPersistence(t *testing.T) {
	listPath := filepath.Join(t.TempDir(), "lista-nativo.txt")
	selected := []string{"bash", "linux-firmware"}
	if err := savePackageSelection(listPath, selected); err != nil {
		t.Fatal(err)
	}
	got := loadSavedSelection(listPath, []string{"bash", "linux-firmware", "vim"})
	if len(got) != 2 || got[0] != "bash" || got[1] != "linux-firmware" {
		t.Fatalf("loadSavedSelection()=%v, want %v", got, selected)
	}
	filtered := loadSavedSelection(listPath, []string{"bash", "vim"})
	if len(filtered) != 1 || filtered[0] != "bash" {
		t.Fatalf("loadSavedSelection() should ignore packages absent from current list: %v", filtered)
	}
}
