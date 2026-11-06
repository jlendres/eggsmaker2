package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

const AppVersion = "26.11.1"

type DistroFamily string

const (
	DistroArch    DistroFamily = "arch"
	DistroDebian  DistroFamily = "debian"
	DistroUnknown DistroFamily = "unknown"
)

type AppState struct {
	Distro      DistroFamily
	DistroName  string
	EggsVersion string
	EggsPath    string
	RealUser    string
	RealHome    string
	HostName    string
	TargetDir   string
	Window      fyne.Window
	Log         *LogView
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("Eggsmaker Go %s\n", AppVersion)
		return
	}
	state := &AppState{}
	state.detectDistro()
	state.initUserAndPaths()
	state.validateEnvironment()

	application := app.NewWithID("org.eggsmaker.go")
	application.Settings().SetTheme(eggsmakerTheme{Theme: theme.DefaultTheme()})
	state.Window = application.NewWindow("Eggsmaker Go")
	state.buildGUI(application)
}

func (s *AppState) detectDistro() {
	s.Distro = DistroUnknown
	f, err := os.Open("/etc/os-release")
	if err == nil {
		defer f.Close()
		values := map[string]string{}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			key, value, ok := strings.Cut(line, "=")
			if ok {
				values[key] = strings.Trim(strings.Trim(value, `"`), `'`)
			}
		}
		s.DistroName = values["PRETTY_NAME"]
		check := strings.ToLower(values["ID"] + " " + values["ID_LIKE"])
		if regexp.MustCompile(`(^|[[:space:]])(arch|archlinux|cachyos|endeavouros|manjaro|garuda|arcolinux|artix|archcraft)([[:space:]]|$)`).MatchString(check) {
			s.Distro = DistroArch
		} else if regexp.MustCompile(`(^|[[:space:]])(debian|ubuntu|linuxmint|mint|pop|biglinux|kali|mx|zorin|elementary|neon|raspbian|deepin)([[:space:]]|$)`).MatchString(check) {
			s.Distro = DistroDebian
		}
	}
	if s.Distro == DistroUnknown {
		if _, err := exec.LookPath("pacman"); err == nil {
			s.Distro = DistroArch
		} else if _, err := exec.LookPath("apt-get"); err == nil {
			s.Distro = DistroDebian
		}
	}
	if s.DistroName == "" {
		s.DistroName = string(s.Distro)
	}
}

func (s *AppState) initUserAndPaths() {
	// La GUI se ejecuta como usuario normal; SUDO_USER permite conservar el
	// destino correcto si, excepcionalmente, se inició desde sudo.
	s.RealUser = os.Getenv("SUDO_USER")
	if s.RealUser == "" {
		s.RealUser = os.Getenv("USER")
	}
	home := os.Getenv("HOME")
	if os.Geteuid() == 0 && s.RealUser != "" {
		if out, err := exec.Command("getent", "passwd", s.RealUser).Output(); err == nil {
			fields := strings.Split(strings.TrimSpace(string(out)), ":")
			if len(fields) > 5 {
				home = fields[5]
			}
		}
	}
	if home == "" {
		home = filepath.Join("/home", s.RealUser)
	}
	s.RealHome = home
	s.HostName, _ = os.Hostname()
	s.TargetDir = filepath.Join(s.RealHome, s.HostName+"-pak")
	if err := os.MkdirAll(s.TargetDir, 0755); err != nil {
		s.TargetDir = os.TempDir()
	}
}

func (s *AppState) validateEnvironment() {
	eggsPath, pathErr := exec.LookPath("eggs")
	if pathErr != nil {
		s.EggsVersion = "no instalado"
		return
	}
	s.EggsPath = eggsPath
	out, err := exec.Command(eggsPath, "--version").CombinedOutput()
	if err == nil {
		fields := strings.Fields(string(out))
		if len(fields) > 0 {
			s.EggsVersion = fields[len(fields)-1]
		}
	}
	if s.EggsVersion == "" {
		s.EggsVersion = "desconocida"
	}
}

func (s *AppState) logf(format string, args ...any) {
	if s.Log != nil {
		s.Log.Append(fmt.Sprintf(format, args...))
	}
}

// eggsmakerTheme mejora el contraste del visor de procesos y usa violeta como
// color principal, incluido el indicador de generación de ISO.
type eggsmakerTheme struct{ fyne.Theme }

func (t eggsmakerTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameDisabled {
		return color.NRGBA{R: 0, G: 220, B: 255, A: 255}
	}
	if name == theme.ColorNamePrimary {
		return color.NRGBA{R: 155, G: 75, B: 220, A: 255}
	}
	return t.Theme.Color(name, variant)
}
