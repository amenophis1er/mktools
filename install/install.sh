#!/usr/bin/env bash

# Installation directory
INSTALL_DIR="$HOME/.local/share/mktools"
BIN_DIR="$HOME/.local/bin"  # Change to user's bin directory
REPO_URL="git@github.com:amenophis1er/mktools.git"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Error handling
set -e  # Exit on error

# Create directories if they don't exist
mkdir -p "$INSTALL_DIR"
mkdir -p "$BIN_DIR"

# Clone or update the repository
if [ -d "$INSTALL_DIR/.git" ]; then
    echo -e "${CYAN}Updating mktools...${NC}"
    cd "$INSTALL_DIR" && git pull
else
    echo -e "${CYAN}Installing mktools...${NC}"
    if ! git clone "$REPO_URL" "$INSTALL_DIR"; then
        echo -e "${RED}Failed to clone repository. Make sure you have access to $REPO_URL${NC}"
        exit 1
    fi
fi

# Create the mktools command
MKTOOLS_CMD="$BIN_DIR/mktools"

# Create the script
if ! cat > "$MKTOOLS_CMD" << 'EOF'
#!/usr/bin/env bash

MKTOOLS_DIR="$HOME/.local/share/mktools"

function list_targets() {
    echo "Available targets:"
    for target in "$MKTOOLS_DIR/targets/"*/; do
        if [ -d "$target" ]; then
            target_name=$(basename "$target")
            echo "  $target_name"
        fi
    done
}

function install_target() {
    local target=$1
    local makefile="Makefile"
    local target_dir="$MKTOOLS_DIR/targets/$target"

    # Check if target exists
    if [ ! -d "$target_dir" ]; then
        echo -e "${RED}Target $target not found${NC}"
        return 1
    fi

    # Create Makefile if it doesn't exist
    if [ ! -f "$makefile" ]; then
        touch "$makefile"
    fi

    # Check if target already included
    if grep -q "include.*$target" "$makefile"; then
        echo -e "${YELLOW}Target $target already included${NC}"
        return 0
    fi

    # Add includes at the top of the Makefile
    temp_file=$(mktemp)
    echo "# Added by mktools" > "$temp_file"
    echo "mktools_path := $MKTOOLS_DIR" >> "$temp_file"
    echo "include \$(mktools_path)/common/*.mk" >> "$temp_file"
    echo "include \$(mktools_path)/targets/$target/*.mk" >> "$temp_file"
    echo "" >> "$temp_file"
    cat "$makefile" >> "$temp_file"
    mv "$temp_file" "$makefile"

    echo -e "${GREEN}Target $target installed successfully${NC}"
}

case "$1" in
    "list")
        list_targets
        ;;
    "install")
        if [ -z "$2" ]; then
            echo "Please specify a target to install"
            echo "Usage: mktools install <target-name>"
            exit 1
        fi
        install_target "$2"
        ;;
    *)
        echo "Usage: mktools [list|install <target-name>]"
        ;;
esac
EOF
then
    echo -e "${RED}Failed to create mktools command${NC}"
    exit 1
fi

# Make the command executable
if ! chmod +x "$MKTOOLS_CMD"; then
    echo -e "${RED}Failed to make mktools command executable${NC}"
    exit 1
fi

# Add BIN_DIR to PATH if not already there
if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
    echo -e "${CYAN}Adding $BIN_DIR to your PATH...${NC}"
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zshrc"
    echo -e "${YELLOW}Please restart your shell or run: export PATH=\"$BIN_DIR:\$PATH\"${NC}"
fi

echo -e "${GREEN}mktools installed successfully!${NC}"
echo "You can now use 'mktools list' to see available targets"
echo "and 'mktools install <target-name>' to install a target"