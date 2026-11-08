#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "$(readlink -f -- "$0")")" && pwd)"
cd "$script_dir"

for command in go dpkg-deb dpkg; do
  if ! command -v "$command" >/dev/null 2>&1; then
    echo "Falta el comando requerido: $command" >&2
    exit 1
  fi
done

version="$(sed -nE 's/^const AppVersion = "([^"]+)"$/\1/p' main.go)"
if [[ -z "$version" ]]; then
  echo "No se pudo determinar AppVersion en main.go." >&2
  exit 1
fi

architecture="$(dpkg --print-architecture)"
build_dir="$(mktemp -d)"
trap 'rm -rf -- "$build_dir"' EXIT
package_root="$build_dir/eggsmaker-go"
mkdir -p "$package_root/DEBIAN" "$script_dir/dist"

echo "==> Ejecutando pruebas"
go test ./...
echo "==> Compilando para Debian ($architecture)"
go build -trimpath -buildvcs=false -o "$build_dir/eggsmaker" .

install -Dm755 "$build_dir/eggsmaker" "$package_root/usr/bin/eggsmaker"
install -Dm644 eggsmaker.desktop "$package_root/usr/share/applications/eggsmaker.desktop"
install -Dm644 eggsmaker.png "$package_root/usr/share/icons/hicolor/256x256/apps/eggsmaker.png"
install -Dm644 README.md "$package_root/usr/share/eggsmaker/README.md"

cat > "$package_root/DEBIAN/control" <<EOF
Package: eggsmaker-go
Version: ${version}
Section: utils
Priority: optional
Architecture: ${architecture}
Depends: libc6, libgl1, libx11-6, libxcursor1, libxrandr2, libxinerama1, libxi6, libxxf86vm1, libxkbcommon0
Maintainer: Eggsmaker Go contributors
Description: Graphical interface for Penguins Eggs ISO creation
 Eggsmaker Go builds custom ISO images with Penguins Eggs and provides
 an integrated graphical configuration editor and package manager.
EOF

output="$script_dir/dist/eggsmaker-go_${version}_${architecture}.deb"
dpkg-deb --build --root-owner-group "$package_root" "$output"
printf 'Paquete Debian creado: %s\n' "$output"
