#!/usr/bin/env bash
#
# Preparar dependencias de compilacion Go + GLFW/OpenGL
# para Soplos OS basado en Debian (Testing/Forky y derivados).
#
# Uso:
#   chmod +x preparar-go-grafico-debian.sh
#   ./preparar-go-grafico-debian.sh
#
# El script no instala Go ni modifica el proyecto. Solo instala
# herramientas y cabeceras de desarrollo necesarias para compilar
# aplicaciones Go que usan go-gl/glfw con X11 y OpenGL/GLX.

set -Eeuo pipefail

readonly REQUIRED_HEADERS=(
  "/usr/include/X11/X.h"
  "/usr/include/GL/glx.h"
  "/usr/include/GL/gl.h"
)

if [[ "${EUID}" -eq 0 ]]; then
  SUDO=()
elif command -v sudo >/dev/null 2>&1; then
  SUDO=(sudo)
else
  echo "ERROR: ejecuta este script como root o instala sudo." >&2
  exit 1
fi

if [[ ! -r /etc/debian_version ]] || ! command -v apt-get >/dev/null 2>&1; then
  echo "ERROR: este script está preparado para Debian y derivados con APT." >&2
  exit 1
fi

echo "==> Sistema detectado: $(. /etc/os-release; echo "${PRETTY_NAME:-Debian derivado}")"
echo "==> Actualizando la lista de paquetes..."
"${SUDO[@]}" apt-get update

echo "==> Instalando herramientas y dependencias de desarrollo..."
"${SUDO[@]}" apt-get install -y --no-install-recommends \
  build-essential \
  pkg-config \
  x11proto-dev \
  libx11-dev \
  libxrandr-dev \
  libxinerama-dev \
  libxcursor-dev \
  libxi-dev \
  libxxf86vm-dev \
  libgl-dev \
  libglx-dev

missing=()
for header in "${REQUIRED_HEADERS[@]}"; do
  [[ -f "$header" ]] || missing+=("$header")
done

if ((${#missing[@]} > 0)); then
  echo
  echo "AVISO: faltan cabeceras después de la instalación:"
  printf '  - %s\n' "${missing[@]}"
  echo
  echo "==> Reinstalando los paquetes que proporcionan las cabeceras..."
  "${SUDO[@]}" apt-get install -y --reinstall x11proto-dev libgl-dev libglx-dev

  missing=()
  for header in "${REQUIRED_HEADERS[@]}"; do
    [[ -f "$header" ]] || missing+=("$header")
  done
fi

if ((${#missing[@]} > 0)); then
  echo
  echo "ERROR: siguen faltando estas cabeceras:"
  printf '  - %s\n' "${missing[@]}"
  echo
  echo "Revisa los repositorios y el estado de los paquetes con:"
  echo "  dpkg -l x11proto-dev libx11-dev libgl-dev libglx-dev"
  echo "  dpkg -V x11proto-dev libgl-dev libglx-dev"
  exit 1
fi

echo
echo "==> Verificación de cabeceras:"
for header in "${REQUIRED_HEADERS[@]}"; do
  ls -l "$header"
done

echo
echo "==> Dependencias gráficas instaladas correctamente."
echo "Ya puedes entrar en el directorio del proyecto Go y ejecutar:"
echo "  go test ./..."
echo
echo "Nota: este script no instala Go ni dependencias propias del proyecto."
