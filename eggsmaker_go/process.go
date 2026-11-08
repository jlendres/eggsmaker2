package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	if os.Geteuid() == 0 {
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

// runUserLogged ejecuta tareas del perfil del usuario (yay/Flatpak --user)
// sin elevarlas a root. Si la aplicación fue iniciada como root, vuelve al
// usuario original detectado mediante SUDO_USER.
func (s *AppState) runUserLogged(ctx context.Context, title, name string, args ...string) error {
	var cmd *exec.Cmd
	if os.Geteuid() == 0 && s.RealUser != "" && s.RealUser != "root" {
		all := append([]string{"-u", s.RealUser, "--", name}, args...)
		cmd = exec.CommandContext(ctx, "sudo", all...)
	} else {
		cmd = exec.CommandContext(ctx, name, args...)
	}
	return s.streamCommand(ctx, title, cmd)
}
