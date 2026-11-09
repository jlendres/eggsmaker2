package main

import "testing"

func TestTranslationsAvailable(t *testing.T) {
	cases := []struct{ lang, source, want string }{
		{"es", "Salir", "Salir"},
		{"en", "Salir", "Exit"},
		{"pt", "Salir", "Sair"},
		{"it", "Salir", "Esci"},
		{"en", "Abrir gestor de particiones", "Open partition manager"},
		{"pt", "No hay gestor de particiones instalado.", "Não há gerenciador de partições instalado."},
		{"it", "Gestor de particiones", "Gestore partizioni"},
	}
	for _, tc := range cases {
		if got := tr(tc.lang, tc.source); got != tc.want {
			t.Errorf("tr(%q, %q) = %q; want %q", tc.lang, tc.source, got, tc.want)
		}
	}
}

func TestValidLanguage(t *testing.T) {
	for _, lang := range []string{"es", "en", "pt", "it"} {
		if !validLanguage(lang) {
			t.Errorf("language %q should be valid", lang)
		}
	}
	if validLanguage("fr") {
		t.Error("French is not configured")
	}
}
