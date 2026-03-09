#!/bin/sh
# SnapGit installer for Linux and macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/ovair/snapgit/main/install.sh | sh
set -e

REPO="ovair/snapgit"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS="$(uname -s)"
case "$OS" in
    Linux*)  OS="linux" ;;
    Darwin*) OS="darwin" ;;
    *)       echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *)       echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest release tag
VERSION="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)"
if [ -z "$VERSION" ]; then
    echo "Error: could not determine latest version."
    exit 1
fi

echo "Installing SnapGit $VERSION ..."

# Download and extract
ASSET="sg_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$ASSET"
TMP="$(mktemp -d)"

echo "Downloading $ASSET ..."
curl -fsSL "$URL" -o "$TMP/$ASSET"
tar -xzf "$TMP/$ASSET" -C "$TMP"

# Install
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP/sg" "$INSTALL_DIR/sg"
else
    echo "Installing to $INSTALL_DIR (requires sudo) ..."
    sudo mv "$TMP/sg" "$INSTALL_DIR/sg"
fi

chmod +x "$INSTALL_DIR/sg"
rm -rf "$TMP"

echo "SnapGit $VERSION installed successfully!"
echo "Run 'sg help' to get started."
