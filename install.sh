#!/bin/bash
set -e

REPO="andragon31/Hati"
INSTALL_DIR="/usr/local/bin"

if [[ "$OSTYPE" == "darwin"* ]]; then
    ARCH=$(uname -m)
    if [ "$ARCH" = "arm64" ]; then
        BIN="hati-darwin-arm64"
    else
        BIN="hati-darwin-amd64"
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    ARCH=$(uname -m)
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
        BIN="hati-linux-arm64"
    else
        BIN="hati-linux-amd64"
    fi
else
    echo "Unsupported OS: $OSTYPE"
    exit 1
fi

TMP=$(mktemp)
URL="https://github.com/${REPO}/releases/latest/download/${BIN}"

echo "Downloading Hati..."
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP" "$INSTALL_DIR/hati"
    echo "Installed to $INSTALL_DIR/hati"
else
    echo "Installing to $INSTALL_DIR requires sudo..."
    sudo mv "$TMP" "$INSTALL_DIR/hati"
    echo "Installed to $INSTALL_DIR/hati"
fi

echo ""
echo "Hati installed! Run:"
echo "  hati version          # Verify"
echo "  hati init            # Initialize"
echo ""
