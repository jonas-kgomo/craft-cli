#!/bin/bash

# CRAFT CLI - Ultra-Lite Version
# No Bubbletea, just beautiful Lipgloss styling

echo "🛠️  Starting CRAFT CLI (Ultra-Lite)"
echo ""

export CLICOLOR_FORCE=1
export COLORTERM=truecolor
export PATH="/opt/homebrew/bin:$PATH"

./craft "$@"
