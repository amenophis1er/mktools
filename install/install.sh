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

function check_version() {
    local current_version=$(cat "$MKTOOLS_DIR/VERSION" 2>/dev/null || echo "unknown")
    local latest_version=$(curl -s https://api.github.com/repos/amenophis1er/mktools/releases/latest | \
        grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')

    echo "Current version: $current_version"
    if [ "$current_version" != "$latest_version" ]; then
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
    "version")
        check_version
        ;;
    "dump"|"test"|*)  # First try to execute as target, fallback to usage
        if [ -z "$1" ]; then
            echo "Usage: mktools [list|install <target-name>|version|<target-name>]"
            exit 1
        fi
        # Check if this is an installed target
        if [ -f "Makefile" ] && grep -q "include.*$1" "Makefile"; then
            make "$1"
        else
            echo "Usage: mktools [list|install <target-name>|version|<target-name>]"
            exit 1
        fi
        ;;
esac
EOF
then
    echo -e "${RED}Failed to create mktools command${NC}"
    exit 1
fi