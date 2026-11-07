#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$(readlink -f "$0")")"

if ! command -v go >/dev/null 2>&1; then
  echo "Error: Go no está instalado. Ejecutá ./install.sh o instalá Go 1.23+." >&2
  exit 1
fi

version="$(go env GOVERSION)"
if [[ "$version" =~ go([0-9]+)\.([0-9]+) ]]; then
  major="${BASH_REMATCH[1]}"
  minor="${BASH_REMATCH[2]}"
  if (( major < 1 || (major == 1 && minor < 23) )); then
    echo "Error: se requiere Go 1.23 o posterior (detectado: $version)." >&2
    exit 1
  fi
else
  echo "No se pudo determinar la versión de Go: $version" >&2
  exit 1
fi

echo "==> Verificando formato Go"
gofmt -w main.go gui.go i18n.go appearance.go process.go progress.go core_test.go i18n_test.go appearance_test.go

echo "==> Ejecutando pruebas"
go test ./...

echo "==> Compilando Eggsmaker Go"
go build -trimpath -buildvcs=false -o eggsmaker .
echo "Compilación terminada: $(pwd)/eggsmaker"
