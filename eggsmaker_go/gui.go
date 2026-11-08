package main

import (
	"context"
	_ "embed"
	"fmt"
	"image/color"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed eggsmaker.png
var eggsmakerLogoPNG []byte

type LogView struct {
	entry  *widget.Entry
	scroll *container.Scroll
}

func (l *LogView) Append(v string) {
	if l == nil || l.entry == nil {
		return
	}
	fyne.Do(func() {
		l.entry.SetText(l.entry.Text + v)
		l.entry.Refresh()
		if l.scroll != nil {
			l.scroll.ScrollToBottom()
		}
	})
}

func (s *AppState) buildGUI(myApp fyne.App) {
	s.Window.Resize(fyne.NewSize(900, 680))
	s.Window.SetMaster()
	s.Window.SetContent(s.buildMainContent(myApp))
	s.Window.ShowAndRun()
}

func (s *AppState) buildMainContent(myApp fyne.App) fyne.CanvasObject {
	oldLog := ""
	if s.Log != nil && s.Log.entry != nil {
		oldLog = s.Log.entry.Text
	}
	s.themeRefresh = nil
	header := canvas.NewText(tr(s.Language, "Eggsmaker Go")+"  •  v"+AppVersion, color.Transparent)
	s.bindColor(colorNameAccentText, func(c color.Color) { header.Color = c; header.Refresh() })
	header.TextSize = 22
	header.TextStyle = fyne.TextStyle{Bold: true}
	header.Alignment = fyne.TextAlignCenter
	logo := canvas.NewImageFromResource(fyne.NewStaticResource("eggsmaker.png", eggsmakerLogoPNG))
	logo.FillMode = canvas.ImageFillContain
	logoBox := container.NewGridWrap(fyne.NewSize(44, 44), logo)
	appearanceBtn := widget.NewButtonWithIcon(tr(s.Language, "Apariencia"), theme.ColorPaletteIcon(), func() { s.showAppearanceDialog(myApp) })
	headerRow := container.NewBorder(nil, nil, nil, container.NewHBox(appearanceBtn, widget.NewLabel(tr(s.Language, "Idioma:")), func() fyne.CanvasObject {
		langSelect := widget.NewSelect([]string{"Español", "English", "Português", "Italiano"}, nil)
		langNames := map[string]string{"es": "Español", "en": "English", "pt": "Português", "it": "Italiano"}
		langCodes := map[string]string{"Español": "es", "English": "en", "Português": "pt", "Italiano": "it"}
		if s.Language == "" {
			s.Language = "es"
		}
		langSelect.SetSelected(langNames[s.Language])
		langSelect.OnChanged = func(name string) {
			if code, ok := langCodes[name]; ok && code != s.Language {
				s.Language = code
				myApp.Preferences().SetString("language", code)
				s.Window.SetContent(s.buildMainContent(myApp))
			}
		}
		return langSelect
	}()), container.NewCenter(container.NewHBox(logoBox, header)))
	status := widget.NewLabel(fmt.Sprintf(tr(s.Language, "Penguins Eggs: %s   |   Distribución: %s"), s.EggsVersion, s.DistroName))
	if _, err := exec.LookPath("eggs"); err != nil {
		status.SetText(status.Text + tr(s.Language, "   |   ADVERTENCIA: eggs no está instalado"))
	}
	console := widget.NewMultiLineEntry()
	console.SetMinRowsVisible(12)
	console.Wrapping = fyne.TextWrapOff
	console.TextStyle = fyne.TextStyle{Monospace: true}
	console.SetText(oldLog)
	if oldLog == "" {
		console.SetText(tr(s.Language, "Eggsmaker Go iniciado.\nLos registros de los procesos aparecerán aquí.\n"))
	}
	console.Disable()
	consoleScroll := container.NewScroll(console)
	s.Log = &LogView{entry: console, scroll: consoleScroll}
	progress := newColorProgress(nil)
	progress.Hide()
	progressBusy := widget.NewProgressBarInfinite()
	progressBusy.Hide()
	progressText := widget.NewLabel(tr(s.Language, "Listo"))
	isoInfo := canvas.NewText("", color.Transparent)
	s.bindColor(colorNameAccentText, func(c color.Color) { isoInfo.Color = c; isoInfo.Refresh() })
	isoInfo.TextSize = 14
	isoInfo.TextStyle = fyne.TextStyle{Bold: true}

	var btnISO *taskButton
	btnISO = newTaskButton(tr(s.Language, "Generar ISO personalizada"), func() { s.startISOWorkflow(progress, progressBusy, progressText, isoInfo, btnISO) })
	btnLive := newTaskButton(tr(s.Language, "Instalación Live (Calamares / Krill)"), func() { s.showInstallerDialog() })
	btnPackages := newTaskButton(tr(s.Language, "Gestor de paquetes"), func() { s.showPackageManager() })
	btnConfig := newTaskButton(tr(s.Language, "Editar configuración de Penguins Eggs"), func() {
		s.showEggsConfigDialog(nil)
	})
	btnClear := widget.NewButton(tr(s.Language, "Limpiar visor"), func() { console.SetText("") })
	btnExit := widget.NewButton(tr(s.Language, "Salir"), func() { s.Window.Close() })

	buttons := container.New(layout.NewGridLayout(2), btnISO, btnLive, btnPackages, btnConfig)
	logHeader := container.NewBorder(nil, nil, widget.NewLabel(tr(s.Language, "Terminal / salida de procesos")), btnClear)
	consoleFrame := canvas.NewRectangle(color.Transparent)
	s.bindColor(colorNameAccentText, func(c color.Color) { consoleFrame.StrokeColor = c; consoleFrame.Refresh() })
	consoleFrame.StrokeWidth = 1.5
	consoleOverride := container.NewThemeOverride(consoleScroll, consoleTheme{
		Theme:      theme.Current(),
		background: themeColor(colorNameConsoleBackground),
		foreground: themeColor(colorNameAccentText),
	})
	s.themeRefresh = append(s.themeRefresh, func() {
		consoleOverride.Theme = consoleTheme{
			Theme:      theme.Current(),
			background: themeColor(colorNameConsoleBackground),
			foreground: themeColor(colorNameAccentText),
		}
		consoleOverride.Refresh()
	})
	consoleView := container.NewStack(consoleFrame, container.NewPadded(consoleOverride))
	content := container.NewBorder(
		container.NewVBox(headerRow, status, widget.NewSeparator(), buttons, widget.NewSeparator(), progressText, isoInfo, progressBusy, progress, logHeader),
		container.NewHBox(layout.NewSpacer(), btnExit), nil, nil,
		consoleView,
	)
	if oldLog == "" {
		s.logf("Sistema detectado: %s (%s)\nUsuario: %s\nDirectorio de listas: %s\n", s.Distro, s.DistroName, s.RealUser, s.TargetDir)
	}
	return content
}

func (s *AppState) eggBinary() string {
	if s.EggsPath != "" {
		return s.EggsPath
	}
	return "eggs"
}

func (s *AppState) showConfirmNoDefault(title, message string, callback func(bool)) {
	var d *dialog.CustomDialog
	finished := false
	finish := func(answer bool) {
		if finished {
			return
		}
		finished = true
		d.Hide()
		callback(answer)
	}
	yes := widget.NewButton(tr(s.Language, "Sí"), func() { finish(true) })
	yes.Importance = widget.LowImportance
	no := widget.NewButton(tr(s.Language, "No"), func() { finish(false) })
	no.Importance = widget.HighImportance
	buttons := container.NewHBox(layout.NewSpacer(), yes, no)
	d = dialog.NewCustomWithoutButtons(title, container.NewVBox(widget.NewLabel(message), buttons), s.Window)
	d.Resize(fyne.NewSize(460, d.MinSize().Height))
	d.Show()
}

func (s *AppState) startISOWorkflow(progress *colorProgress, busy *widget.ProgressBarInfinite, label *widget.Label, isoInfo *canvas.Text, button *taskButton) {
	if _, err := exec.LookPath("eggs"); err != nil {
		dialog.ShowError(fmt.Errorf("Penguins Eggs no está instalado o no está en PATH"), s.Window)
		return
	}
	s.showConfirmNoDefault(tr(s.Language, "Configuración de escritorio"), tr(s.Language, "¿Deseás copiar la configuración del escritorio con 'eggs tools skel'?"), func(skel bool) {
		continueClone := func() {
			s.showConfirmNoDefault(tr(s.Language, "Clonación"), tr(s.Language, "¿Deseás crear un clon del sistema (--clone)?"), func(clone bool) {
				s.showConfirmNoDefault(tr(s.Language, "Confirmar generación"), tr(s.Language, "Se limpiarán los temporales de Eggs y se iniciará la creación de la ISO. ¿Continuar?"), func(ok bool) {
					if ok {
						go s.generateISO(clone, progress, busy, label, isoInfo, button)
					}
				})
			})
		}
		askConfig := func() {
			s.showConfirmNoDefault(tr(s.Language, "Configuración de Penguins Eggs"), tr(s.Language, "¿Deseás editar la configuración con 'eggs config'? Los cambios se guardan desde el propio asistente."), func(edit bool) {
				if !edit {
					continueClone()
					return
				}
				s.showEggsConfigDialog(continueClone)
			})
		}
		if skel {
			go func() {
				err := s.runLogged(context.Background(), "Copiar configuración de escritorio", "eggs", "tools", "skel")
				if err != nil {
					fyne.Do(func() { dialog.ShowError(err, s.Window) })
					return
				}
				fyne.Do(func() { askConfig() })
			}()
		} else {
			askConfig()
		}
	})
}

func (s *AppState) generateISO(clone bool, progress *colorProgress, busy *widget.ProgressBarInfinite, label *widget.Label, isoInfo *canvas.Text, button *taskButton) {
	fyne.Do(func() { button.busy = true; button.setActive(true) })
	defer fyne.Do(func() { button.busy = false; button.setActive(false) })
	fyne.Do(func() {
		busy.Show()
		busy.Start()
		label.SetText(tr(s.Language, "Generando ISO: progreso indeterminado; se muestran etapas y salida real."))
	})
	args := []string{"remaster"}
	if clone {
		args = append(args, "--clone")
	}
	eggsBin := s.EggsPath
	if eggsBin == "" {
		eggsBin = "eggs"
	}
	command := shellQuote(eggsBin) + " kill && " + shellQuote(eggsBin) + " tools clean && " + shellQuote(eggsBin) + " " + strings.Join(args, " ")
	startedAt := time.Now()
	fyne.Do(func() { isoInfo.Text = tr(s.Language, "Generando ISO..."); isoInfo.Refresh() })
	ctx := context.Background()
	err := s.runLogged(ctx, "Preparar y generar ISO", "/bin/bash", "-lc", command)
	fyne.Do(func() {
		busy.Stop()
		busy.Hide()
		if err != nil {
			label.SetText(tr(s.Language, "La generación terminó con errores"))
		} else {
			label.SetText(tr(s.Language, "Generación finalizada"))
		}
	})
	if err != nil {
		fyne.Do(func() { dialog.ShowError(err, s.Window) })
		return
	}
	iso, err := latestISO("/home/eggs")
	if err == nil {
		if info, e := os.Stat(iso); e != nil {
			err = e
		} else if info.ModTime().Before(startedAt) {
			err = fmt.Errorf("no se detectó una ISO nueva después de iniciar la generación")
		}
	}
	if err != nil {
		fyne.Do(func() {
			isoInfo.Text = fmt.Sprintf(tr(s.Language, "Tiempo de procesamiento: %s"), time.Since(startedAt).Round(time.Second))
			isoInfo.Refresh()
			dialog.ShowInformation(tr(s.Language, "Generación finalizada"), tr(s.Language, "No se encontró un archivo ISO en /home/eggs. Revisá el visor de salida."), s.Window)
		})
		return
	}
	info, statErr := os.Stat(iso)
	if statErr != nil {
		err = statErr
		fyne.Do(func() { dialog.ShowError(err, s.Window) })
		return
	}
	elapsed := time.Since(startedAt).Round(time.Second)
	isoDetails := fmt.Sprintf(tr(s.Language, "Tiempo de procesamiento: %s    |    Tamaño de la ISO: %s"), elapsed, formatBytes(info.Size()))
	fyne.Do(func() {
		isoInfo.Text = isoDetails
		isoInfo.Refresh()
		s.showConfirmNoDefault(tr(s.Language, "ISO generada"), tr(s.Language, "Archivo encontrado:\n")+iso+"\n\n"+isoDetails+"\n\n¿Deseás copiarlo a otra ubicación?", func(copyIt bool) {
			if copyIt {
				s.chooseCopyDestination(iso, progress, label)
			} else {
				s.askCleanup()
			}
		})
	})
}

func latestISO(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	type item struct {
		path string
		mod  time.Time
	}
	var isos []item
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".iso") {
			info, err := e.Info()
			if err == nil {
				isos = append(isos, item{filepath.Join(dir, e.Name()), info.ModTime()})
			}
		}
	}
	if len(isos) == 0 {
		return "", fmt.Errorf("no hay archivos ISO en %s", dir)
	}
	sort.Slice(isos, func(i, j int) bool { return isos[i].mod.After(isos[j].mod) })
	return isos[0].path, nil
}

func (s *AppState) chooseCopyDestination(source string, progress *colorProgress, label *widget.Label) {
	fileDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, s.Window)
			return
		}
		if uri == nil {
			s.askCleanup()
			return
		}
		dest := filepath.Join(uri.Path(), filepath.Base(source))
		if filepath.Clean(dest) == filepath.Clean(source) {
			dialog.ShowInformation(tr(s.Language, "Copiar ISO"), tr(s.Language, "El destino seleccionado es el mismo archivo de origen."), s.Window)
			return
		}
		if _, statErr := os.Stat(dest); statErr == nil {
			s.showConfirmNoDefault(tr(s.Language, "Sobrescribir archivo"), tr(s.Language, "Ya existe ")+dest+". ¿Deseás reemplazarlo?", func(ok bool) {
				if ok {
					go s.copyWithProgress(source, dest, progress, label)
				}
			})
			return
		}
		go s.copyWithProgress(source, dest, progress, label)
	}, s.Window)
	fileDialog.Resize(fyne.NewSize(900, 650))
	fileDialog.Show()
}

func (s *AppState) copyWithProgress(source, dest string, progress *colorProgress, label *widget.Label) {
	in, err := os.Open(source)
	if err != nil {
		fyne.Do(func() { dialog.ShowError(err, s.Window) })
		return
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		fyne.Do(func() { dialog.ShowError(err, s.Window) })
		return
	}
	out, err := os.CreateTemp(filepath.Dir(dest), ".eggsmaker-copy-*.part")
	if err != nil {
		fyne.Do(func() { dialog.ShowError(err, s.Window) })
		return
	}
	tempPath := out.Name()
	defer os.Remove(tempPath)
	fyne.Do(func() { progress.Show(); progress.SetValue(0); label.SetText(tr(s.Language, "Copiando ISO...")) })
	buf := make([]byte, 4*1024*1024)
	var copied int64
	last := time.Now()
	for {
		n, readErr := in.Read(buf)
		if n > 0 {
			written, werr := out.Write(buf[:n])
			copied += int64(written)
			if werr != nil {
				err = werr
				break
			}
			if written != n {
				err = io.ErrShortWrite
				break
			}
			if info.Size() > 0 && (time.Since(last) > 250*time.Millisecond || copied == info.Size()) {
				pct := float64(copied) / float64(info.Size())
				copiedNow, totalNow := copied, info.Size()
				fyne.Do(func() {
					progress.SetValue(pct)
					label.SetText(fmt.Sprintf(tr(s.Language, "Copiando: %.1f%%  •  %s / %s"), pct*100, formatBytes(copiedNow), formatBytes(totalNow)))
				})
				last = time.Now()
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			err = readErr
			break
		}
	}
	if err == nil {
		err = out.Sync()
	}
	closeErr := out.Close()
	if err == nil {
		err = closeErr
	}
	if err == nil {
		err = os.Rename(tempPath, dest)
	}
	fyne.Do(func() {
		progress.Hide()
		if err != nil {
			label.SetText(tr(s.Language, "Error durante la copia"))
		} else {
			progress.SetValue(1)
			label.SetText(tr(s.Language, "Copia completada: ") + dest)
		}
	})
	if err != nil {
		fyne.Do(func() { dialog.ShowError(err, s.Window) })
		return
	}
	fyne.Do(func() {
		s.showConfirmNoDefault(tr(s.Language, "Copia completada"), tr(s.Language, "¿Deseás copiar la ISO a otra ubicación?"), func(again bool) {
			if again {
				s.chooseCopyDestination(source, progress, label)
			} else {
				s.askCleanup()
			}
		})
	})
}

func formatBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}

func (s *AppState) askCleanup() {
	s.showConfirmNoDefault(tr(s.Language, "Limpieza"), tr(s.Language, "¿Deseás eliminar los archivos temporales de construcción con 'eggs kill'?"), func(ok bool) {
		if !ok {
			return
		}
		go func() {
			err := s.runLogged(context.Background(), "Limpiar temporales", "eggs", "kill")
			fyne.Do(func() {
				if err != nil {
					dialog.ShowError(err, s.Window)
					return
				}
				dialog.ShowInformation(tr(s.Language, "Limpieza finalizada"), tr(s.Language, "El proceso de limpieza de archivos temporales finalizó correctamente."), s.Window)
			})
		}()
	})
}

func (s *AppState) showInstallerDialog() {
	choices := []string{tr(s.Language, "Calamares (interfaz gráfica)"), tr(s.Language, "Krill (terminal interactiva)")}
	selectBox := widget.NewSelect(choices, nil)
	selectBox.SetSelected(choices[0])
	dialog.ShowCustomConfirm(tr(s.Language, "Instalación Live"), tr(s.Language, "Ejecutar"), tr(s.Language, "Cancelar"), container.NewVBox(widget.NewLabel(tr(s.Language, "Seleccioná el instalador:")), selectBox), func(ok bool) {
		if !ok {
			return
		}
		if selectBox.Selected == choices[0] {
			go func() { _ = s.runLogged(context.Background(), "Calamares", "eggs", "sysinstall", "calamares") }()
		} else {
			if err := s.runInteractive(shellQuote(s.eggBinary())+" sysinstall krill", "Eggsmaker Go - Krill"); err != nil {
				dialog.ShowError(err, s.Window)
			}
		}
	}, s.Window)
}

func (s *AppState) showPackageManager() {
	if s.Distro != DistroArch && s.Distro != DistroDebian {
		dialog.ShowInformation(tr(s.Language, "Gestor de paquetes"), tr(s.Language, "Distribución no reconocida: sólo se admiten Debian y Arch."), s.Window)
		return
	}
	label := widget.NewLabel(tr(s.Language, "Las listas se guardan en ") + s.TargetDir + tr(s.Language, ". Primero recopilá las listas del sistema."))
	collect := widget.NewButton(tr(s.Language, "1. Recopilar listas"), func() {
		go func() {
			if err := s.collectPackageLists(); err != nil {
				fyne.Do(func() { dialog.ShowError(err, s.Window) })
				return
			}
			fyne.Do(func() { label.SetText(tr(s.Language, "Listas actualizadas en ") + s.TargetDir) })
		}()
	})
	install := widget.NewButton(tr(s.Language, "2. Instalar lista nativa"), func() { s.chooseInstallList("lista-nativo.txt", "nativa") })
	aurInstall := widget.NewButton(tr(s.Language, "Instalar lista AUR (Arch)"), func() { s.chooseInstallList("lista-aur.txt", "aur") })
	chaoticInstall := widget.NewButton(tr(s.Language, "Instalar lista Chaotic-AUR"), func() { s.chooseInstallList("lista-chaotic-aur.txt", "chaotic") })
	flatpakInstall := widget.NewButton(tr(s.Language, "Instalar lista Flatpak"), func() { s.chooseInstallList("lista-flatpak.txt", "flatpak") })
	remove := widget.NewButton(tr(s.Language, "3. Desinstalar paquetes"), func() { s.chooseRemovalList("lista-nativo.txt", tr(s.Language, "Paquetes del sistema"), false) })
	aurRemove := widget.NewButton(tr(s.Language, "4. Desinstalar paquetes AUR"), func() { s.chooseRemovalList("lista-aur.txt", tr(s.Language, "Paquetes AUR"), false) })
	chaoticRemove := widget.NewButton(tr(s.Language, "Desinstalar paquetes Chaotic-AUR"), func() { s.chooseRemovalList("lista-chaotic-aur.txt", tr(s.Language, "Paquetes Chaotic-AUR"), false) })
	flatpakRemove := widget.NewButton(tr(s.Language, "5. Desinstalar Flatpaks"), func() { s.chooseRemovalList("lista-flatpak.txt", tr(s.Language, "Aplicaciones Flatpak"), true) })
	clearLists := widget.NewButton(tr(s.Language, "Vaciar listas guardadas"), func() {
		s.showConfirmNoDefault(tr(s.Language, "Vaciar listas"), tr(s.Language, "Se eliminarán únicamente los archivos de listas de paquetes de Eggsmaker. ¿Continuar?"), func(ok bool) {
			if !ok {
				return
			}
			names := []string{"lista-nativo.txt", "lista-aur.txt", "lista-chaotic-aur.txt", "lista-flatpak.txt"}
			for _, name := range names {
				path := filepath.Join(s.TargetDir, name)
				_ = os.Remove(path)
				_ = os.Remove(selectionFilePath(path))
			}
			label.SetText(tr(s.Language, "Listas eliminadas."))
		})
	})
	openDir := widget.NewButton(tr(s.Language, "Abrir carpeta de listas"), func() {
		if _, err := exec.LookPath("xdg-open"); err == nil {
			_ = exec.Command("xdg-open", s.TargetDir).Start()
		}
	})
	content := container.NewVBox(label, collect, install, aurInstall, chaoticInstall, flatpakInstall, remove, aurRemove, chaoticRemove, flatpakRemove, clearLists, openDir,
		widget.NewLabel(tr(s.Language, "La instalación/restauración requiere listas compatibles con la misma familia de distribución.")))
	managerDialog := dialog.NewCustom(tr(s.Language, "Gestor de paquetes"), tr(s.Language, "Cerrar"), content, s.Window)
	managerDialog.Resize(fyne.NewSize(760, 700))
	managerDialog.Show()
}

func (s *AppState) collectPackageLists() error {
	if err := os.MkdirAll(s.TargetDir, 0755); err != nil {
		return err
	}
	run := func(name string, args ...string) ([]byte, error) { return exec.Command(name, args...).Output() }
	var native, aur, chaotic, flatpak []byte
	var err error
	if s.Distro == DistroArch {
		native, err = run("pacman", "-Qqe")
		if err != nil {
			return fmt.Errorf("pacman -Qqe: %w", err)
		}
		aur, _ = run("pacman", "-Qqem")
		// Clasificar paquetes explícitos del repositorio Chaotic-AUR.
		installed := map[string]bool{}
		for _, p := range strings.Fields(string(native)) {
			installed[p] = true
		}
		chaoticRepo, _ := run("pacman", "-Sl", "chaotic-aur")
		var chaoticNames []string
		for _, line := range strings.Split(string(chaoticRepo), "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 2 && fields[0] == "chaotic-aur" && installed[fields[1]] {
				chaoticNames = append(chaoticNames, fields[1])
			}
		}
		chaotic = []byte(strings.Join(chaoticNames, "\n"))
		if len(chaoticNames) > 0 {
			chaotic = append(chaotic, '\n')
		}
		// Separar AUR y Chaotic-AUR de la lista nativa.
		aurSet := map[string]bool{}
		for _, p := range strings.Fields(string(aur)) {
			aurSet[p] = true
		}
		chaoticSet := map[string]bool{}
		for _, p := range chaoticNames {
			chaoticSet[p] = true
		}
		var filtered []string
		for _, p := range strings.Fields(string(native)) {
			if !aurSet[p] && !chaoticSet[p] {
				filtered = append(filtered, p)
			}
		}
		native = []byte(strings.Join(filtered, "\n"))
		if len(filtered) > 0 {
			native = append(native, '\n')
		}
	} else {
		native, err = run("apt-mark", "showmanual")
		if err != nil {
			return fmt.Errorf("apt-mark showmanual: %w", err)
		}
	}
	if _, e := exec.LookPath("flatpak"); e == nil {
		flatpak, _ = run("flatpak", "list", "--app", "--columns=application")
	}
	files := map[string][]byte{"lista-nativo.txt": native, "lista-aur.txt": aur, "lista-chaotic-aur.txt": chaotic, "lista-flatpak.txt": flatpak}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(s.TargetDir, name), data, 0644); err != nil {
			return err
		}
	}
	s.logf("\n[Gestor] Listas recopiladas: nativos=%d, AUR=%d, Chaotic-AUR=%d, Flatpak=%d\n", len(strings.Fields(string(native))), len(strings.Fields(string(aur))), len(strings.Fields(string(chaotic))), len(strings.Fields(string(flatpak))))
	return nil
}

func readPackageFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var result []string
	for _, line := range strings.Split(string(data), "\n") {
		p := strings.TrimSpace(line)
		if p == "" || seen[p] {
			continue
		}
		if !regexpPackage.MatchString(p) {
			continue
		}
		seen[p] = true
		result = append(result, p)
	}
	return result, nil
}

var regexpPackage = regexp.MustCompile(`^[A-Za-z0-9@._+:-]+$`)

func selectionFilePath(path string) string { return path + ".seleccion" }

func loadSavedSelection(path string, packages []string) []string {
	selected, err := readPackageFile(selectionFilePath(path))
	if err != nil {
		return nil
	}
	allowed := make(map[string]bool, len(packages))
	for _, p := range packages {
		allowed[p] = true
	}
	var result []string
	for _, p := range selected {
		if allowed[p] {
			result = append(result, p)
		}
	}
	return result
}

func savePackageSelection(path string, selected []string) error {
	data := strings.Join(selected, "\n")
	if len(selected) > 0 {
		data += "\n"
	}
	return os.WriteFile(selectionFilePath(path), []byte(data), 0644)
}

func packageSelectionUI(lang, path string, packages []string, instruction string) (*widget.CheckGroup, fyne.CanvasObject) {
	checks := widget.NewCheckGroup(packages, nil)
	checks.Horizontal = false
	checks.SetSelected(loadSavedSelection(path, packages))
	selectAll := widget.NewButton(tr(lang, "Seleccionar todos"), func() { checks.SetSelected(append([]string(nil), packages...)) })
	selectNone := widget.NewButton(tr(lang, "Deseleccionar todos"), func() { checks.SetSelected(nil) })
	saveLabel := widget.NewLabel(tr(lang, "La selección guardada se conservará para la próxima sesión."))
	save := widget.NewButton(tr(lang, "Guardar selección"), func() {
		if err := savePackageSelection(path, checks.Selected); err != nil {
			saveLabel.SetText(tr(lang, "Error al guardar: ") + err.Error())
			return
		}
		saveLabel.SetText(fmt.Sprintf(tr(lang, "Selección guardada: %d elementos."), len(checks.Selected)))
	})
	controls := container.NewHBox(selectAll, selectNone, save)
	content := container.NewBorder(container.NewVBox(widget.NewLabel(instruction), controls), saveLabel, nil, nil, container.NewScroll(checks))
	return checks, content
}

func (s *AppState) chooseRemovalList(fileName, title string, flatpak bool) {
	path := filepath.Join(s.TargetDir, fileName)
	packages, err := readPackageFile(path)
	if err != nil {
		dialog.ShowInformation(title, tr(s.Language, "No existe la lista. Primero utilizá 'Recopilar listas'."), s.Window)
		return
	}
	if len(packages) == 0 {
		dialog.ShowInformation(title, tr(s.Language, "La lista está vacía."), s.Window)
		return
	}
	checks, content := packageSelectionUI(s.Language, path, packages, tr(s.Language, "Seleccioná los elementos que querés quitar:"))
	listDialog := dialog.NewCustomConfirm(title, tr(s.Language, "Desinstalar seleccionados"), tr(s.Language, "Cerrar"), content, func(ok bool) {
		if !ok || len(checks.Selected) == 0 {
			return
		}
		selected := append([]string(nil), checks.Selected...)
		s.showConfirmNoDefault(tr(s.Language, "Confirmar desinstalación"), fmt.Sprintf(tr(s.Language, "Se quitarán %d elementos. ¿Continuar?"), len(selected)), func(confirmed bool) {
			if !confirmed {
				return
			}
			go func() {
				var err error
				if flatpak {
					args := append([]string{"uninstall", "-y"}, selected...)
					err = s.runUserLogged(context.Background(), "Desinstalar Flatpaks", "flatpak", args...)
				} else if s.Distro == DistroArch {
					args := append([]string{"-Rns", "--noconfirm"}, selected...)
					err = s.runLogged(context.Background(), "Desinstalar paquetes", "pacman", args...)
				} else {
					args := append([]string{"purge", "-y"}, selected...)
					err = s.runLogged(context.Background(), "Desinstalar paquetes", "apt-get", args...)
				}
				if err != nil {
					fyne.Do(func() { dialog.ShowError(err, s.Window) })
				}
			}()
		})
	}, s.Window)
	listDialog.Resize(fyne.NewSize(800, 620))
	listDialog.Show()
}

func (s *AppState) chooseInstallList(fileName, kind string) {
	path := filepath.Join(s.TargetDir, fileName)
	packages, err := readPackageFile(path)
	if err != nil {
		dialog.ShowInformation(tr(s.Language, "Instalar paquetes"), tr(s.Language, "No existe ")+path+tr(s.Language, ". Primero recopilá las listas."), s.Window)
		return
	}
	if len(packages) == 0 {
		dialog.ShowInformation(tr(s.Language, "Instalar paquetes"), tr(s.Language, "La lista está vacía."), s.Window)
		return
	}
	checks, content := packageSelectionUI(s.Language, path, packages, tr(s.Language, "Seleccioná los paquetes que querés instalar/restaurar:"))
	listDialog := dialog.NewCustomConfirm(tr(s.Language, "Instalar paquetes"), tr(s.Language, "Instalar seleccionados"), tr(s.Language, "Cerrar"), content, func(ok bool) {
		if !ok || len(checks.Selected) == 0 {
			return
		}
		selected := append([]string(nil), checks.Selected...)
		s.showConfirmNoDefault(tr(s.Language, "Confirmar instalación"), fmt.Sprintf(tr(s.Language, "Se instalarán %d paquetes de la lista %s. ¿Continuar?"), len(selected), kind), func(confirmed bool) {
			if !confirmed {
				return
			}
			go func() {
				var err error
				switch kind {
				case "aur":
					if s.Distro != DistroArch {
						err = fmt.Errorf("%s", tr(s.Language, "AUR sólo está disponible en Arch"))
					} else if _, e := exec.LookPath("yay"); e != nil {
						err = fmt.Errorf("no se encontró yay")
					} else {
						err = s.runUserLogged(context.Background(), "Instalar paquetes AUR", "yay", append([]string{"-S", "--needed", "--noconfirm"}, selected...)...)
					}
				case "chaotic":
					if s.Distro != DistroArch {
						err = fmt.Errorf("%s", tr(s.Language, "Chaotic-AUR sólo está disponible en Arch"))
					} else {
						repoArgs := append([]string{"-S", "--needed", "--noconfirm"}, func() []string {
							var a []string
							for _, p := range selected {
								a = append(a, "chaotic-aur/"+p)
							}
							return a
						}()...)
						err = s.runLogged(context.Background(), "Instalar paquetes Chaotic-AUR", "pacman", repoArgs...)
					}
				case "flatpak":
					if _, e := exec.LookPath("flatpak"); e != nil {
						err = fmt.Errorf("%s", tr(s.Language, "Flatpak no está instalado"))
					} else {
						err = s.runUserLogged(context.Background(), "Instalar Flatpaks", "flatpak", append([]string{"install", "--user", "--noninteractive", "-y"}, selected...)...)
					}
				default:
					if s.Distro == DistroArch {
						err = s.runLogged(context.Background(), tr(s.Language, "Instalar paquetes"), "pacman", append([]string{"-S", "--needed", "--noconfirm"}, selected...)...)
					} else {
						err = s.runLogged(context.Background(), tr(s.Language, "Instalar paquetes"), "apt-get", append([]string{"install", "-y"}, selected...)...)
					}
				}
				if err != nil {
					fyne.Do(func() { dialog.ShowError(err, s.Window) })
				}
			}()
		})
	}, s.Window)
	listDialog.Resize(fyne.NewSize(800, 620))
	listDialog.Show()
}

// taskButton dibuja su propio texto para poder resaltarlo mientras la acción está
// activa. Sus colores se consultan al tema, de modo que siguen la apariencia elegida.
type taskButton struct {
	widget.BaseWidget
	label  *canvas.Text
	active bool
	busy   bool
	action func()
}

func newTaskButton(text string, action func()) *taskButton {
	b := &taskButton{label: canvas.NewText(text, themeColor(colorNameButtonText)), action: action}
	b.label.Alignment = fyne.TextAlignCenter
	b.label.TextStyle = fyne.TextStyle{Bold: true}
	b.ExtendBaseWidget(b)
	return b
}

func (b *taskButton) background() color.Color {
	if b.active {
		return themeColor(theme.ColorNamePrimary)
	}
	return themeColor(colorNameButtonBackground)
}

func (b *taskButton) textColor() color.Color {
	if b.active {
		return color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	}
	return themeColor(colorNameButtonText)
}

func (b *taskButton) CreateRenderer() fyne.WidgetRenderer {
	c := b.background()
	b.label.Color = b.textColor()
	return &taskButtonRenderer{button: b, left: canvas.NewCircle(c), right: canvas.NewCircle(c), middle: canvas.NewRectangle(c)}
}
func (b *taskButton) MinSize() fyne.Size { return fyne.NewSize(180, 42) }

type taskButtonRenderer struct {
	button      *taskButton
	left, right *canvas.Circle
	middle      *canvas.Rectangle
}

func (r *taskButtonRenderer) Layout(size fyne.Size) {
	d := size.Height
	r.left.Resize(fyne.NewSize(d, d))
	r.left.Move(fyne.NewPos(0, 0))
	r.right.Resize(fyne.NewSize(d, d))
	r.right.Move(fyne.NewPos(size.Width-d, 0))
	r.middle.Resize(fyne.NewSize(size.Width-d, d))
	r.middle.Move(fyne.NewPos(d/2, 0))
	r.button.label.Resize(size)
	r.button.label.Move(fyne.NewPos(0, 0))
}
func (r *taskButtonRenderer) MinSize() fyne.Size { return r.button.MinSize() }
func (r *taskButtonRenderer) Refresh() {
	c := r.button.background()
	r.left.FillColor, r.right.FillColor, r.middle.FillColor = c, c, c
	r.button.label.Color = r.button.textColor()
	r.left.Refresh()
	r.right.Refresh()
	r.middle.Refresh()
	r.button.label.Refresh()
	r.Layout(r.button.Size())
}
func (r *taskButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.middle, r.left, r.right, r.button.label}
}
func (r *taskButtonRenderer) Destroy() {}
func (b *taskButton) Tapped(*fyne.PointEvent) {
	b.setActive(true)
	if b.action != nil {
		b.action()
	}
	time.AfterFunc(700*time.Millisecond, func() {
		fyne.Do(func() {
			if !b.busy {
				b.setActive(false)
			}
		})
	})
}
func (b *taskButton) setActive(active bool) {
	b.active = active
	b.Refresh()
}
