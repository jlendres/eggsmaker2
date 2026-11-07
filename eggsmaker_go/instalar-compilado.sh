#!/usr/bin/env bash
# Instala Eggsmaker Go desde el binario ya compilado, sin ejecutar go build.
set -euo pipefail
cd "$(dirname "$(readlink -f "$0")")"

if [[ ${EUID:-$(id -u)} -eq 0 ]]; then
  echo "Ejecutá ./instalar-compilado.sh como usuario normal; se solicitará sudo para instalar."
  exit 1
fi
if ! command -v sudo >/dev/null 2>&1; then
  echo "Error: se necesita sudo para instalar la aplicación." >&2
  exit 1
fi
if [[ ! -f ./eggsmaker || ! -x ./eggsmaker ]]; then
  echo "Error: no se encuentra el binario compilado ./eggsmaker." >&2
  exit 1
fi
if [[ ./main.go -nt ./eggsmaker || ./gui.go -nt ./eggsmaker || ./process.go -nt ./eggsmaker ]]; then
  echo "Error: el binario es anterior al código fuente modificado." >&2
  echo "En el equipo de desarrollo ejecutá ./compilar.sh y luego distribuí el binario actualizado." >&2
  exit 1
fi

sudo make install-files PREFIX=/usr/local
printf '\nInstalación del compilado completada. No se realizó ninguna compilación.\n'
