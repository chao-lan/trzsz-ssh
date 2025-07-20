#!/bin/bash

# tssh zsh completion installation script

echo "Installing tssh zsh completion..."

# Create completion directory if it doesn't exist
COMPLETION_DIR="$HOME/.zsh/completions"
mkdir -p "$COMPLETION_DIR"

# Copy completion file
cp _tssh_completion "$COMPLETION_DIR/_tssh"

# Add to .zshrc if not already present
ZSHRC="$HOME/.zshrc"
FPATH_LINE="fpath=(\$HOME/.zsh/completions \$fpath)"
AUTOLOAD_LINE="autoload -U compinit && compinit"

if ! grep -q "fpath=.*\.zsh/completions" "$ZSHRC"; then
    echo "" >> "$ZSHRC"
    echo "# tssh completion" >> "$ZSHRC"
    echo "$FPATH_LINE" >> "$ZSHRC"
fi

if ! grep -q "autoload -U compinit" "$ZSHRC"; then
    echo "$AUTOLOAD_LINE" >> "$ZSHRC"
fi

# Since you have alias ssh=tssh, we need to make sure completion works for both
echo "" >> "$ZSHRC"
echo "# Make completion work for ssh alias (which points to tssh)" >> "$ZSHRC"
echo "compdef _tssh ssh" >> "$ZSHRC"

echo "Installation complete!"
echo ""
echo "Please run the following command to reload your zsh configuration:"
echo "source ~/.zshrc"
echo ""
echo "Or restart your terminal."
echo ""
echo "After that, you should be able to use tab completion with:"
echo "  ssh --de<TAB>  (will complete to --debug)"
echo "  ssh --up<TAB>  (will complete to --upload-file)"
echo "  ssh hostname<TAB>  (will complete SSH hosts from ~/.ssh/config)"