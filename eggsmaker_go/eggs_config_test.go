package main

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestEggsConfigYAMLRoundTrip(t *testing.T) {
	config := defaultEggsConfig()
	config.Remaster.User = `custom "user"`
	config.Remaster.Password = `p@ss: "word"`
	config.Remaster.Compression.Algorithm = "xz"
	config.Remaster.Compression.Level = 6
	config.Remaster.ISOPrefix = "my-linux"
	config.Remaster.Installer = "krill"
	ramMode := false
	config.Remaster.RamMode = &ramMode

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	var decoded eggsConfigData
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Remaster.User != config.Remaster.User ||
		decoded.Remaster.Password != config.Remaster.Password ||
		decoded.Remaster.Compression.Algorithm != config.Remaster.Compression.Algorithm ||
		decoded.Remaster.Compression.Level != config.Remaster.Compression.Level ||
		decoded.Remaster.ISOPrefix != config.Remaster.ISOPrefix ||
		decoded.Remaster.Installer != config.Remaster.Installer ||
		decoded.Remaster.RamMode == nil || *decoded.Remaster.RamMode {
		t.Fatalf("configuration YAML round trip changed values: %+v", decoded.Remaster)
	}
}

func TestEggsConfigTranslations(t *testing.T) {
	for _, lang := range []string{"en", "pt", "it"} {
		for _, source := range []string{"Configuración", "Guardar", "Cancelar", "Sí", "Usuario Live:", "Modo RAM:", "Exclusiones para la ISO"} {
			if _, ok := eggsConfigTranslations[lang][source]; !ok {
				t.Errorf("missing %s translation for %q", lang, source)
			}
		}
		if got := tr("pt", "No"); got != "Não" {
			t.Errorf("Portuguese No translation = %q", got)
		}
	}
}
