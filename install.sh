#!/bin/sh
# AutoLang installer — downloads the latest release binary for your platform.
# Usage: curl -fsSL https://raw.githubusercontent.com/Kyeong6/autolang/main/install.sh | sh

set -e

REPO="Kyeong6/autolang"
BIN="autolang"
INSTALL_DIR="${AUTOLANG_INSTALL_DIR:-/usr/local/bin}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

case "$OS" in
  darwin|linux) ;;
  *)
    echo "Unsupported OS: $OS" >&2
    echo "On Windows, use WSL or download from:" >&2
    echo "  https://github.com/$REPO/releases/latest" >&2
    exit 1
    ;;
esac

# Fetch latest version tag
VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$VERSION" ]; then
  echo "Failed to fetch latest version." >&2
  exit 1
fi

ARCHIVE="${BIN}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$ARCHIVE"

echo "Installing $BIN $VERSION ($OS/$ARCH)..."

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

curl -fsSL "$URL" -o "$TMP/$ARCHIVE"
tar -xzf "$TMP/$ARCHIVE" -C "$TMP"

# Install binary
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP/$BIN" "$INSTALL_DIR/$BIN"
else
  sudo mv "$TMP/$BIN" "$INSTALL_DIR/$BIN"
fi

chmod +x "$INSTALL_DIR/$BIN"

echo ""
echo "✓ $BIN installed to $INSTALL_DIR/$BIN"
echo ""
echo "Next steps:"
echo "  1. Set your translation provider:"
echo "     export AUTOLANG_PROVIDER=deepl"
echo "     export AUTOLANG_API_KEY=your-key"
echo ""
echo "  2. Add shell integration to ~/.zshrc:"
echo "     echo 'eval \"\$(autolang init)\"' >> ~/.zshrc"
echo "     source ~/.zshrc"
