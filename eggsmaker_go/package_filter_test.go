package main

import (
	"errors"
	"fmt"
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
