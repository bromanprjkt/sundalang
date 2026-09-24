#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
WHITE='\033[1;37m'
NC='\033[0m'

get_term_width() {
    if command -v tput &> /dev/null; then
        tput cols 2>/dev/null || echo 80
    else
        echo 80
    fi
}

center_text() {
    local text="$1"
    local color="${2:-$NC}"
    local width=$(get_term_width)
    local text_len=${#text}
    local padding=$(( (width - text_len) / 2 ))
    [ $padding -lt 0 ] && padding=0
    printf "%${padding}s" ""
    echo -e "${color}${text}${NC}"
}

draw_box() {
    local width=$(get_term_width)
    local box_width=44
    local padding=$(( (width - box_width) / 2 ))
    [ $padding -lt 0 ] && padding=0
    local pad=$(printf "%${padding}s" "")
    
    echo -e "${CYAN}${pad}╭──────────────────────────────────────────╮${NC}"
    echo -e "${CYAN}${pad}│                                          │${NC}"
    echo -e "${CYAN}${pad}│${WHITE}            ⬢  SUNDALANG  ⬢              ${CYAN}│${NC}"
    echo -e "${CYAN}${pad}│                                          │${NC}"
    echo -e "${CYAN}${pad}│${NC}   Bahasa Pemrograman Sunda Pandeglang   ${CYAN}│${NC}"
    echo -e "${CYAN}${pad}│${NC}          Installer for Unix/Mac         ${CYAN}│${NC}"
    echo -e "${CYAN}${pad}│                                          │${NC}"
    echo -e "${CYAN}${pad}╰──────────────────────────────────────────╯${NC}"
}

draw_success_box() {
    local width=$(get_term_width)
    local box_width=36
    local padding=$(( (width - box_width) / 2 ))
    [ $padding -lt 0 ] && padding=0
    local pad=$(printf "%${padding}s" "")
    
    echo -e "${GREEN}${pad}╭──────────────────────────────────╮${NC}"
    echo -e "${GREEN}${pad}│                                  │${NC}"
    echo -e "${GREEN}${pad}│${WHITE}      ✓ INSTALASI SUKSES!        ${GREEN}│${NC}"
    echo -e "${GREEN}${pad}│                                  │${NC}"
    echo -e "${GREEN}${pad}╰──────────────────────────────────╯${NC}"
}

draw_divider() {
    local width=$(get_term_width)
    local div_len=50
    [ $div_len -gt $((width - 10)) ] && div_len=$((width - 10))
    local padding=$(( (width - div_len) / 2 ))
    [ $padding -lt 0 ] && padding=0
    printf "%${padding}s" ""
    printf "${GRAY}"
    for ((i=0; i<div_len; i++)); do printf "─"; done
    printf "${NC}\n"
}

status_msg() {
    local status="$1"
    local text="$2"
    local color="$3"
    local msg="[${status}] ${text}"
    center_text "$msg" "$color"
}

clear
echo ""
draw_box
echo ""

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    linux*)
        OS="linux"
        ;;
    darwin*)
        OS="darwin"
        ;;
    *)
        echo ""
        status_msg "ERROR" "OS teu disupport: $OS" "$RED"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo ""
        status_msg "ERROR" "Architecture teu disupport: $ARCH" "$RED"
        exit 1
        ;;
esac

if [ "$OS" = "darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        BINARY_NAME="sundalang-macos-arm64"
    else
        BINARY_NAME="sundalang-macos"
    fi
elif [ "$OS" = "linux" ]; then
    if [ "$ARCH" = "arm64" ]; then
        BINARY_NAME="sundalang-linux-arm64"
    else
        BINARY_NAME="sundalang"
    fi
fi

center_text "Platform: $OS-$ARCH" "$YELLOW"
echo ""

INSTALL_DIR="$HOME/.sundalang/bin"
BINARY_PATH="$INSTALL_DIR/sundalang"

if [ ! -d "$INSTALL_DIR" ]; then
    status_msg "INFO" "Nyieun folder instalasi..." "$YELLOW"
    mkdir -p "$INSTALL_DIR"
    status_msg "OK" "Folder instalasi dibuat" "$GREEN"
fi

echo ""
draw_divider
echo ""
status_msg "INFO" "Nyari versi terbaru..." "$YELLOW"
echo ""

RELEASE_INFO=$(curl -fsSL https://api.github.com/repos/bromanprjkt/sundalang/releases/latest)
VERSION=$(echo "$RELEASE_INFO" | grep '"tag_name"' | cut -d '"' -f 4)

if [ -z "$VERSION" ]; then
    status_msg "ERROR" "Gagal manggihan versi terbaru" "$RED"
    exit 1
fi

status_msg "OK" "Kapanggih: SundaLang $VERSION" "$GREEN"

DOWNLOAD_URL=$(echo "$RELEASE_INFO" | grep "browser_download_url.*$BINARY_NAME\"" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo ""
    status_msg "ERROR" "Gagal manggihan binary pikeun $OS-$ARCH" "$RED"
    exit 1
fi

echo ""
draw_divider
echo ""
status_msg "INFO" "Ngeundeur SundaLang $VERSION..." "$YELLOW"
echo ""

TMP_DIR="${TMPDIR:-/tmp}"
[ ! -d "$TMP_DIR" ] && TMP_DIR="$HOME"
TEMP_FILE="$TMP_DIR/sundalang-$$.tmp"

if ! curl -fsSL -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
    status_msg "ERROR" "Gagal ngundeur binary" "$RED"
    rm -f "$TEMP_FILE"
    exit 1
fi

mv "$TEMP_FILE" "$BINARY_PATH"
chmod +x "$BINARY_PATH"
status_msg "OK" "Hasil ngundeur" "$GREEN"

echo ""
draw_divider
echo ""
status_msg "INFO" "Ngatur PATH..." "$YELLOW"
echo ""

SHELL_RC=""
if [ -n "$ZSH_VERSION" ] || [ -f "$HOME/.zshrc" ]; then
    SHELL_RC="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ] || [ -f "$HOME/.bashrc" ]; then
    SHELL_RC="$HOME/.bashrc"
else
    SHELL_RC="$HOME/.profile"
fi

if ! grep -q ".sundalang/bin" "$SHELL_RC" 2>/dev/null; then
    echo "" >> "$SHELL_RC"
    echo 'export PATH="$HOME/.sundalang/bin:$PATH"' >> "$SHELL_RC"
    status_msg "OK" "PATH geus diupdate dina $SHELL_RC" "$GREEN"
    echo ""
    center_text "PENTING: Jalankeun command ieu:" "$YELLOW"
    center_text "source $SHELL_RC" "$CYAN"
    center_text "Atawa tutup tur buka deui terminal" "$YELLOW"
else
    status_msg "OK" "PATH geus aya SundaLang" "$GREEN"
fi

echo ""
echo ""
draw_success_box
echo ""

center_text "Kumaha carana make:" "$CYAN"
echo ""
center_text "1. Jalankeun: source $SHELL_RC" "$WHITE"
center_text "2. Test: sundalang --version" "$WHITE"
center_text "3. Jalankeun file .sl: sundalang namafile.sl" "$WHITE"
echo ""
draw_divider
echo ""
center_text "Lokasi: $BINARY_PATH" "$GRAY"
echo ""
center_text "Uninstall:" "$GRAY"
center_text "curl -fsSL https://raw.githubusercontent.com/bromanprjkt/sundalang/main/uninstall.sh | bash" "$GRAY"
echo ""
