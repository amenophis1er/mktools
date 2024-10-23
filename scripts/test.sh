#!/usr/bin/env bash

# Exit on error
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Test directories
TEST_DIR="${PROJECT_ROOT}/test-output/clean-install"
TEST_DIR_IDEMPOTENT="${PROJECT_ROOT}/test-output/idempotent"
TEST_DIR_EXISTING="${PROJECT_ROOT}/test-output/existing-makefile"

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}Cleaning up test directories...${NC}"
    rm -rf "${PROJECT_ROOT}/test-output"
}

# Error handler
handle_error() {
    echo -e "\n${RED}Error: Test failed${NC}"
    echo "See above for details"
    cleanup
    exit 1
}

# Set up error handling
trap handle_error ERR

# Helper function to run tests
run_test() {
    local test_name=$1
    local test_cmd=$2

    echo -e "\n${YELLOW}Running test: ${test_name}${NC}"
    echo -e "${CYAN}Executing automated test command:${NC} $test_cmd"
    eval "$test_cmd"
    echo -e "${GREEN}✓ Test passed: ${test_name}${NC}"
}

# Helper function to verify installed files
verify_installation() {
    local dir=$1

    # Check common files
    [ -f "$dir/src/common/colors.mk" ] || (echo "Missing colors.mk" && return 1)

    # Check each target from the source project
    for target_dir in "${PROJECT_ROOT}/src/targets/*/"; do
        if [ -f "${target_dir}__init__.mk" ] && [ -f "${target_dir}__vars__.mk" ]; then
            target=$(basename "$target_dir")
            echo "Checking target: $target"
            [ -f "$dir/src/targets/$target/__init__.mk" ] || (echo "Missing $target/__init__.mk" && return 1)
            [ -f "$dir/src/targets/$target/__vars__.mk" ] || (echo "Missing $target/__vars__.mk" && return 1)
        fi
    done

    # Check if Makefile exists and contains help target
    [ -f "$dir/Makefile" ] || (echo "Missing Makefile" && return 1)
    grep -q "^help:" "$dir/Makefile" || (echo "Missing help target in Makefile" && return 1)

    return 0
}

# Helper function to automate installation
automate_install() {
    local dir=$1
    cd "$dir"
    echo -e "${YELLOW}Running automated installation in: $dir${NC}"
    echo "y" | ${PROJECT_ROOT}/dist/install.sh
}

# Main test sequence
main() {
    echo -e "${YELLOW}Starting automated test sequence...${NC}"

    # Clean start
    cleanup
    mkdir -p "$TEST_DIR" "$TEST_DIR_IDEMPOTENT" "$TEST_DIR_EXISTING"

    # Test 1: Generate installer
    run_test "Installer Generation" "
        cd '$PROJECT_ROOT' && \
        bash scripts/generate_installer.sh && \
        test -f dist/install.sh
    "

    # Test 2: Clean installation
    run_test "Clean Installation" "
        cd '$TEST_DIR' && \
        automate_install '$TEST_DIR' && \
        verify_installation '$TEST_DIR' && \
        make help
    "

    # Test 3: Idempotency
    run_test "Idempotency" "
        cd '$TEST_DIR_IDEMPOTENT' && \
        automate_install '$TEST_DIR_IDEMPOTENT' && \
        cp Makefile Makefile.first && \
        automate_install '$TEST_DIR_IDEMPOTENT' && \
        diff Makefile Makefile.first && \
        verify_installation '$TEST_DIR_IDEMPOTENT'
    "

    # Test 4: Existing Makefile
    run_test "Existing Makefile" "
        cd '$TEST_DIR_EXISTING' && \
        echo -e '.PHONY: test\n\ntest:\n\techo \"test\"' > Makefile && \
        automate_install '$TEST_DIR_EXISTING' && \
        verify_installation '$TEST_DIR_EXISTING' && \
        make test
    "

    # Cleanup
    cleanup

    echo -e "\n${GREEN}All automated tests passed successfully!${NC}"
}

# Run main function
main "$@"