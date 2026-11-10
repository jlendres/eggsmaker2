package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRunEachPackageContinuesAfterErrors(t *testing.T) {
	selected := []string{"first", "missing", "last"}
	var attempted []string
	var logged strings.Builder
	wantErr := errors.New("package unavailable")
	err := runEachPackage("instalar", selected, func(name string) error {
		attempted = append(attempted, name)
		if name == "missing" {
			return wantErr
		}
		return nil
	}, func(format string, args ...any) {
		_, _ = fmt.Fprintf(&logged, format, args...)
	})
	if !reflect.DeepEqual(attempted, selected) {
		t.Fatalf("attempted packages = %v; want %v", attempted, selected)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("runEachPackage() error = %v; want wrapped package error", err)
	}
	if !strings.Contains(logged.String(), "Finalizó instalar: 1 de 3 paquetes tuvieron errores") {
		t.Fatalf("failure summary was not logged: %s", logged.String())
	}
}

func TestPackageBatchContinuesAfterPackageError(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "calls.log")
	manager := filepath.Join(dir, "fake-pacman")
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$PACKAGE_BATCH_TEST_LOG\"\ncase \"$*\" in *fail-package*) exit 7 ;; esac\n"
	if err := os.WriteFile(manager, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PACKAGE_BATCH_TEST_LOG", logPath)

	cmd, err := packageBatchCommand(context.Background(), "pacman-install", manager, []string{"first-package", "fail-package", "last-package"})
	if err != nil {
		t.Fatal(err)
	}
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("package batch succeeded despite a package failure")
	}
	if !strings.Contains(string(output), "se continúa con el siguiente") {
		t.Fatalf("batch did not report continuation: %s", output)
	}
	calls, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"first-package", "fail-package", "last-package"} {
		if !strings.Contains(string(calls), name) {
			t.Errorf("package %q was not attempted; calls: %s", name, calls)
		}
	}
}

func TestFlatpakInstallFilterSkipsInstalledAndMissingApps(t *testing.T) {
	dir := t.TempDir()
	flatpak := filepath.Join(dir, "flatpak")
	script := `#!/bin/sh
case "$1" in
  remote-ls) printf '%s\n' org.example.installed org.example.available ;;
  list) printf '%s\n' org.example.installed ;;
  *) exit 2 ;;
esac
`
	if err := os.WriteFile(flatpak, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	state := &AppState{}
	got, err := state.filterSelectedPackages(context.Background(), "flatpak", true, []string{
		"org.example.installed",
		"org.example.available",
		"org.example.missing",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{"org.example.available"}) {
		t.Fatalf("filtered apps = %v; want [org.example.available]", got)
	}
}

func TestArchRepoPackages(t *testing.T) {
	output := []byte("core bash 5.2-1\nextra vim 9.1-1\nchaotic-aur app 1.0-1\ninvalid\n")
	if got, want := archRepoPackages(output, "chaotic-aur"), map[string]struct{}{"app": {}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("chaotic repo packages = %v; want %v", got, want)
	}
	if got, want := archRepoPackages(output, ""), map[string]struct{}{"bash": {}, "vim": {}, "app": {}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("all repo packages = %v; want %v", got, want)
	}
}

func TestDebianInstalledPackages(t *testing.T) {
	output := []byte("ii \tadduser\nrc \tremoved-package\nii  bash:amd64\n")
	got := debianInstalledPackages(output)
	if _, ok := got["removed-package"]; ok {
		t.Fatal("package in config-files-only state should not be considered installed")
	}
	for _, name := range []string{"adduser", "bash:amd64"} {
		if _, ok := got[name]; !ok {
			t.Errorf("installed package %q missing from result %v", name, got)
		}
	}
}

func TestFilterPackageNamesKeepsAvailableAndDeduplicates(t *testing.T) {
	selected := []string{"bash", "missing", "vim", "bash"}
	available := map[string]struct{}{"bash:amd64": {}, "vim": {}}
	found, missing := filterPackageNames(selected, available)
	if !reflect.DeepEqual(found, []string{"bash", "vim"}) {
		t.Errorf("found = %v; want [bash vim]", found)
	}
	if !reflect.DeepEqual(missing, []string{"missing"}) {
		t.Errorf("missing = %v; want [missing]", missing)
	}
}
