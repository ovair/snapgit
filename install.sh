#!/bin/bash
# SnapGit installer for Linux and macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/ovair/snapgit/main/install.sh | bash
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

# Build asset name (strip leading 'v' from version)
STRIPPED_VERSION="${VERSION#v}"
ASSET="sg_${STRIPPED_VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$ASSET"
CHECKSUMS_URL="https://github.com/$REPO/releases/download/$VERSION/checksums.txt"
TMP="$(mktemp -d)"

echo "Downloading $ASSET ..."
curl -fsSL "$URL" -o "$TMP/$ASSET"
curl -fsSL "$CHECKSUMS_URL" -o "$TMP/checksums.txt"

# Verify checksum
echo "Verifying checksum ..."
EXPECTED="$(grep "$ASSET" "$TMP/checksums.txt" | awk '{print $1}')"
if [ -z "$EXPECTED" ]; then
    echo "Error: checksum not found for $ASSET"
    rm -rf "$TMP"
    exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
    ACTUAL="$(sha256sum "$TMP/$ASSET" | awk '{print $1}')"
elif command -v shasum >/dev/null 2>&1; then
    ACTUAL="$(shasum -a 256 "$TMP/$ASSET" | awk '{print $1}')"
else
    echo "Warning: no sha256sum or shasum found, skipping checksum verification"
    ACTUAL="$EXPECTED"
fi

if [ "$ACTUAL" != "$EXPECTED" ]; then
    echo "Error: checksum mismatch!"
    echo "  Expected: $EXPECTED"
    echo "  Got:      $ACTUAL"
    rm -rf "$TMP"
    exit 1
fi

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
