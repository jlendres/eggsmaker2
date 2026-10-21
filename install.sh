#!/bin/bash
#
# Eggsmaker 2 Universal Installer
# Supports Arch, Debian/Ubuntu, Fedora, and other major distributions
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Installation mode
SYSTEM_INSTALL=false

# Check for active virtual environment
if [[ -n "$VIRTUAL_ENV" ]]; then
    echo -e "${YELLOW}Virtual environment detected: $VIRTUAL_ENV${NC}"
    echo -e "${YELLOW}Unloading virtual environment for installation...${NC}"
    
    OLD_VENV="$VIRTUAL_ENV"
    
    # Remove venv bin from PATH
    PATH=$(echo "$PATH" | sed -e "s|$VIRTUAL_ENV/bin:||g" -e "s|:$VIRTUAL_ENV/bin||g")
    unset VIRTUAL_ENV
    
    # Note: We cannot truly 'deactivate' the parent shell's environment from a script,
    # but we can ensure this script runs with a clean environment.
fi

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --system)
            SYSTEM_INSTALL=true
            shift
            ;;
        --help|-h)
            echo "Eggsmaker 2 Installer"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --system    Install system-wide (requires root)"
            echo "  --help      Show this help message"
            echo ""
            echo "Without --system, installs to ~/.local/"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}╔════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Eggsmaker 2 Installer             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════╝${NC}"
echo ""

# Detect distribution
if [ -f /etc/os-release ]; then
    . /etc/os-release
    DISTRO_ID=$ID
    DISTRO_NAME=$NAME
else
    echo -e "${RED}Cannot detect distribution${NC}"
    exit 1
fi

echo -e "${GREEN}Detected: $DISTRO_NAME${NC}"

# Check if running as root for system install
if [ "$SYSTEM_INSTALL" = true ] && [ "$EUID" -ne 0 ]; then
    echo -e "${RED}System install requires root privileges${NC}"
    echo "Please run with sudo: sudo $0 --system"
    exit 1
fi

# Detect package manager and install dependencies
echo -e "\n${YELLOW}Installing system dependencies...${NC}"

case $DISTRO_ID in
    arch|manjaro|endeavouros)
        if [ "$SYSTEM_INSTALL" = true ]; then
            pacman -S --needed --noconfirm python python-pip git rsync
        else
            echo "Checking dependencies..."
            for pkg in python python-pip git rsync; do
                if ! pacman -Q $pkg &>/dev/null; then
                    echo -e "${YELLOW}Please install $pkg manually: sudo pacman -S $pkg${NC}"
                fi
            done
        fi
        ;;
    debian|ubuntu|linuxmint|pop)
        if [ "$SYSTEM_INSTALL" = true ]; then
            apt-get update
            apt-get install -y python3 python3-pip python3-venv git rsync
        else
            echo "Checking dependencies..."
            for pkg in python3 python3-pip python3-venv git rsync; do
                if ! dpkg -l | grep -q "^ii  $pkg "; then
                    echo -e "${YELLOW}Please install $pkg manually: sudo apt install $pkg${NC}"
                fi
            done
        fi
        ;;
    fedora|rhel|centos|almalinux|rocky)
        if [ "$SYSTEM_INSTALL" = true ]; then
            dnf install -y python3 python3-pip git rsync
        else
            echo "Checking dependencies..."
            for pkg in python3 python3-pip git rsync; do
                if ! rpm -q $pkg &>/dev/null; then
                    echo -e "${YELLOW}Please install $pkg manually: sudo dnf install $pkg${NC}"
                fi
            done
        fi
        ;;
    *)
        echo -e "${YELLOW}Unknown distribution: $DISTRO_ID${NC}"
        echo "Please ensure python3, pip, git, and rsync are installed"
        ;;
esac

# Determine installation paths
if [ "$SYSTEM_INSTALL" = true ]; then
    INSTALL_DIR="/opt/eggsmaker2"
    BIN_DIR="/usr/local/bin"
    DESKTOP_DIR="/usr/share/applications"
    ICON_DIR="/usr/share/pixmaps"
else
    INSTALL_DIR="$HOME/.local/share/eggsmaker2"
    BIN_DIR="$HOME/.local/bin"
    DESKTOP_DIR="$HOME/.local/share/applications"
    ICON_DIR="$HOME/.local/share/icons"
fi

echo -e "\n${YELLOW}Installing to: $INSTALL_DIR${NC}"

# Create directories
mkdir -p "$INSTALL_DIR"
mkdir -p "$BIN_DIR"
mkdir -p "$DESKTOP_DIR"
mkdir -p "$ICON_DIR"

# Copy application files
echo -e "\n${YELLOW}Copying application files...${NC}"
cp -r . "$INSTALL_DIR/"

# Create virtual environment
echo -e "\n${YELLOW}Creating Python virtual environment...${NC}"
python3 -m venv "$INSTALL_DIR/venv"

# Activate venv and install dependencies
echo -e "\n${YELLOW}Installing Python dependencies...${NC}"
source "$INSTALL_DIR/venv/bin/activate"
pip install --upgrade pip
pip install -r "$INSTALL_DIR/requirements.txt"
deactivate

# Create launcher wrapper
echo -e "\n${YELLOW}Creating launcher...${NC}"
cat > "$BIN_DIR/eggsmaker2" <<EOF
#!/bin/bash
# Eggsmaker 2 Launcher Wrapper
source "$INSTALL_DIR/venv/bin/activate"
cd "$INSTALL_DIR"
python3 "$INSTALL_DIR/main.py" "\$@"
EOF

chmod +x "$BIN_DIR/eggsmaker2"

# Install desktop file
echo -e "\n${YELLOW}Installing desktop entry...${NC}"
sed "s|Exec=eggsmaker2|Exec=$BIN_DIR/eggsmaker2|g" eggsmaker2.desktop > "$DESKTOP_DIR/eggsmaker2.desktop"
chmod +x "$DESKTOP_DIR/eggsmaker2.desktop"

# Install icon
if [ -f "assets/assets/eggsmaker2.png" ]; then
    cp "assets/assets/eggsmaker2.png" "$ICON_DIR/eggsmaker2.png"
fi

# Update desktop database
if command -v update-desktop-database &> /dev/null; then
    echo -e "\n${YELLOW}Updating desktop database...${NC}"
    if [ "$SYSTEM_INSTALL" = true ]; then
        update-desktop-database "$DESKTOP_DIR"
    else
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    fi
fi

echo -e "\n${GREEN}╔════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   Installation Complete!           ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}To launch Eggsmaker 2:${NC}"
echo -e "  ${GREEN}eggsmaker2${NC}"
echo ""
echo -e "${BLUE}Or find it in your application menu${NC}"
echo ""



if [[ -n "$OLD_VENV" ]]; then
     echo -e "${YELLOW}Note: You were in a virtual environment before running this script.${NC}"
     echo -e "${YELLOW}Please run 'deactivate' and 'source <venv>/bin/activate' to refresh if needed.${NC}"
fi
