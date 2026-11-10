package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const packageBatchScript = `operation=$1
manager=$2
shift 2
status=0
for package do
  printf '[Gestor] Procesando %s\n' "$package"
  case "$operation" in
    pacman-install) "$manager" -S --needed --noconfirm -- "$package" ;;
    pacman-remove) "$manager" -Rns --noconfirm -- "$package" ;;
    apt-install) "$manager" install -y -- "$package" ;;
    apt-remove) "$manager" purge -y -- "$package" ;;
    chaotic-install) "$manager" -S --needed --noconfirm "chaotic-aur/$package" ;;
    chaotic-remove) "$manager" -Rns --noconfirm -- "$package" ;;
    *) printf '[Gestor] Operación desconocida: %s\n' "$operation" >&2; exit 2 ;;
  esac
  result=$?
  if [ "$result" -ne 0 ]; then
    printf '[Gestor] Error al procesar %s (código %s); se continúa con el siguiente.\n' "$package" "$result" >&2
    status=1
  fi
done
exit "$status"
`

func packageBatchCommand(ctx context.Context, operation, manager string, packages []string) (*exec.Cmd, error) {
	switch operation {
	case "pacman-install", "pacman-remove", "apt-install", "apt-remove", "chaotic-install", "chaotic-remove":
	default:
		return nil, fmt.Errorf("operación de paquetes no válida: %s", operation)
	}
	if len(packages) == 0 {
		return nil, fmt.Errorf("no hay paquetes para procesar")
	}
	args := []string{"-c", packageBatchScript, "eggsmaker-package-batch", operation, manager}
	args = append(args, packages...)
	return exec.CommandContext(ctx, "/bin/bash", args...), nil
}

func (s *AppState) runPrivilegedPackageBatch(ctx context.Context, title, operation string, packages []string) error {
	managerName := "pacman"
	switch operation {
	case "apt-install", "apt-remove":
		managerName = "apt-get"
	}
	manager, err := exec.LookPath(managerName)
	if err != nil {
		return fmt.Errorf("%s no está disponible: %w", managerName, err)
	}
	cmd, err := packageBatchCommand(ctx, operation, manager, packages)
	if err != nil {
		return err
	}
	if os.Geteuid() == 0 {
		return s.streamCommand(ctx, title, cmd)
	}
	if pkexec, err := exec.LookPath("pkexec"); err == nil {
		all := append([]string{cmd.Path}, cmd.Args[1:]...)
		return s.streamCommand(ctx, title, exec.CommandContext(ctx, pkexec, all...))
	}
	if _, err := exec.LookPath("sudo"); err != nil {
		return fmt.Errorf("se requiere pkexec o sudo para procesar paquetes")
	}
	terminalArgs := append([]string{"/bin/bash"}, cmd.Args[1:]...)
	quotedArgs := make([]string, len(terminalArgs))
	for i, arg := range terminalArgs {
		quotedArgs[i] = shellQuote(arg)
	}
	command := "sudo " + strings.Join(quotedArgs, " ")
	s.logf("\n── %s ──\nSe abrirá una terminal para autorizar y procesar el lote completo.\n", title)
	return s.runInteractiveAsUser(command, title)
}

func runEachPackage(action string, packages []string, run func(string) error, logf func(string, ...any)) error {
	failures := make([]error, 0)
	for _, name := range packages {
		logf("[Gestor] %s: %s\n", action, name)
		if err := run(name); err != nil {
			logf("[Gestor] Error al %s %s: %v\n", action, name, err)
			failures = append(failures, fmt.Errorf("%s: %w", name, err))
		}
	}
	if len(failures) == 0 {
		logf("[Gestor] Finalizó %s para los %d paquetes seleccionados.\n", action, len(packages))
		return nil
	}
	logf("[Gestor] Finalizó %s: %d de %d paquetes tuvieron errores; se intentaron los demás.\n", action, len(failures), len(packages))
	return errors.Join(failures...)
}

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
		if name == "flatpak" {
			cmd, err := s.userCommand(ctx, name, args...)
			if err != nil {
				return nil, err
			}
			return cmd.Output()
		}
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
	alreadyInstalledCount := 0
	if kind == "flatpak" && installing && len(found) > 0 {
		cmd, err := s.userCommand(ctx, "flatpak", "list", "--app", "--columns=application")
		if err != nil {
			return nil, fmt.Errorf("no se pudieron consultar las aplicaciones Flatpak instaladas: %w", err)
		}
		installedOutput, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("no se pudieron consultar las aplicaciones Flatpak instaladas: %w", err)
		}
		var alreadyInstalled []string
		alreadyInstalled, found = filterPackageNames(found, packageNamesFromOutput(installedOutput))
		alreadyInstalledCount = len(alreadyInstalled)
		for _, name := range alreadyInstalled {
			s.logf("[Gestor] Se omite %s: ya está instalado.\n", name)
		}
	}
	action := "desinstalación"
	if installing {
		action = "instalación"
	}
	for _, name := range missing {
		s.logf("[Gestor] Se omite %s: no está disponible para %s.\n", name, action)
	}
	if skipped := len(missing) + alreadyInstalledCount; skipped > 0 {
		s.logf("[Gestor] Omitidos %d de %d paquetes; se procesarán los %d disponibles.\n", skipped, len(selected), len(found))
	}
	if len(found) == 0 && len(missing) > 0 && alreadyInstalledCount == 0 {
		s.logf("[Gestor] No hay paquetes disponibles para %s.\n", action)
	}
	return found, nil
}
