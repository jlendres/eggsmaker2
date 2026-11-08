package main

import (
	"image/color"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

// El tema base de Fyne consulta la app actual; las pruebas necesitan una app de prueba.
func TestMain(m *testing.M) {
	test.NewApp()
	os.Exit(m.Run())
}

func TestDefaultAppearanceKeepsOriginalColors(t *testing.T) {
	th := newEggsmakerTheme(defaultAppearance())
	want := map[fyne.ThemeColorName]color.NRGBA{
		theme.ColorNamePrimary:         {R: 255, G: 105, B: 0, A: 255},
		theme.ColorNameDisabled:        {R: 255, G: 125, B: 0, A: 255},
		theme.ColorNameInputBackground: {R: 42, G: 42, B: 48, A: 255},
		colorNameButtonBackground:      {R: 55, G: 55, B: 62, A: 255},
		colorNameButtonText:            {R: 245, G: 245, B: 245, A: 255},
		colorNameAccentText:            {R: 255, G: 125, B: 0, A: 255},
	}
	for _, variant := range []fyne.ThemeVariant{theme.VariantDark, theme.VariantLight} {
		for name, c := range want {
			if got := th.Color(name, variant); got != c {
				t.Errorf("default mode, variant %v, %s = %v; want %v", variant, name, got, c)
			}
		}
	}
}

func TestSystemModeFollowsVariant(t *testing.T) {
	th := newEggsmakerTheme(appearanceSettings{Mode: themeModeSystem, Accent: defaultAccentKey})
	dark := th.Color(colorNameButtonBackground, theme.VariantDark)
	light := th.Color(colorNameButtonBackground, theme.VariantLight)
	if dark == light {
		t.Fatal("system mode should use different button colors for dark and light variants")
	}
	if got := th.Color(colorNameButtonBackground, theme.VariantDark); got != (color.NRGBA{R: 55, G: 55, B: 62, A: 255}) {
		t.Errorf("dark variant button background = %v", got)
	}
	if got := th.Color(theme.ColorNameBackground, theme.VariantLight); got != theme.DefaultTheme().Color(theme.ColorNameBackground, theme.VariantLight) {
		t.Errorf("light variant background should come from the base theme, got %v", got)
	}
}

func TestForcedModesIgnoreSystemVariant(t *testing.T) {
	dark := newEggsmakerTheme(appearanceSettings{Mode: themeModeDark, Accent: defaultAccentKey})
	light := newEggsmakerTheme(appearanceSettings{Mode: themeModeLight, Accent: defaultAccentKey})
	for _, variant := range []fyne.ThemeVariant{theme.VariantDark, theme.VariantLight} {
		if got, want := dark.Color(theme.ColorNameBackground, variant), theme.DefaultTheme().Color(theme.ColorNameBackground, theme.VariantDark); got != want {
			t.Errorf("forced dark, system variant %v: background = %v; want %v", variant, got, want)
		}
		if got, want := light.Color(theme.ColorNameBackground, variant), theme.DefaultTheme().Color(theme.ColorNameBackground, theme.VariantLight); got != want {
			t.Errorf("forced light, system variant %v: background = %v; want %v", variant, got, want)
		}
	}
}

func TestAccentSelection(t *testing.T) {
	blue, _ := findAccent("blue")
	th := newEggsmakerTheme(appearanceSettings{Mode: themeModeDefault, Accent: "blue"})
	if got := th.Color(theme.ColorNamePrimary, theme.VariantDark); got != blue.primary {
		t.Errorf("blue primary = %v; want %v", got, blue.primary)
	}
	custom := newEggsmakerTheme(appearanceSettings{Mode: themeModeDefault, Accent: accentCustom, Custom: "#112233"})
	if got := custom.Color(theme.ColorNamePrimary, theme.VariantDark); got != (color.NRGBA{R: 0x11, G: 0x22, B: 0x33, A: 255}) {
		t.Errorf("custom primary = %v", got)
	}
}

func TestCustomAppearanceBackgroundColors(t *testing.T) {
	th := newEggsmakerTheme(appearanceSettings{
		Mode:              themeModeDefault,
		Accent:            defaultAccentKey,
		AppBackground:     "#123456",
		ConsoleBackground: "#654321",
	})
	if got := th.Color(theme.ColorNameBackground, theme.VariantDark); got != (color.NRGBA{R: 0x12, G: 0x34, B: 0x56, A: 255}) {
		t.Errorf("custom main background = %v", got)
	}
	if got := th.Color(colorNameConsoleBackground, theme.VariantDark); got != (color.NRGBA{R: 0x65, G: 0x43, B: 0x21, A: 255}) {
		t.Errorf("custom console background = %v", got)
	}
	if got := th.Size(theme.SizeNameInputRadius); got != 8 {
		t.Errorf("button/input corner radius = %v; want 8", got)
	}
}

func TestParseAndFormatHexColor(t *testing.T) {
	c, ok := parseHexColor("#FF6900")
	if !ok || c != (color.NRGBA{R: 255, G: 105, B: 0, A: 255}) {
		t.Fatalf("parseHexColor = %v, %v", c, ok)
	}
	if got := formatHexColor(c); got != "#ff6900" {
		t.Errorf("formatHexColor = %q", got)
	}
	for _, bad := range []string{"", "#fff", "zzzzzz", "#12345678"} {
		if _, ok := parseHexColor(bad); ok {
			t.Errorf("parseHexColor(%q) should fail", bad)
		}
	}
}

func TestNormalizeAppearance(t *testing.T) {
	got := normalizeAppearance(appearanceSettings{Mode: "rainbow", Accent: "nope", Custom: "xx"})
	if got != defaultAppearance() {
		t.Errorf("invalid values should fall back to defaults, got %+v", got)
	}
	got = normalizeAppearance(appearanceSettings{
		Mode:              themeModeDefault,
		Accent:            defaultAccentKey,
		AppBackground:     "not-a-color",
		ConsoleBackground: "also-invalid",
	})
	if got.AppBackground != "" || got.ConsoleBackground != defaultConsoleBackground {
		t.Errorf("invalid background colors should fall back to defaults, got %+v", got)
	}
	got = normalizeAppearance(appearanceSettings{Mode: themeModeSystem, Accent: accentCustom})
	if got.Accent != defaultAccentKey || got.Mode != themeModeSystem {
		t.Errorf("custom accent without color should fall back to orange, got %+v", got)
	}
	keep := appearanceSettings{
		Mode:              themeModeLight,
		Accent:            accentCustom,
		Custom:            "#abcdef",
		ConsoleBackground: defaultConsoleBackground,
	}
	if got := normalizeAppearance(keep); got != keep {
		t.Errorf("valid settings should be preserved, got %+v", got)
	}
}

func TestAppearanceTranslations(t *testing.T) {
	for _, lang := range []string{"en", "pt", "it"} {
		for _, src := range append([]string{"Apariencia", "Tema:", "Color de acento:", "Personalizado", "Color de fondo de la ventana principal…", "Color de fondo de la salida de procesos…"}, func() []string {
			var all []string
			for _, l := range themeModeLabels {
				all = append(all, l)
			}
			for _, p := range accentPresets {
				all = append(all, p.label)
			}
			return all
		}()...) {
			if tr(lang, src) == "" {
				t.Errorf("empty translation for %q (%s)", src, lang)
			}
		}
	}
	if tr("en", "Oscuro") != "Dark" {
		t.Error("Oscuro should translate to Dark")
	}
}
