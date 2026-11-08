package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/yaml.v3"
)

const (
	eggsCustomConfigPath  = "/etc/penguins-eggs.d/custom.yaml"
	eggsCustomExcludePath = "/etc/penguins-eggs.d/custom.exclude.list"
)

type eggsConfigData struct {
	Remaster struct {
		User        string `yaml:"user,omitempty"`
		Password    string `yaml:"password"`
		Compression struct {
			Algorithm string `yaml:"algorithm"`
			Level     int    `yaml:"level,omitempty"`
		} `yaml:"compression"`
		ISOPrefix string `yaml:"iso_prefix,omitempty"`
		Installer string `yaml:"installer,omitempty"`
		RamMode   *bool  `yaml:"ram_mode,omitempty"`
	} `yaml:"remaster"`
}

func defaultEggsConfig() eggsConfigData {
	var config eggsConfigData
	config.Remaster.User = "live"
	config.Remaster.Password = "evolution"
	config.Remaster.Compression.Algorithm = "zstd"
	config.Remaster.Compression.Level = 3
	config.Remaster.Installer = "krill"
	if _, err := exec.LookPath("calamares"); err == nil {
		config.Remaster.Installer = "calamares"
	}
	ramMode := true
	config.Remaster.RamMode = &ramMode
	return config
}

func loadEggsConfig() (eggsConfigData, string, error) {
	config := defaultEggsConfig()
	data, err := os.ReadFile(eggsCustomConfigPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return eggsConfigData{}, "", fmt.Errorf("no se pudo leer %s: %w", eggsCustomConfigPath, err)
	}
	if err == nil {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return eggsConfigData{}, "", fmt.Errorf("configuración YAML inválida en %s: %w", eggsCustomConfigPath, err)
		}
	}

	excludes, err := os.ReadFile(eggsCustomExcludePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return eggsConfigData{}, "", fmt.Errorf("no se pudo leer %s: %w", eggsCustomExcludePath, err)
	}
	return config, string(excludes), nil
}

func (s *AppState) saveEggsConfigFile(path string, contents []byte) error {
	cmd := s.privilegedCommand(context.Background(), "tee", path)
	cmd.Stdin = bytes.NewReader(contents)
	var stderr bytes.Buffer
	cmd.Stdout = io.Discard
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail != "" {
			return fmt.Errorf("%s: %w: %s", path, err, detail)
		}
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

func (s *AppState) showEggsConfigDialog(onClose func()) {
	config, excludes, err := loadEggsConfig()
	if err != nil {
		dialog.ShowError(err, s.Window)
		return
	}
	lang := s.Language

	user := widget.NewEntry()
	user.SetText(config.Remaster.User)
	password := widget.NewPasswordEntry()
	password.SetText(config.Remaster.Password)
	algorithm := widget.NewSelect([]string{"zstd", "xz", "lz4", "gzip"}, nil)
	algorithm.SetSelected(config.Remaster.Compression.Algorithm)
	level := widget.NewEntry()
	level.SetText(strconv.Itoa(config.Remaster.Compression.Level))
	level.SetPlaceHolder("3")
	prefix := widget.NewEntry()
	prefix.SetText(config.Remaster.ISOPrefix)
	prefix.SetPlaceHolder(tr(lang, "Automático"))
	ramMode := widget.NewCheck(tr(lang, "Activar modo RAM"), nil)
	ramMode.SetChecked(config.Remaster.RamMode == nil || *config.Remaster.RamMode)
	installer := widget.NewSelect([]string{"krill", "calamares"}, nil)
	installer.SetSelected(config.Remaster.Installer)

	levelItem := widget.NewFormItem(tr(lang, "Nivel de compresión:"), level)
	formItems := func() []*widget.FormItem {
		items := []*widget.FormItem{
			widget.NewFormItem(tr(lang, "Usuario Live:"), user),
			widget.NewFormItem(tr(lang, "Contraseña Live:"), password),
			widget.NewFormItem(tr(lang, "Algoritmo de compresión:"), algorithm),
			levelItem,
			widget.NewFormItem(tr(lang, "Prefijo de ISO:"), prefix),
			widget.NewFormItem(tr(lang, "Modo RAM:"), ramMode),
			widget.NewFormItem(tr(lang, "Instalador:"), installer),
		}
		if algorithm.Selected != "zstd" {
			items = append(items[:3], items[4:]...)
		}
		return items
	}
	form := widget.NewForm(formItems()...)
	algorithm.OnChanged = func(value string) {
		form.Items = formItems()
		form.Refresh()
	}
	settingsTab := container.NewTabItem(tr(lang, "Configuración"), form)

	excludeList := widget.NewMultiLineEntry()
	excludeList.SetText(excludes)
	excludeList.Wrapping = fyne.TextWrapOff
	excludeList.SetMinRowsVisible(12)
	excludesTab := container.NewTabItem(tr(lang, "Exclusiones para la ISO"), container.NewBorder(
		widget.NewLabel(tr(lang, "Una ruta por línea; las líneas que comienzan con # son comentarios.")),
		nil, nil, nil, excludeList,
	))
	tabs := container.NewAppTabs(settingsTab, excludesTab)
	tabs.SetTabLocation(container.TabLocationTop)

	var d *dialog.CustomDialog
	save := func() {
		if user.Text == "" {
			dialog.ShowError(fmt.Errorf("%s", tr(lang, "El usuario Live no puede quedar vacío.")), s.Window)
			return
		}
		if password.Text == "" {
			dialog.ShowError(fmt.Errorf("%s", tr(lang, "La contraseña Live no puede quedar vacía.")), s.Window)
			return
		}
		compressionLevel, err := strconv.Atoi(strings.TrimSpace(level.Text))
		if err != nil || compressionLevel < 1 {
			dialog.ShowError(fmt.Errorf("%s", tr(lang, "El nivel de compresión debe ser un entero positivo.")), s.Window)
			return
		}
		if !slices.Contains([]string{"zstd", "xz", "lz4", "gzip"}, algorithm.Selected) ||
			!slices.Contains([]string{"krill", "calamares"}, installer.Selected) {
			dialog.ShowError(fmt.Errorf("%s", tr(lang, "Seleccioná un algoritmo y un instalador.")), s.Window)
			return
		}

		config.Remaster.User = user.Text
		config.Remaster.Password = password.Text
		config.Remaster.Compression.Algorithm = algorithm.Selected
		config.Remaster.Compression.Level = compressionLevel
		config.Remaster.ISOPrefix = strings.TrimSpace(prefix.Text)
		config.Remaster.Installer = installer.Selected
		ramModeValue := ramMode.Checked
		config.Remaster.RamMode = &ramModeValue
		yamlData, err := yaml.Marshal(config)
		if err != nil {
			dialog.ShowError(err, s.Window)
			return
		}
		excludeData := []byte(excludeList.Text)
		d.Hide()
		progress := dialog.NewProgressInfinite(tr(lang, "Guardando configuración"), tr(lang, "Guardando archivos de Penguins Eggs…"), s.Window)
		progress.Show()
		go func() {
			err := s.saveEggsConfigFile(eggsCustomConfigPath, yamlData)
			if err == nil {
				if saveErr := s.saveEggsConfigFile(eggsCustomExcludePath, excludeData); saveErr != nil {
					err = fmt.Errorf("%s: %w", tr(lang, "Se guardó la configuración, pero falló la lista de exclusiones"), saveErr)
				}
			}
			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					d.Show()
					dialog.ShowError(err, s.Window)
					return
				}
				if onClose != nil {
					onClose()
				} else {
					dialog.ShowInformation(tr(lang, "Configuración guardada"), tr(lang, "La configuración de Penguins Eggs se guardó correctamente."), s.Window)
				}
			})
		}()
	}
	cancel := widget.NewButton(tr(lang, "Cancelar"), func() {
		d.Hide()
		if onClose != nil {
			onClose()
		}
	})
	cancel.Importance = widget.LowImportance
	saveButton := widget.NewButton(tr(lang, "Guardar"), save)
	saveButton.Importance = widget.HighImportance
	d = dialog.NewCustomWithoutButtons(tr(lang, "Configuración de Penguins Eggs"), tabs, s.Window)
	d.SetButtons([]fyne.CanvasObject{cancel, saveButton})
	d.Resize(fyne.NewSize(640, 520))
	d.Show()
}

var eggsConfigTranslations = map[string]map[string]string{
	"en": {
		"Automático": "Automatic", "Activar modo RAM": "Enable RAM mode", "Sí": "Yes", "No": "No",
		"Nivel de compresión:": "Compression level:", "Usuario Live:": "Live user:", "Contraseña Live:": "Live password:", "Modo RAM:": "RAM mode:",
		"Algoritmo de compresión:": "Compression algorithm:", "Prefijo de ISO:": "ISO prefix:", "Instalador:": "Installer:",
		"Configuración": "Settings", "Exclusiones para la ISO": "ISO exclusions",
		"Una ruta por línea; las líneas que comienzan con # son comentarios.": "One path per line; lines beginning with # are comments.",
		"Guardar": "Save", "Cancelar": "Cancel", "El usuario Live no puede quedar vacío.": "The live user cannot be empty.",
		"La contraseña Live no puede quedar vacía.":                      "The live password cannot be empty.",
		"El nivel de compresión debe ser un entero positivo.":            "The compression level must be a positive integer.",
		"Seleccioná un algoritmo y un instalador.":                       "Select an algorithm and an installer.",
		"Se guardó la configuración, pero falló la lista de exclusiones": "The configuration was saved, but saving the exclusions list failed",
		"Configuración guardada":                                         "Configuration saved", "La configuración de Penguins Eggs se guardó correctamente.": "Penguins Eggs configuration was saved successfully.",
		"Guardando configuración": "Saving configuration", "Guardando archivos de Penguins Eggs…": "Saving Penguins Eggs files…",
	},
	"pt": {
		"Automático": "Automático", "Activar modo RAM": "Ativar modo RAM", "Sí": "Sim", "No": "Não",
		"Nivel de compresión:": "Nível de compressão:", "Usuario Live:": "Usuário Live:", "Contraseña Live:": "Senha Live:", "Modo RAM:": "Modo RAM:",
		"Algoritmo de compresión:": "Algoritmo de compressão:", "Prefijo de ISO:": "Prefixo da ISO:", "Instalador:": "Instalador:",
		"Configuración": "Configuração", "Exclusiones para la ISO": "Exclusões da ISO",
		"Una ruta por línea; las líneas que comienzan con # son comentarios.": "Um caminho por linha; linhas que começam com # são comentários.",
		"Guardar": "Salvar", "Cancelar": "Cancelar", "El usuario Live no puede quedar vacío.": "O usuário Live não pode ficar vazio.",
		"La contraseña Live no puede quedar vacía.":                      "A senha Live não pode ficar vazia.",
		"El nivel de compresión debe ser un entero positivo.":            "O nível de compressão deve ser um inteiro positivo.",
		"Seleccioná un algoritmo y un instalador.":                       "Selecione um algoritmo e um instalador.",
		"Se guardó la configuración, pero falló la lista de exclusiones": "A configuração foi salva, mas falhou ao salvar a lista de exclusões",
		"Configuración guardada":                                         "Configuração salva", "La configuración de Penguins Eggs se guardó correctamente.": "A configuração do Penguins Eggs foi salva.",
		"Guardando configuración": "Salvando configuração", "Guardando archivos de Penguins Eggs…": "Salvando arquivos do Penguins Eggs…",
	},
	"it": {
		"Automático": "Automatico", "Activar modo RAM": "Attiva modalità RAM", "Sí": "Sì", "No": "No",
		"Nivel de compresión:": "Livello di compressione:", "Usuario Live:": "Utente Live:", "Contraseña Live:": "Password Live:", "Modo RAM:": "Modalità RAM:",
		"Algoritmo de compresión:": "Algoritmo di compressione:", "Prefijo de ISO:": "Prefisso ISO:", "Instalador:": "Programma di installazione:",
		"Configuración": "Configurazione", "Exclusiones para la ISO": "Esclusioni ISO",
		"Una ruta por línea; las líneas que comienzan con # son comentarios.": "Un percorso per riga; le righe che iniziano con # sono commenti.",
		"Guardar": "Salva", "Cancelar": "Annulla", "El usuario Live no puede quedar vacío.": "L'utente Live non può essere vuoto.",
		"La contraseña Live no puede quedar vacía.":                      "La password Live non può essere vuota.",
		"El nivel de compresión debe ser un entero positivo.":            "Il livello di compressione deve essere un numero intero positivo.",
		"Seleccioná un algoritmo y un instalador.":                       "Seleziona un algoritmo e un programma di installazione.",
		"Se guardó la configuración, pero falló la lista de exclusiones": "La configurazione è stata salvata, ma non è stato possibile salvare l'elenco delle esclusioni",
		"Configuración guardada":                                         "Configurazione salvata", "La configuración de Penguins Eggs se guardó correctamente.": "La configurazione di Penguins Eggs è stata salvata.",
		"Guardando configuración": "Salvataggio configurazione", "Guardando archivos de Penguins Eggs…": "Salvataggio dei file di Penguins Eggs…",
	},
}
