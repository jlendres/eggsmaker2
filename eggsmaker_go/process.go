package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// commandArgs mantiene los argumentos separados: no se interpolan nombres de
// archivos o paquetes dentro de una shell.
func (s *AppState) privilegedCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	if name == "eggs" && s.EggsPath != "" {
		name = s.EggsPath
	}
	if os.Geteuid() == 0 {
		return exec.CommandContext(ctx, name, args...)
	}
	if _, err := exec.LookPath("pkexec"); err == nil {
		all := append([]string{name}, args...)
		return exec.CommandContext(ctx, "pkexec", all...)
	}
	all := append([]string{name}, args...)
	return exec.CommandContext(ctx, "sudo", append([]string{"-n"}, all...)...)
}

func (s *AppState) runLogged(ctx context.Context, title, name string, args ...string) error {
	return s.streamCommand(ctx, title, s.privilegedCommand(ctx, name, args...))
}

func (s *AppState) streamCommand(ctx context.Context, title string, cmd *exec.Cmd) error {
	s.logf("\n── %s ──\n$ %s\n", title, strings.Join(cmd.Args, " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		s.logf("[ERROR] No se pudo iniciar: %v\n", err)
		return err
	}
	output := make(chan []byte, 128)
	read := func(r io.Reader) {
		buf := make([]byte, 4096)
		for {
			n, readErr := r.Read(buf)
			if n > 0 {
				output <- append([]byte(nil), buf[:n]...)
			}
			if readErr != nil {
				return
			}
		}
	}
	var readers sync.WaitGroup
	readers.Add(2)
	go func() { defer readers.Done(); read(stdout) }()
	go func() { defer readers.Done(); read(stderr) }()
	done := make(chan error, 1)
	go func() { readers.Wait(); done <- cmd.Wait() }()
	for {
		select {
		case chunk := <-output:
			s.logf("%s", chunk)
		case err := <-done:
			for {
				select {
				case chunk := <-output:
					s.logf("%s", chunk)
				default:
					if err != nil {
						s.logf("[ERROR] Proceso finalizado: %v\n", err)
					} else {
						s.logf("[OK] Proceso finalizado correctamente.\n")
					}
					return err
				}
			}
		case <-ctx.Done():
			_ = cmd.Process.Kill()
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (s *AppState) runInteractive(command, title string) error {
	return s.runInteractiveCommand(command, title, true)
}

func (s *AppState) runInteractiveAsUser(command, title string) error {
	return s.runInteractiveCommand(command, title, false)
}

func (s *AppState) runInteractiveCommand(command, title string, elevate bool) error {
	terms := []string{"konsole", "xfce4-terminal", "gnome-terminal", "kitty", "alacritty", "xterm", "x-terminal-emulator"}
	term := ""
	for _, candidate := range terms {
		if _, err := exec.LookPath(candidate); err == nil {
			term = candidate
			break
		}
	}
	if term == "" {
		err := fmt.Errorf("no se encontró una terminal gráfica para ejecutar %s", command)
		s.logf("[ERROR] %v\n", err)
		return err
	}
	var args []string
	shell := ""
	if !elevate || os.Geteuid() == 0 {
		shell = command
	} else if _, err := exec.LookPath("pkexec"); err == nil {
		shell = "pkexec /bin/bash -lc " + shellQuote(command)
	} else {
		shell = "sudo /bin/bash -lc " + shellQuote(command)
	}
	wrapped := shell + "; rc=$?; echo; echo 'Proceso finalizado (código '$rc').'; read -r -p 'Presione Enter para cerrar...' _; exit $rc"
	switch term {
	case "konsole":
		args = []string{"--nofork", "--workdir", ".", "-p", "tabtitle=" + title, "-e", "/bin/bash", "-lc", wrapped}
	case "xfce4-terminal":
		args = []string{"--title=" + title, "--disable-server", "-e", "/bin/bash -lc " + shellQuote(wrapped)}
	case "gnome-terminal":
		args = []string{"--wait", "--title=" + title, "--", "/bin/bash", "-lc", wrapped}
	case "kitty":
		args = []string{"--title", title, "/bin/bash", "-lc", wrapped}
	case "alacritty":
		args = []string{"--title", title, "-e", "/bin/bash", "-lc", wrapped}
	case "x-terminal-emulator":
		args = []string{"-e", "/bin/bash -lc " + shellQuote(wrapped)}
	default:
		args = []string{"-T", title, "-e", "/bin/bash", "-lc", wrapped}
	}
	s.logf("\n── %s ──\nSe abrió una terminal interactiva para: %s\n", title, command)
	if err := exec.Command(term, args...).Run(); err != nil {
		s.logf("[ERROR] %v\n", err)
		return err
	}
	return nil
}

func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }

// userCommand ejecuta una tarea dentro del perfil del usuario de escritorio.
func (s *AppState) userCommand(ctx context.Context, name string, args ...string) (*exec.Cmd, error) {
	if os.Geteuid() == 0 && s.RealUser != "" && s.RealUser != "root" {
		account, err := user.Lookup(s.RealUser)
		if err != nil {
			return nil, fmt.Errorf("no se pudo identificar al usuario %s: %w", s.RealUser, err)
		}
		runtimeDir := filepath.Join("/run/user", account.Uid)
		environment := []string{
			"HOME=" + account.HomeDir,
			"USER=" + account.Username,
			"LOGNAME=" + account.Username,
			"XDG_RUNTIME_DIR=" + runtimeDir,
			"DBUS_SESSION_BUS_ADDRESS=unix:path=" + filepath.Join(runtimeDir, "bus"),
		}
		all := append([]string{"-u", account.Username, "--", "env"}, environment...)
		all = append(all, name)
		all = append(all, args...)
		return exec.CommandContext(ctx, "runuser", all...), nil
	}
	return exec.CommandContext(ctx, name, args...), nil
}

// runUserLogged ejecuta tareas del perfil del usuario (yay/Flatpak --user)
// sin elevarlas a root. Si la aplicación fue iniciada como root, vuelve al
// usuario original detectado mediante SUDO_USER.
func (s *AppState) runUserLogged(ctx context.Context, title, name string, args ...string) error {
	cmd, err := s.userCommand(ctx, name, args...)
	if err != nil {
		return err
	}
	return s.streamCommand(ctx, title, cmd)
}

func (s *AppState) runUserInteractive(title, name string, args ...string) error {
	cmd, err := s.userCommand(context.Background(), name, args...)
	if err != nil {
		return err
	}
	quotedArgs := make([]string, len(cmd.Args))
	for i, arg := range cmd.Args {
		quotedArgs[i] = shellQuote(arg)
	}
	return s.runInteractiveAsUser(strings.Join(quotedArgs, " "), title)
}
