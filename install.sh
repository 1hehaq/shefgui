#!/bin/bash

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DESKTOP_FILE="$PROJECT_DIR/shef.desktop"
ICON_FILE="$PROJECT_DIR/evilgopher.png"
BINARY_FILE="$PROJECT_DIR/shef"

echo "🔧 Setting up Shef GUI desktop integration..."

if [ ! -f "$ICON_FILE" ]; then
    echo "❌ Error: evilgopher.png not found in $PROJECT_DIR"
    exit 1
fi

if [ ! -f "$DESKTOP_FILE" ]; then
    echo "❌ Error: shef.desktop not found in $PROJECT_DIR"
    exit 1
fi

if [ ! -f "$BINARY_FILE" ]; then
    echo "🔨 Building shef binary..."
    cd "$PROJECT_DIR"
    go build -o shef .
    echo "✅ Binary built successfully"
fi

chmod +x "$BINARY_FILE"

sed -i "s|Exec=.*|Exec=$BINARY_FILE|g" "$DESKTOP_FILE"
sed -i "s|Icon=.*|Icon=$ICON_FILE|g" "$DESKTOP_FILE"

chmod +x "$DESKTOP_FILE"

if [ -d "$HOME/Desktop" ]; then
    cp "$DESKTOP_FILE" "$HOME/Desktop/"
    chmod +x "$HOME/Desktop/shef.desktop"
    echo "✅ Desktop launcher created: ~/Desktop/shef-gui.desktop"
fi

APPS_DIR="$HOME/.local/share/applications"
mkdir -p "$APPS_DIR"
cp "$DESKTOP_FILE" "$APPS_DIR/"
echo "✅ Application menu entry created: $APPS_DIR/shef-gui.desktop"

if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$APPS_DIR"
    echo "✅ Desktop database updated"
fi

echo ""
echo "🎉 Shef GUI desktop integration complete!"
echo ""
echo "You can now:"
echo "  • Double-click the desktop icon to launch"
echo "  • Find 'Shef GUI' in your applications menu"
echo "  • Search for 'Shodan' in your app launcher"
echo ""
echo "Icon: evilgopher.png"
echo "Binary: $BINARY_FILE"
