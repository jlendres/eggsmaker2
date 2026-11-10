package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func packageNamesFromOutput(output []byte) map[string]struct{} {
	names := make(map[string]struct{})
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		names[fields[len(fields)-1]] = struct{}{}
	}
	return names
}

func archRepoPackages(output []byte, repo string) map[string]struct{} {
	names := make(map[string]struct{})
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || (repo != "" && fields[0] != repo) {
			continue
		}
		names[fields[1]] = struct{}{}
	}
	return names
}

func debianInstalledPackages(output []byte) map[string]struct{} {
	names := make(map[string]struct{})
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || !strings.HasPrefix(fields[0], "ii") {
			continue
		}
		names[fields[1]] = struct{}{}
	}
	return names
}

func filterPackageNames(selected []string, available map[string]struct{}) (found, missing []string) {
	normalizedAvailable := make(map[string]struct{}, len(available))
	for name := range available {
		normalizedAvailable[packageNameKey(name)] = struct{}{}
	}
	seen := make(map[string]struct{}, len(selected))
	for _, name := range selected {
		if _, duplicate := seen[name]; duplicate {
			continue
		}
		seen[name] = struct{}{}
		if _, ok := normalizedAvailable[packageNameKey(name)]; ok {
			found = append(found, name)
		} else {
			missing = append(missing, name)
		}
	}
	return found, missing
}

func packageNameKey(name string) string {
	if base, _, ok := strings.Cut(name, ":"); ok {
		return base
	}
	return name
}

func (s *AppState) availablePackageNames(ctx context.Context, kind string, installing bool, selected []string) (map[string]struct{}, error) {
	run := func(name string, args ...string) ([]byte, error) {
		return exec.CommandContext(ctx, name, args...).Output()
	}

	switch {
	case kind == "flatpak":
		if installing {
			output, err := run("flatpak", "remote-ls", "--app", "--columns=application")
			if err != nil {
				return nil, fmt.Errorf("no se pudo consultar las aplicaciones Flatpak disponibles: %w", err)
			}
			return packageNamesFromOutput(output), nil
		}
		output, err := run("flatpak", "list", "--app", "--columns=application")
		if err != nil {
			return nil, fmt.Errorf("no se pudieron consultar las aplicaciones Flatpak instaladas: %w", err)
		}
		return packageNamesFromOutput(output), nil
	case s.Distro == DistroArch && kind == "aur" && installing:
		if _, err := exec.LookPath("yay"); err != nil {
			return nil, fmt.Errorf("no se encontró yay")
		}
		available := make(map[string]struct{})
		for _, name := range selected {
			output, err := exec.CommandContext(ctx, "yay", "-Si", name).CombinedOutput()
			if err == nil {
				available[name] = struct{}{}
				continue
			}
			s.logf("[Gestor] Se omite %s de AUR: no se pudo verificar el paquete (%v)", name, err)
			if detail := strings.TrimSpace(string(output)); detail != "" {
				s.logf(" — %s", detail)
			}
			s.logf("\n")
		}
		return available, nil
	case s.Distro == DistroArch:
		if installing {
			output, err := run("pacman", "-Sl")
			if err != nil {
				return nil, fmt.Errorf("no se pudieron consultar los paquetes disponibles en los repositorios: %w", err)
			}
			if kind == "chaotic" {
				return archRepoPackages(output, "chaotic-aur"), nil
			}
			available := archRepoPackages(output, "")
			if kind == "nativa" {
				delete(available, "")
			}
			return available, nil
		}
		output, err := run("pacman", "-Qq")
		if err != nil {
			return nil, fmt.Errorf("no se pudieron consultar los paquetes instalados: %w", err)
		}
		return packageNamesFromOutput(output), nil
	case s.Distro == DistroDebian:
		if installing {
			output, err := run("apt-cache", "pkgnames")
			if err != nil {
				return nil, fmt.Errorf("no se pudieron consultar los paquetes disponibles en APT: %w", err)
			}
			return packageNamesFromOutput(output), nil
		}
		output, err := run("dpkg-query", "-W", "-f=${db:Status-Abbrev}\\t${binary:Package}\\n")
		if err != nil {
			return nil, fmt.Errorf("no se pudieron consultar los paquetes instalados: %w", err)
		}
		return debianInstalledPackages(output), nil
	default:
		return nil, fmt.Errorf("familia de distribución no compatible para el gestor de paquetes")
	}
}

func (s *AppState) filterSelectedPackages(ctx context.Context, kind string, installing bool, selected []string) ([]string, error) {
	available, err := s.availablePackageNames(ctx, kind, installing, selected)
	if err != nil {
		return nil, err
	}
	found, missing := filterPackageNames(selected, available)
	action := "desinstalación"
	if installing {
		action = "instalación"
	}
	for _, name := range missing {
		s.logf("[Gestor] Se omite %s: no está disponible para %s.\n", name, action)
	}
	if len(missing) > 0 {
		s.logf("[Gestor] Omitidos %d de %d paquetes; se procesarán los %d disponibles.\n", len(missing), len(selected), len(found))
	}
	if len(found) == 0 {
		s.logf("[Gestor] No hay paquetes disponibles para %s.\n", action)
	}
	return found, nil
}
