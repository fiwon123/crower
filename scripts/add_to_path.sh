#!/bin/sh
# Add the directory containing this script to PATH in the correct shell config.

# Get the directory where this script is located (POSIX compatible)
SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"

# Detect shell
DETECTED_SHELL=""
if [ -n "$SHELL" ]; then
    DETECTED_SHELL="$(basename "$SHELL")"
fi
if [ -z "$DETECTED_SHELL" ]; then
    DETECTED_SHELL="$(basename "$(getent passwd "$(id -un)" 2>/dev/null | cut -d: -f7)")" || true
fi
if [ -z "$DETECTED_SHELL" ]; then
    DETECTED_SHELL="sh"
fi

# Check if already in PATH in any common config file
if grep -Fq "$SCRIPT_DIR" "$HOME/.bashrc" 2>/dev/null || \
   grep -Fq "$SCRIPT_DIR" "$HOME/.zshrc" 2>/dev/null || \
   grep -Fq "$SCRIPT_DIR" "$HOME/.profile" 2>/dev/null || \
   grep -Fq "$SCRIPT_DIR" "$HOME/.cshrc" 2>/dev/null || \
   grep -Fq "$SCRIPT_DIR" "$HOME/.config/fish/config.fish" 2>/dev/null; then
    echo "Folder $SCRIPT_DIR already in PATH"
    exit 0
fi

case "$DETECTED_SHELL" in
    zsh)
        CONFIG_FILE="$HOME/.zshrc"
        LINE="export PATH=\"$SCRIPT_DIR:\$PATH\""
        SOURCE_CMD="source ~/.zshrc"
        ;;
    bash)
        CONFIG_FILE="$HOME/.bashrc"
        LINE="export PATH=\"$SCRIPT_DIR:\$PATH\""
        SOURCE_CMD="source ~/.bashrc"
        ;;
    fish)
        mkdir -p "$HOME/.config/fish"
        CONFIG_FILE="$HOME/.config/fish/config.fish"
        LINE="set -gx PATH \"$SCRIPT_DIR\" \$PATH"
        SOURCE_CMD="source ~/.config/fish/config.fish"
        ;;
    csh|tcsh)
        CONFIG_FILE="$HOME/.cshrc"
        LINE="setenv PATH \"$SCRIPT_DIR:\$PATH\""
        SOURCE_CMD="source ~/.cshrc"
        ;;
    *)
        CONFIG_FILE="$HOME/.profile"
        LINE="export PATH=\"$SCRIPT_DIR:\$PATH\""
        SOURCE_CMD=". ~/.profile"
        ;;
esac

echo "$LINE" >> "$CONFIG_FILE"
echo "Folder $SCRIPT_DIR added to PATH in $CONFIG_FILE"
echo "Run '$SOURCE_CMD' or restart terminal to apply changes."
