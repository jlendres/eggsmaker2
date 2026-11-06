#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$(readlink -f "$0")")"

APP_NAME="Eggsmaker_Go"
ARCH="$(uname -m)"
if [[ "$ARCH" != "x86_64" ]]; then
  echo "Error: por ahora el empaquetado AppImage está preparado para x86_64 (detectado: $ARCH)." >&2
  exit 1
fi

for cmd in go curl file pkg-config; do
  command -v "$cmd" >/dev/null 2>&1 || { echo "Error: falta el comando requerido: $cmd" >&2; exit 1; }
done

if [[ "$(go env GOOS)" != "linux" || "$(go env GOARCH)" != "amd64" ]]; then
  echo "Error: compilá desde Linux x86_64 (GOOS=linux, GOARCH=amd64)." >&2
  exit 1
fi

VERSION="$(sed -nE 's/^const AppVersion = "([^"]+)"/\1/p' main.go | head -n1)"
if [[ -z "$VERSION" ]]; then
  echo "Error: no se pudo obtener AppVersion desde main.go." >&2
  exit 1
fi

BUILD_DIR="$PWD/build/appimage"
TOOLS_DIR="$BUILD_DIR/tools"
APPDIR="$BUILD_DIR/AppDir"
OUTPUT_DIR="$PWD/dist"
mkdir -p "$TOOLS_DIR" "$OUTPUT_DIR"

if ! command -v gcc >/dev/null 2>&1 || ! pkg-config --exists gl x11 xcursor xrandr xinerama xi xkbcommon 2>/dev/null; then
  echo "Aviso: Fyne necesita toolchain C y bibliotecas de desarrollo OpenGL/X11." >&2
  echo "En Debian: build-essential pkg-config libgl-dev libglx-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxkbcommon-dev" >&2
  echo "En Arch: base-devel pkgconf mesa libglvnd libx11 libxcursor libxrandr libxinerama libxi libxkbcommon" >&2
  exit 1
fi

fetch_tool() {
  local target="$1" url="$2"
  if [[ ! -x "$target" ]]; then
    echo "==> Descargando $(basename "$target")"
    curl --fail --location --retry 3 --output "$target.tmp" "$url"
    mv "$target.tmp" "$target"
    chmod +x "$target"
  fi
}

LINUXDEPLOY="${LINUXDEPLOY:-$TOOLS_DIR/linuxdeploy-x86_64.AppImage}"
APPIMAGETOOL="${APPIMAGETOOL:-$TOOLS_DIR/appimagetool-x86_64.AppImage}"
fetch_tool "$LINUXDEPLOY" "https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-x86_64.AppImage"
fetch_tool "$APPIMAGETOOL" "https://github.com/AppImage/appimagetool/releases/download/continuous/appimagetool-x86_64.AppImage"

echo "==> Ejecutando pruebas"
go test ./...

echo "==> Compilando Eggsmaker Go $VERSION"
CGO_ENABLED=1 go build -trimpath -buildvcs=false -o "$BUILD_DIR/eggsmaker" .

rm -rf "$APPDIR"
mkdir -p "$APPDIR/usr/bin" "$APPDIR/usr/share/applications" "$APPDIR/usr/share/icons/hicolor/256x256/apps"
install -m755 "$BUILD_DIR/eggsmaker" "$APPDIR/usr/bin/eggsmaker"
sed 's|^Exec=.*|Exec=eggsmaker|' eggsmaker.desktop > "$APPDIR/usr/share/applications/eggsmaker.desktop"
chmod 644 "$APPDIR/usr/share/applications/eggsmaker.desktop"
install -m644 eggsmaker.png "$APPDIR/usr/share/icons/hicolor/256x256/apps/eggsmaker.png"

# linuxdeploy recopila las bibliotecas compartidas necesarias y crea AppRun.
# linuxdeploy busca el empaquetador con el nombre exacto "appimagetool".
ln -sf "$APPIMAGETOOL" "$(dirname "$APPIMAGETOOL")/appimagetool"
export PATH="$(dirname "$APPIMAGETOOL"):$PATH"
export APPIMAGE_EXTRACT_AND_RUN=1
export ARCH=x86_64
export LDAI_OUTPUT="$OUTPUT_DIR/${APP_NAME}-${VERSION}-x86_64.AppImage"
"$LINUXDEPLOY" \
  --appdir "$APPDIR" \
  --executable "$APPDIR/usr/bin/eggsmaker" \
  --desktop-file "$APPDIR/usr/share/applications/eggsmaker.desktop" \
  --icon-file "$APPDIR/usr/share/icons/hicolor/256x256/apps/eggsmaker.png" \
  --output appimage

# El plugin oficial permite fijar explícitamente la ruta de salida.
if [[ ! -s "$OUTPUT_DIR/${APP_NAME}-${VERSION}-x86_64.AppImage" ]]; then
  echo "Error: linuxdeploy terminó pero no se encontró el AppImage generado." >&2
  exit 1
fi
chmod 755 "$OUTPUT_DIR/${APP_NAME}-${VERSION}-x86_64.AppImage"
file "$OUTPUT_DIR/${APP_NAME}-${VERSION}-x86_64.AppImage"
printf "\nAppImage generado: $OUTPUT_DIR/${APP_NAME}-${VERSION}-x86_64.AppImage"
echo "Para ejecutarlo: chmod +x dist/${APP_NAME}-${VERSION}-x86_64.AppImage && ./dist/${APP_NAME}-${VERSION}-x86_64.AppImage"
