#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "$(readlink -f -- "$0")")" && pwd)"
cd "$script_dir"

if [[ ${EUID:-$(id -u)} -eq 0 ]]; then
  echo "No ejecutes makepkg como root. Ejecutá este script como usuario normal." >&2
  exit 1
fi
for command in go makepkg tar zstd; do
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

echo "==> Ejecutando pruebas"
go test ./...

build_dir="$(mktemp -d)"
trap 'rm -rf -- "$build_dir"' EXIT
mkdir -p "$build_dir/source/eggsmaker_go" "$script_dir/dist"
tar -C "$script_dir" \
  --exclude='./.git' \
  --exclude='./eggsmaker' \
  --exclude='./dist' \
  --exclude='./build' \
  --exclude='./*.AppImage' \
  -cf - . | tar -C "$build_dir/source/eggsmaker_go" -xf -
tar -C "$build_dir/source" -czf "$build_dir/eggsmaker-go-${version}.tar.gz" eggsmaker_go

cat > "$build_dir/PKGBUILD" <<EOF
pkgname=eggsmaker-go
pkgver=${version}
pkgrel=1
pkgdesc='Graphical interface for Penguins Eggs ISO creation'
arch=('x86_64' 'aarch64')
url='https://github.com/jlendres/eggsmaker2'
license=('unknown')
depends=('libglvnd' 'libx11' 'libxcursor' 'libxrandr' 'libxinerama' 'libxi' 'libxxf86vm' 'libxkbcommon')
makedepends=('go' 'gcc' 'pkgconf')
source=("eggsmaker-go-\$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
  cd "\$srcdir/eggsmaker_go"
  go build -trimpath -buildvcs=false -o eggsmaker .
}

package() {
  cd "\$srcdir/eggsmaker_go"
  install -Dm755 eggsmaker "\$pkgdir/usr/bin/eggsmaker"
  install -Dm644 eggsmaker.desktop "\$pkgdir/usr/share/applications/eggsmaker.desktop"
  install -Dm644 eggsmaker.png "\$pkgdir/usr/share/icons/hicolor/256x256/apps/eggsmaker.png"
  install -Dm644 README.md "\$pkgdir/usr/share/eggsmaker/README.md"
}
EOF

package_path="$(cd "$build_dir" && PKGDEST="$build_dir" PKGEXT='.pkg.tar.zst' makepkg --packagelist)"
(cd "$build_dir" && PKGDEST="$build_dir" PKGEXT='.pkg.tar.zst' makepkg --force --cleanbuild --noconfirm)

if [[ ! -f "$package_path" ]]; then
  echo "makepkg terminó sin generar el paquete esperado." >&2
  exit 1
fi
output="$script_dir/dist/$(basename "$package_path")"
install -m 0644 "$package_path" "$output"
printf 'Paquete Arch creado: %s\n' "$output"
