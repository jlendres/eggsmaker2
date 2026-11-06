#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$(readlink -f "$0")")"

if [[ ${EUID:-$(id -u)} -ne 0 ]]; then
  if ! command -v sudo >/dev/null 2>&1; then
    echo "Se necesitan permisos de administrador. Instalá sudo o ejecutá como root." >&2
    exit 1
  fi
  exec sudo "$0" "$@"
fi

make uninstall-files PREFIX=/usr/local
printf 'Eggsmaker Go fue desinstalado. No se eliminó la carpeta de datos del usuario.\n'
