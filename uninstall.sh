#!/bin/bash
#
# Eggsmaker 2 Uninstaller
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Uninstallation mode
SYSTEM_UNINSTALL=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --system)
            SYSTEM_UNINSTALL=true
            shift
            ;;
        --help|-h)
            echo "Eggsmaker 2 Uninstaller"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --system    Uninstall system-wide (requires root)"
            echo "  --help      Show this help message"
            echo ""
            echo "Without --system, uninstalls from ~/.local/"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}╔════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Eggsmaker 2 Uninstaller          ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════╝${NC}"
echo ""

# Check if running as root for system uninstall
if [ "$SYSTEM_UNINSTALL" = true ] && [ "$EUID" -ne 0 ]; then
    echo -e "${RED}System uninstall requires root privileges${NC}"
    echo "Please run with sudo: sudo $0 --system"
    exit 1
fi

# Determine paths
if [ "$SYSTEM_UNINSTALL" = true ]; then
    INSTALL_DIR="/opt/eggsmaker2"
    BIN_FILE="/usr/local/bin/eggsmaker2"
    DESKTOP_FILE="/usr/share/applications/eggsmaker2.desktop"
    ICON_FILE="/usr/share/pixmaps/eggsmaker2.png"
else
    INSTALL_DIR="$HOME/.local/share/eggsmaker2"
    BIN_FILE="$HOME/.local/bin/eggsmaker2"
    DESKTOP_FILE="$HOME/.local/share/applications/eggsmaker2.desktop"
    ICON_FILE="$HOME/.local/share/icons/eggsmaker2.png"
fi

echo -e "${YELLOW}Removing files...${NC}"

# Remove files
if [ -d "$INSTALL_DIR" ]; then
    echo "Removing $INSTALL_DIR"
    rm -rf "$INSTALL_DIR"
fi

if [ -f "$BIN_FILE" ]; then
    echo "Removing $BIN_FILE"
    rm -f "$BIN_FILE"
fi

if [ -f "$DESKTOP_FILE" ]; then
    echo "Removing $DESKTOP_FILE"
    rm -f "$DESKTOP_FILE"
fi

if [ -f "$ICON_FILE" ]; then
    echo "Removing $ICON_FILE"
    rm -f "$ICON_FILE"
fi

# Update desktop database
if command -v update-desktop-database &> /dev/null; then
    echo -e "${YELLOW}Updating desktop database...${NC}"
    if [ "$SYSTEM_UNINSTALL" = true ]; then
        update-desktop-database /usr/share/applications
    else
        update-desktop-database "$HOME/.local/share/applications" 2>/dev/null || true
    fi
fi

echo -e "\n${GREEN}Uninstallation Complete!${NC}"
