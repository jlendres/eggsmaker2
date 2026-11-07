package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Modos de tema disponibles.
//   - default: apariencia original de Eggsmaker (widgets oscuros sobre cualquier tema).
//   - system:  los colores siguen el tema claro/oscuro del sistema.
//   - dark / light: tema fijo, independiente del sistema.
const (
	themeModeDefault = "default"
	themeModeSystem  = "system"
	themeModeDark    = "dark"
	themeModeLight   = "light"

	accentCustom     = "custom"
	defaultAccentKey = "orange"

	prefThemeMode    = "theme_mode"
	prefAccent       = "accent"
	prefAccentCustom = "accent_custom"
)

// Nombres de color propios del tema, usados por los widgets personalizados
// (botones de tareas, barra de progreso) para seguir el tema activo.
const (
	colorNameButtonBackground fyne.ThemeColorName = "eggsmakerButtonBackground"
	colorNameButtonText       fyne.ThemeColorName = "eggsmakerButtonText"
	colorNameAccentText       fyne.ThemeColorName = "eggsmakerAccentText"
)

var themeModeKeys = []string{themeModeDefault, themeModeSystem, themeModeDark, themeModeLight}

var themeModeLabels = map[string]string{
	themeModeDefault: "Predeterminado (actual)",
	themeModeSystem:  "Seguir al sistema",
	themeModeDark:    "Oscuro",
	themeModeLight:   "Claro",
}

type accentPreset struct {
	key     string
	label   string
	primary color.NRGBA
	text    color.NRGBA
}

var accentPresets = []accentPreset{
	{"orange", "Naranja (predeterminado)", color.NRGBA{R: 255, G: 105, B: 0, A: 255}, color.NRGBA{R: 255, G: 125, B: 0, A: 255}},
	{"blue", "Azul", color.NRGBA{R: 33, G: 150, B: 243, A: 255}, color.NRGBA{R: 66, G: 165, B: 245, A: 255}},
	{"green", "Verde", color.NRGBA{R: 46, G: 160, B: 67, A: 255}, color.NRGBA{R: 63, G: 185, B: 80, A: 255}},
	{"purple", "Violeta", color.NRGBA{R: 142, G: 68, B: 230, A: 255}, color.NRGBA{R: 163, G: 113, B: 247, A: 255}},
	{"red", "Rojo", color.NRGBA{R: 229, G: 57, B: 53, A: 255}, color.NRGBA{R: 244, G: 81, B: 81, A: 255}},
	{"cyan", "Cian", color.NRGBA{R: 0, G: 172, B: 193, A: 255}, color.NRGBA{R: 38, G: 198, B: 218, A: 255}},
}

func findAccent(key string) (accentPreset, bool) {
	for _, p := range accentPresets {
		if p.key == key {
			return p, true
		}
	}
	return accentPreset{}, false
}

// appearanceSettings guarda la configuración de colores elegida por la persona.
type appearanceSettings struct {
	Mode   string // default | system | dark | light
	Accent string // clave de accentPresets o "custom"
	Custom string // color personalizado en formato #RRGGBB (vacío si nunca se eligió)
}

func defaultAppearance() appearanceSettings {
	return appearanceSettings{Mode: themeModeDefault, Accent: defaultAccentKey}
}

func validThemeMode(mode string) bool {
	for _, k := range themeModeKeys {
		if k == mode {
			return true
		}
	}
	return false
}

// normalizeAppearance descarta valores desconocidos o corruptos y vuelve a los predeterminados.
func normalizeAppearance(a appearanceSettings) appearanceSettings {
	if !validThemeMode(a.Mode) {
		a.Mode = themeModeDefault
	}
	if _, ok := parseHexColor(a.Custom); !ok {
		a.Custom = ""
	}
	if a.Accent == accentCustom {
		if a.Custom == "" {
			a.Accent = defaultAccentKey
		}
	} else if _, ok := findAccent(a.Accent); !ok {
		a.Accent = defaultAccentKey
	}
	return a
}

func loadAppearance(p fyne.Preferences) appearanceSettings {
	return normalizeAppearance(appearanceSettings{
		Mode:   p.StringWithFallback(prefThemeMode, themeModeDefault),
		Accent: p.StringWithFallback(prefAccent, defaultAccentKey),
		Custom: p.StringWithFallback(prefAccentCustom, ""),
	})
}

func (a appearanceSettings) save(p fyne.Preferences) {
	p.SetString(prefThemeMode, a.Mode)
	p.SetString(prefAccent, a.Accent)
	p.SetString(prefAccentCustom, a.Custom)
}

// accentColors devuelve el color primario (botones activos, barra de progreso)
// y el color de texto de acento (títulos, consola, marco).
func (a appearanceSettings) accentColors() (primary, text color.NRGBA) {
	if a.Accent == accentCustom {
		if c, ok := parseHexColor(a.Custom); ok {
			return c, c
		}
	}
	preset, ok := findAccent(a.Accent)
	if !ok {
		preset = accentPresets[0]
	}
	return preset.primary, preset.text
}

func parseHexColor(s string) (color.NRGBA, bool) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(s) != 6 {
		return color.NRGBA{}, false
	}
	var ch [3]uint8
	for i := 0; i < 3; i++ {
		v, err := strconv.ParseUint(s[i*2:i*2+2], 16, 8)
		if err != nil {
			return color.NRGBA{}, false
		}
		ch[i] = uint8(v)
	}
	return color.NRGBA{R: ch[0], G: ch[1], B: ch[2], A: 255}, true
}

func formatHexColor(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}

// eggsmakerTheme aplica el acento elegido y la política de tema (predeterminado,
// sistema, oscuro o claro) sobre el tema base de Fyne.
type eggsmakerTheme struct {
	fyne.Theme
	appearance appearanceSettings
}

func newEggsmakerTheme(a appearanceSettings) eggsmakerTheme {
	return eggsmakerTheme{Theme: theme.DefaultTheme(), appearance: normalizeAppearance(a)}
}

// effectiveVariant devuelve la variante realmente usada: oscuro/claro la fuerzan;
// predeterminado y sistema respetan la variante que informa el sistema.
func (t eggsmakerTheme) effectiveVariant(v fyne.ThemeVariant) fyne.ThemeVariant {
	switch t.appearance.Mode {
	case themeModeDark:
		return theme.VariantDark
	case themeModeLight:
		return theme.VariantLight
	}
	return v
}

func (t eggsmakerTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	v := t.effectiveVariant(variant)
	primary, text := t.appearance.accentColors()
	// En modo predeterminado los widgets propios son siempre oscuros (apariencia
	// original); en los demás modos siguen la variante efectiva.
	darkWidgets := t.appearance.Mode == themeModeDefault || v == theme.VariantDark
	if !darkWidgets {
		// Sobre fondos claros el acento primario ofrece mejor contraste.
		text = primary
	}
	switch name {
	case theme.ColorNameDisabled:
		return text
	case theme.ColorNamePrimary:
		return primary
	case theme.ColorNameInputBackground:
		if darkWidgets {
			return color.NRGBA{R: 42, G: 42, B: 48, A: 255}
		}
	case colorNameAccentText:
		return text
	case colorNameButtonBackground:
		if darkWidgets {
			return color.NRGBA{R: 55, G: 55, B: 62, A: 255}
		}
		return color.NRGBA{R: 222, G: 222, B: 228, A: 255}
	case colorNameButtonText:
		if darkWidgets {
			return color.NRGBA{R: 245, G: 245, B: 245, A: 255}
		}
		return color.NRGBA{R: 28, G: 28, B: 32, A: 255}
	}
	return t.Theme.Color(name, v)
}

// themeColor consulta un color del tema activo (incluidos los nombres propios).
func themeColor(name fyne.ThemeColorName) color.Color {
	return theme.Color(name)
}

// bindColor aplica un color del tema ahora y cada vez que cambie la apariencia.
func (s *AppState) bindColor(name fyne.ThemeColorName, apply func(color.Color)) {
	apply(themeColor(name))
	s.themeRefresh = append(s.themeRefresh, func() { apply(themeColor(name)) })
}

// applyAppearance guarda la configuración y la aplica a toda la aplicación.
func (s *AppState) applyAppearance(myApp fyne.App) {
	s.Appearance = normalizeAppearance(s.Appearance)
	s.Appearance.save(myApp.Preferences())
	myApp.Settings().SetTheme(newEggsmakerTheme(s.Appearance))
	for _, refresh := range s.themeRefresh {
		refresh()
	}
	if s.Window != nil && s.Window.Content() != nil {
		s.Window.Content().Refresh()
	}
}

func keyForLabel(keys []string, labels map[string]string, lang, label string) string {
	for _, k := range keys {
		if tr(lang, labels[k]) == label {
			return k
		}
	}
	return keys[0]
}

func (s *AppState) showAppearanceDialog(myApp fyne.App) {
	lang := s.Language

	accentKeys := make([]string, 0, len(accentPresets)+1)
	accentLabels := map[string]string{accentCustom: "Personalizado"}
	for _, p := range accentPresets {
		accentKeys = append(accentKeys, p.key)
		accentLabels[p.key] = p.label
	}
	accentKeys = append(accentKeys, accentCustom)

	options := func(keys []string, labels map[string]string) []string {
		out := make([]string, 0, len(keys))
		for _, k := range keys {
			out = append(out, tr(lang, labels[k]))
		}
		return out
	}

	updating := false
	modeSelect := widget.NewSelect(options(themeModeKeys, themeModeLabels), nil)
	modeSelect.SetSelected(tr(lang, themeModeLabels[s.Appearance.Mode]))
	accentSelect := widget.NewSelect(options(accentKeys, accentLabels), nil)
	accentSelect.SetSelected(tr(lang, accentLabels[s.Appearance.Accent]))

	setAccentSelect := func(key string) {
		updating = true
		accentSelect.SetSelected(tr(lang, accentLabels[key]))
		updating = false
	}

	pickCustom := func() {
		picker := dialog.NewColorPicker(tr(lang, "Color personalizado"), tr(lang, "Elegí el color de acento de la aplicación."), func(c color.Color) {
			s.Appearance.Custom = formatHexColor(c)
			s.Appearance.Accent = accentCustom
			setAccentSelect(accentCustom)
			s.applyAppearance(myApp)
		}, s.Window)
		picker.Advanced = true
		picker.Show()
	}

	modeSelect.OnChanged = func(label string) {
		if updating {
			return
		}
		s.Appearance.Mode = keyForLabel(themeModeKeys, themeModeLabels, lang, label)
		s.applyAppearance(myApp)
	}
	accentSelect.OnChanged = func(label string) {
		if updating {
			return
		}
		key := keyForLabel(accentKeys, accentLabels, lang, label)
		if key == accentCustom && s.Appearance.Custom == "" {
			// Todavía no hay color personalizado: se vuelve a la opción anterior
			// y se abre el selector; si se cancela, nada cambia.
			setAccentSelect(s.Appearance.Accent)
			pickCustom()
			return
		}
		s.Appearance.Accent = key
		s.applyAppearance(myApp)
	}

	reset := widget.NewButton(tr(lang, "Restablecer valores predeterminados"), func() {
		s.Appearance = defaultAppearance()
		updating = true
		modeSelect.SetSelected(tr(lang, themeModeLabels[s.Appearance.Mode]))
		accentSelect.SetSelected(tr(lang, accentLabels[s.Appearance.Accent]))
		updating = false
		s.applyAppearance(myApp)
	})
	custom := widget.NewButton(tr(lang, "Elegir color personalizado…"), pickCustom)

	hint := widget.NewLabel(tr(lang, "«Predeterminado» conserva la apariencia original. «Seguir al sistema» adapta los colores al tema claro u oscuro del sistema."))
	hint.Wrapping = fyne.TextWrapWord

	form := widget.NewForm(
		widget.NewFormItem(tr(lang, "Tema:"), modeSelect),
		widget.NewFormItem(tr(lang, "Color de acento:"), accentSelect),
	)
	content := container.NewVBox(form, custom, reset, hint)
	d := dialog.NewCustom(tr(lang, "Apariencia"), tr(lang, "Cerrar"), content, s.Window)
	d.Resize(fyne.NewSize(480, d.MinSize().Height))
	d.Show()
}

var appearanceTranslations = map[string]map[string]string{
	"en": {
		"Apariencia": "Appearance", "Tema:": "Theme:", "Color de acento:": "Accent color:",
		"Predeterminado (actual)": "Default (current)", "Seguir al sistema": "Follow system", "Oscuro": "Dark", "Claro": "Light",
		"Naranja (predeterminado)": "Orange (default)", "Azul": "Blue", "Verde": "Green", "Violeta": "Purple", "Rojo": "Red", "Cian": "Cyan", "Personalizado": "Custom",
		"Elegir color personalizado…": "Choose custom color…", "Restablecer valores predeterminados": "Restore defaults",
		"Color personalizado": "Custom color", "Elegí el color de acento de la aplicación.": "Choose the application's accent color.",
		"«Predeterminado» conserva la apariencia original. «Seguir al sistema» adapta los colores al tema claro u oscuro del sistema.": "“Default” keeps the original look. “Follow system” adapts the colors to the system's light or dark theme.",
	},
	"pt": {
		"Apariencia": "Aparência", "Tema:": "Tema:", "Color de acento:": "Cor de destaque:",
		"Predeterminado (actual)": "Padrão (atual)", "Seguir al sistema": "Seguir o sistema", "Oscuro": "Escuro", "Claro": "Claro",
		"Naranja (predeterminado)": "Laranja (padrão)", "Azul": "Azul", "Verde": "Verde", "Violeta": "Roxo", "Rojo": "Vermelho", "Cian": "Ciano", "Personalizado": "Personalizado",
		"Elegir color personalizado…": "Escolher cor personalizada…", "Restablecer valores predeterminados": "Restaurar padrões",
		"Color personalizado": "Cor personalizada", "Elegí el color de acento de la aplicación.": "Escolha a cor de destaque do aplicativo.",
		"«Predeterminado» conserva la apariencia original. «Seguir al sistema» adapta los colores al tema claro u oscuro del sistema.": "«Padrão» mantém a aparência original. «Seguir o sistema» adapta as cores ao tema claro ou escuro do sistema.",
	},
	"it": {
		"Apariencia": "Aspetto", "Tema:": "Tema:", "Color de acento:": "Colore di accento:",
		"Predeterminado (actual)": "Predefinito (attuale)", "Seguir al sistema": "Segui il sistema", "Oscuro": "Scuro", "Claro": "Chiaro",
		"Naranja (predeterminado)": "Arancione (predefinito)", "Azul": "Blu", "Verde": "Verde", "Violeta": "Viola", "Rojo": "Rosso", "Cian": "Ciano", "Personalizado": "Personalizzato",
		"Elegir color personalizado…": "Scegli colore personalizzato…", "Restablecer valores predeterminados": "Ripristina predefiniti",
		"Color personalizado": "Colore personalizzato", "Elegí el color de acento de la aplicación.": "Scegli il colore di accento dell'applicazione.",
		"«Predeterminado» conserva la apariencia original. «Seguir al sistema» adapta los colores al tema claro u oscuro del sistema.": "«Predefinito» mantiene l'aspetto originale. «Segui il sistema» adatta i colori al tema chiaro o scuro del sistema.",
	},
}
