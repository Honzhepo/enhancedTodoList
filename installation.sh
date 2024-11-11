#!/bin/bash

# Define directories and file paths
CONFIG_DIR="$HOME/.config/enhancedTodoList"
ASSETS_DIR="./assets"  # Assumes `assets` folder is in the current directory
BIN_DIR="/usr/local/bin"
ICON_DIR="$HOME/.local/share/icons"
APP_ICON="enhancedTodo.png"  # Name of your icon file (ensure it's in the same directory as the script)

# Step 1: Build the Go application
echo "Building the Go application..."
go build -o enhancedTodoList || { echo "Build failed"; exit 1; }

# Step 2: Move the binary to $HOME/bin
echo "Installing the binary to $BIN_DIR..."
sudo mkdir -p "$BIN_DIR"
sudo mv enhancedTodoList "$BIN_DIR/enhancedTodoList"

# Step 3: Create the configuration directory and copy assets
echo "Setting up configuration directory at $CONFIG_DIR..."
mkdir -p "$CONFIG_DIR"
cp -r "$ASSETS_DIR/"* "$CONFIG_DIR"

# Step 4: Place the icon in the user's local icon directory
echo "Installing icon to $ICON_DIR..."
mkdir -p "$ICON_DIR"
cp "$APP_ICON" "$ICON_DIR/enhancedTodo.png"

# Step 5: Create the .desktop entry
echo "Creating .desktop entry..."
mkdir -p "$HOME/.local/share/applications"
cat <<EOF > "$HOME/.local/share/applications/enhancedTodo.desktop"
[Desktop Entry]
Name=Enhanced Todo
Comment=A simple todo list application
Exec=$HOME/bin/enhancedTodoList
Icon=$ICON_DIR/enhancedTodo.png
Terminal=false
Type=Application
Categories=Utility;
EOF

chmod +x "$HOME/.local/share/applications/enhancedTodo.desktop"

# Step 6: Refresh the desktop database (optional but recommended)
echo "Updating desktop database..."
update-desktop-database ~/.local/share/applications

echo "Installation complete! You can now launch 'Enhanced Todo' from Rofi or your desktop environment's application launcher."
