BINARY_NAME := eggsmaker
PREFIX ?= /usr/local
BINDIR := $(PREFIX)/bin
APPDIR := $(PREFIX)/share/applications
DATADIR := $(PREFIX)/share/eggsmaker
ICONDIR := $(PREFIX)/share/icons/hicolor/256x256/apps

.PHONY: all build appimage install install-files uninstall uninstall-files clean test
all: build

build:
	go build -trimpath -buildvcs=false -o $(BINARY_NAME) .

appimage:
	./compilar_appimage.sh

# install-files instala un binario ya compilado; permite compilar como usuario
# normal y elevar privilegios sólo para copiar archivos al sistema.
install: build install-files

install-files:
	install -Dm755 $(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	install -Dm644 eggsmaker.desktop $(DESTDIR)$(APPDIR)/$(BINARY_NAME).desktop
	install -Dm644 eggsmaker.png $(DESTDIR)$(ICONDIR)/eggsmaker.png
	install -Dm644 README.md $(DESTDIR)$(DATADIR)/README.md
	@if command -v update-desktop-database >/dev/null 2>&1; then update-desktop-database $(DESTDIR)$(APPDIR) 2>/dev/null || true; fi
	@echo "Eggsmaker Go instalado. Ejecutalo como usuario normal desde el menú o con: eggsmaker"

uninstall: uninstall-files

uninstall-files:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	rm -f $(DESTDIR)$(APPDIR)/$(BINARY_NAME).desktop
	rm -f $(DESTDIR)$(ICONDIR)/eggsmaker.png
	rm -rf $(DESTDIR)$(DATADIR)
	@if command -v update-desktop-database >/dev/null 2>&1; then update-desktop-database $(DESTDIR)$(APPDIR) 2>/dev/null || true; fi

clean:
	rm -f $(BINARY_NAME)

test:
	go test ./...
