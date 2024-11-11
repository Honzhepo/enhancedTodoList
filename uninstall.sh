#!/bin/bash

# Define paths
BIN_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/enhancedTodoList"
ICON_DIR="$HOME/.local/share/icons"
APP_ICON="$ICON_DIR/enhancedTodo.png"
DESKTOP_FILE="$HOME/.local/share/applications/enhancedTodo.desktop"

# Step 1: Remove the binary from $HOME/bin
echo "Removing the binary from $BIN_DIR..."
rm -f "$BIN_DIR/enhancedTodo"

# Step 2: Remove the configuration directory and its contents
echo "Deleting the configuration directory at $CONFIG_DIR..."
rm -rf "$CONFIG_DIR"

# Step 3: Remove the icon from the local icon directory
echo "Removing the icon from $ICON_DIR..."
rm -f "$APP_ICON"

# Step 4: Remove the .desktop entry
echo "Deleting .desktop entry at $DESKTOP_FILE..."
rm -f "$DESKTOP_FILE"

# Step 5: Refresh the desktop database (optional)
echo "Updating desktop database..."
update-desktop-database "$HOME/.local/share/applications"

echo "Uninstallation complete!"
