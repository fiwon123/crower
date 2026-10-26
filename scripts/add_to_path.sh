#!/bin/sh
# Add the directory containing this script to PATH, in the one config file the
# user's own shell actually reads.
#
# A line in another shell's config does not help anyone: a zsh user never reads
# ~/.bashrc. Detection and the idempotency check therefore both target a single
# file, so a stale entry left by an older version of this script cannot make a
# run report success while the command is still missing.

# Directory of this script. POSIX only: BASH_SOURCE is not available under sh.
SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"

# Which shell is the user running?
detect_shell() {
    if [ -n "$SHELL" ]; then
        basename "$SHELL"
        return
    fi
    user_shell="$(getent passwd "$(id -un)" 2>/dev/null | cut -d: -f7)"
    if [ -n "$user_shell" ]; then
        basename "$user_shell"
        return
    fi
    echo sh
}

DETECTED_SHELL="$(detect_shell)"
if [ -z "$DETECTED_SHELL" ]; then
    DETECTED_SHELL="sh"
fi

# The config file and the syntax for that shell — and only that shell.
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

# Exact-line match. A substring hit elsewhere in the file is not a PATH entry,
# and even a correct line in somebody else's config is not this shell's PATH.
if [ -f "$CONFIG_FILE" ] && grep -Fqx "$LINE" "$CONFIG_FILE" 2>/dev/null; then
    echo "Folder $SCRIPT_DIR is already configured in $CONFIG_FILE"
else
    mkdir -p "$(dirname -- "$CONFIG_FILE")"
    printf '%s\n' "$LINE" >> "$CONFIG_FILE"
    echo "Folder $SCRIPT_DIR added to PATH in $CONFIG_FILE"
fi

# Name the entries that belong to shells the user is not running. They are why
# a machine can look configured while the command is still missing.
for other in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile" "$HOME/.cshrc" \
             "$HOME/.config/fish/config.fish"; do
    if [ "$other" != "$CONFIG_FILE" ] && [ -f "$other" ] &&
       grep -Fq "$SCRIPT_DIR" "$other" 2>/dev/null; then
        echo "Note: $other also mentions this folder, but $DETECTED_SHELL does not read it."
    fi
done

# The config file only takes effect in the next shell. Say what to run.
case ":$PATH:" in
    *":$SCRIPT_DIR:"*)
        echo "It is active in this shell already."
        ;;
    *)
        echo "Run '$SOURCE_CMD' or open a new terminal to use it now."
        ;;
esac
