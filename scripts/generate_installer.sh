#!/usr/bin/env bash

# Exit on error
set -e

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Output directory
OUTPUT_DIR="${PROJECT_ROOT}/dist"
mkdir -p "$OUTPUT_DIR"

# Function to encode file content to base64
encode_file() {
    local file=$1
    if [ -f "$file" ]; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            base64 < "$file"
        else
            base64 -w 0 < "$file"
        fi
    else
        echo "File not found: $file" >&2
        exit 1
    fi
}

# Function to convert string to uppercase
to_upper() {
    echo "$1" | tr '[:lower:]' '[:upper:]'
}

# Copy template
cp "${PROJECT_ROOT}/templates/installer.sh.template" "${OUTPUT_DIR}/install.sh"

# Start with common files
sed -i.bak \
    -e "s|__COLORS_MK__|$(encode_file "${PROJECT_ROOT}/src/common/colors.mk")|g" \
    "${OUTPUT_DIR}/install.sh"

# Create a temporary file for the processed installer
temp_installer=$(mktemp)
cp "${OUTPUT_DIR}/install.sh" "$temp_installer"

# Find all targets (directories in src/targets that contain __init__.mk)
for target_dir in "${PROJECT_ROOT}"/src/targets/*/; do
    if [ -f "${target_dir}__init__.mk" ]; then
        target=$(basename "$target_dir")
        target_upper=$(to_upper "$target")

        # Replace placeholder with target variables
        sed -i.bak -e "/__TARGETS_VARS__/i\\
# Embedded content for target: $target\\
${target_upper}_INIT=\"__${target_upper}_INIT__\"\\
${target_upper}_VARS=\"__${target_upper}_VARS__\"\\
" "$temp_installer"

        # Replace content placeholders
        sed -i.bak \
            -e "s|__${target_upper}_INIT__|$(encode_file "${PROJECT_ROOT}/src/targets/${target}/__init__.mk")|g" \
            -e "s|__${target_upper}_VARS__|$(encode_file "${PROJECT_ROOT}/src/targets/${target}/__vars__.mk")|g" \
            "$temp_installer"
    fi
done

# Remove the placeholder line
sed -i.bak "/__TARGETS_VARS__/d" "$temp_installer"

# Move the processed file to final location
mv "$temp_installer" "${OUTPUT_DIR}/install.sh"

# Cleanup backup files
rm -f "${OUTPUT_DIR}/install.sh.bak" "${temp_installer}.bak"

# Make installer executable
chmod +x "${OUTPUT_DIR}/install.sh"

echo "Installer script generated at ${OUTPUT_DIR}/install.sh"