#!/bin/sh
# Tetrix CE CLI Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/well-prado/tetrix-cli/main/scripts/install.sh | sh
set -e

REPO="well-prado/tetrix-cli"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="tetrix"

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

echo "Detecting system: ${OS}/${ARCH}"

# Get latest release tag
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')
if [ -z "$LATEST" ]; then
  echo "Error: Could not determine latest release."
  exit 1
fi

echo "Latest version: ${LATEST}"

# Download
ARCHIVE="${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${LATEST}/${ARCHIVE}"

echo "Downloading ${URL}..."
TMP_DIR=$(mktemp -d)
curl -fsSL "$URL" -o "${TMP_DIR}/${ARCHIVE}"

# Extract
tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "${TMP_DIR}"

# Install
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
rm -rf "$TMP_DIR"

echo ""
echo "Tetrix CLI installed successfully!"
echo ""
echo "  Version: ${LATEST}"
echo "  Binary:  ${INSTALL_DIR}/${BINARY_NAME}"
echo ""
echo "Get started:"
echo "  tetrix install"
echo ""
