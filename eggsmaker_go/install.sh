#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$(readlink -f "$0")")"

if [[ ${EUID:-$(id -u)} -eq 0 ]]; then
  echo "Ejecutá ./install.sh como usuario normal; se solicitará sudo sólo al instalar."
  exit 1
fi
if ! command -v sudo >/dev/null 2>&1; then
  echo "Se requiere sudo para instalar los archivos del sistema." >&2
  exit 1
fi

os_release_file="${OS_RELEASE_FILE:-/etc/os-release}"
if [[ -r "$os_release_file" ]]; then
  # shellcheck disable=SC1090
  source "$os_release_file"
fi
family="${ID:-} ${ID_LIKE:-}"

if ! command -v go >/dev/null 2>&1; then
  if [[ "$family" =~ (arch|cachyos|manjaro|endeavouros|garuda|arcolinux|artix) ]]; then
    sudo pacman -S --needed go
  elif [[ "$family" =~ (debian|ubuntu|mint|biglinux|pop|kali|mx|zorin|elementary) ]]; then
    sudo apt-get update
    sudo apt-get install -y golang-go
  else
    echo "Distribución no reconocida. Instalá Go 1.23 o posterior manualmente." >&2
    exit 1
  fi
fi

version="$(go env GOVERSION)"
if [[ "$version" =~ go([0-9]+)\.([0-9]+) ]]; then
  major="${BASH_REMATCH[1]}"; minor="${BASH_REMATCH[2]}"
  if (( major < 1 || (major == 1 && minor < 23) )); then
    echo "Se requiere Go 1.23 o posterior (detectado: $version)." >&2
    echo "En Debian estable puede ser necesario instalar Go desde https://go.dev/dl/." >&2
    exit 1
  fi
else
  echo "No se pudo determinar la versión de Go: $version" >&2
  exit 1
fi

if [[ "$family" =~ (arch|cachyos|manjaro|endeavouros|garuda|arcolinux|artix) ]]; then
  packages=(base-devel gcc make pkgconf mesa libglvnd libx11 libxcursor libxrandr libxinerama libxi libxkbcommon)
  missing=()
  for dep in "${packages[@]}"; do
    pacman -Q "$dep" >/dev/null 2>&1 || missing+=("$dep")
  done
  if ((${#missing[@]})); then sudo pacman -S --needed "${missing[@]}"; fi
elif [[ "$family" =~ (debian|ubuntu|mint|biglinux|pop|kali|mx|zorin|elementary) ]]; then
  packages=(build-essential pkg-config libgl-dev libglx-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxkbcommon-dev)
  missing=()
  for dep in "${packages[@]}"; do
    dpkg-query -W -f='${Status}' "$dep" 2>/dev/null | grep -q 'install ok installed' || missing+=("$dep")
  done
  if ((${#missing[@]})); then
    sudo apt-get update
    sudo apt-get install -y "${missing[@]}"
  fi
else
  echo "Distribución no reconocida. Instalá las dependencias de compilación de Fyne manualmente." >&2
  exit 1
fi

./compilar.sh
sudo make install-files PREFIX=/usr/local
printf '\nInstalación completada. Eggsmaker Go aparecerá en el menú de aplicaciones.\n'
