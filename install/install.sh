#!/usr/bin/env bash

VERSION=${1:-latest}
INSTALL_DIR="$HOME/.local/share/mktools"
BIN_DIR="$HOME/.local/bin"
REPO_URL="https://github.com/amenophis1er/mktools"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

# Error handling
set -e

echo -e "${CYAN}Installing mktools...${NC}"

# Create directories if they don't exist
echo -e "${CYAN}Creating directories...${NC}"
mkdir -p "$INSTALL_DIR" || { echo -e "${RED}Failed to create $INSTALL_DIR${NC}"; exit 1; }
mkdir -p "$BIN_DIR" || { echo -e "${RED}Failed to create $BIN_DIR${NC}"; exit 1; }

# Clone repository
echo -e "${CYAN}Cloning repository...${NC}"
rm -rf "$INSTALL_DIR.tmp"
git clone "$REPO_URL" "$INSTALL_DIR.tmp" || {
    echo -e "${RED}Failed to clone repository${NC}"
    exit 1
}

# Copy files to installation directory
echo -e "${CYAN}Installing files...${NC}"
rsync -a --delete "$INSTALL_DIR.tmp/" "$INSTALL_DIR/" || {
    echo -e "${RED}Failed to copy files${NC}"
    rm -rf "$INSTALL_DIR.tmp"
    exit 1
}
rm -rf "$INSTALL_DIR.tmp"

# Install version file
echo "$VERSION" > "$INSTALL_DIR/VERSION"

# Create the mktools command
MKTOOLS_CMD="$BIN_DIR/mktools"

# Create the script
echo -e "${CYAN}Creating mktools command...${NC}"
cat > "$MKTOOLS_CMD" << 'EOF' || { echo -e "${RED}Failed to create mktools command${NC}"; exit 1; }
#!/usr/bin/env bash

MKTOOLS_DIR="$HOME/.local/share/mktools"

function list_targets() {
    echo -e "${CYAN}Available targets:${NC}"
    echo -e "\n${YELLOW}dump${NC}"
    echo "  Description: Creates a context dump of your project structure"
    echo "  Features:"
    echo "    - Directory structure listing"
    echo "    - File listing"
    echo "    - Content of text files (excludes binary files)"
    echo -e "  ${GREEN}Usage:${NC}"
    echo "    1. Install:  mktools install dump"
    echo "    2. Use:      make dump"
    echo
    echo -e "${CYAN}Installation Instructions:${NC}"
    echo "  To use any target, you need to install it first:"
    echo "    mktools install <target-name>"
    echo
    echo "  This will update your Makefile to include the target."
    echo "  After installation, use 'make <target-name>' to run the target."
    echo
    echo -e "${YELLOW}Note:${NC} Targets are used with 'make' command after installation,"
    echo "      not with 'mktools' command directly."
}

function check_version() {
    local current_version=$(cat "$MKTOOLS_DIR/VERSION" 2>/dev/null || echo "unknown")

    # Don't check for updates if we just installed the latest version
    if [ "$current_version" = "latest" ]; then
        echo "Current version: latest (freshly installed)"
        return 0
    fi

    local latest_version=$(curl -s https://api.github.com/repos/amenophis1er/mktools/releases/latest | \
        grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')

    echo "Current version: $current_version"
    if [ -n "$latest_version" ] && [ "$current_version" != "$latest_version" ]; then
        echo "New version $latest_version available!"
        echo "Run: curl -sSL https://raw.githubusercontent.com/amenophis1er/mktools/main/install/install.sh | bash"
    fi
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

    # Create temp file
    temp_file=$(mktemp)

    # Preserve any comments at the start of the file
    sed -n '/^[[:space:]]*#/p' "$makefile" > "$temp_file"

    # Add PHONY declarations
    grep "^\.PHONY:" "$makefile" >> "$temp_file"

    # Add mktools includes
    echo "" >> "$temp_file"
    echo "# Added by mktools" >> "$temp_file"
    echo "mktools_path := $MKTOOLS_DIR" >> "$temp_file"
    echo "include \$(mktools_path)/common/colors.mk" >> "$temp_file"
    echo "include \$(mktools_path)/targets/$target/*.mk" >> "$temp_file"
    echo "" >> "$temp_file"

    # Add the default target if it exists (preserving its position)
    grep "^default:" "$makefile" >> "$temp_file"

    # Add the rest of the file, excluding what we've already added
    grep -v "^\.PHONY:" "$makefile" | \
    grep -v "^default:" | \
    grep -v "^[[:space:]]*#" >> "$temp_file"

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
    "version")
        check_version
        ;;
    "dump"|"test"|*)  # First try to execute as target, fallback to usage
        if [ -z "$1" ]; then
            echo -e "${YELLOW}Usage: mktools [list|install <target-name>|version]${NC}"
            exit 1
        fi
        # If it looks like they're trying to use a target directly
        if [ -d "$MKTOOLS_DIR/targets/$1" ]; then
            echo -e "${RED}Note: Targets cannot be run directly with mktools${NC}"
            echo -e "${GREEN}To use the '$1' target:${NC}"
            echo "  1. First install it:  mktools install $1"
            echo "  2. Then use it with:  make $1"
            exit 1
        else
            echo -e "${YELLOW}Usage: mktools [list|install <target-name>|version]${NC}"
            exit 1
        fi
        ;;
esac
EOF

# Make the command executable
chmod +x "$MKTOOLS_CMD" || { echo -e "${RED}Failed to make mktools command executable${NC}"; exit 1; }

echo -e "${GREEN}mktools installed successfully!${NC}"
echo "You can now use 'mktools list' to see available targets"
echo "and 'mktools install <target-name>' to install a target"