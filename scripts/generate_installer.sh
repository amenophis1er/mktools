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

# Copy template and replace placeholders
cp "${PROJECT_ROOT}/templates/installer.sh.template" "${OUTPUT_DIR}/install.sh"

# Update file content placeholders
sed -i.bak \
    -e "s|__COLORS_MK__|$(encode_file "${PROJECT_ROOT}/src/common/colors.mk")|g" \
    -e "s|__VERSION_MK__|$(encode_file "${PROJECT_ROOT}/src/common/version.mk")|g" \
    -e "s|__DUMP_MK__|$(encode_file "${PROJECT_ROOT}/src/targets/dump/dump.mk")|g" \
    -e "s|__VARS_MK__|$(encode_file "${PROJECT_ROOT}/src/targets/dump/vars.mk")|g" \
    -e "s|__TARGETS_MK__|$(encode_file "${PROJECT_ROOT}/src/common/targets.mk")|g" \
    "${OUTPUT_DIR}/install.sh"

# Cleanup backup file
rm -f "${OUTPUT_DIR}/install.sh.bak"

# Make installer executable
chmod +x "${OUTPUT_DIR}/install.sh"

echo "Installer script generated at ${OUTPUT_DIR}/install.sh"