#!/usr/bin/env bash

# Exit on error
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Test directories
TEST_DIR="${PROJECT_ROOT}/test-output/clean-install"
TEST_DIR_IDEMPOTENT="${PROJECT_ROOT}/test-output/idempotent"
TEST_DIR_EXISTING="${PROJECT_ROOT}/test-output/existing-makefile"
TEST_DIR_SELECTIVE="${PROJECT_ROOT}/test-output/selective-install"

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

# Helper function to automate installation
automate_install() {
    local dir=$1
    local choice=$2

    echo -e "${YELLOW}Running automated installation:${NC}"
    echo -e "Directory: $dir"
    echo -e "Choice: $choice"

    # Create answers file
    echo "y" > "${PROJECT_ROOT}/test-output/answers"
    echo "$choice" >> "${PROJECT_ROOT}/test-output/answers"

    cd "$dir"
    ${PROJECT_ROOT}/dist/install.sh < "${PROJECT_ROOT}/test-output/answers"
}

# Main test sequence
main() {
    echo -e "${YELLOW}Starting automated test sequence...${NC}"

    # Clean start
    cleanup
    mkdir -p "$TEST_DIR" "$TEST_DIR_IDEMPOTENT" "$TEST_DIR_EXISTING" "$TEST_DIR_SELECTIVE"

    # Test 1: Generate installer
    run_test "Installer Generation" "
        cd '$PROJECT_ROOT' && \
        bash scripts/generate_installer.sh && \
        test -f dist/install.sh
    "

    # Test 2: Clean installation (all targets)
    run_test "Clean Installation" "
        automate_install '$TEST_DIR' '0' && \
        cd '$TEST_DIR' && \
        make -n dump && \
        make -n version
    "

    # Test 3: Idempotency
    run_test "Idempotency" "
        automate_install '$TEST_DIR_IDEMPOTENT' '0' && \
        cd '$TEST_DIR_IDEMPOTENT' && \
        cp Makefile Makefile.first && \
        automate_install '$TEST_DIR_IDEMPOTENT' '0' && \
        diff Makefile Makefile.first
    "

    # Test 4: Existing Makefile
    run_test "Existing Makefile" "
        cd '$TEST_DIR_EXISTING' && \
        echo -e '.PHONY: test\n\ntest:\n\techo \"test\"' > Makefile && \
        automate_install '$TEST_DIR_EXISTING' '0' && \
        make -n test && \
        make -n dump && \
        make -n version && \
        [ \$(grep -c 'test:' Makefile) -eq 1 ]
    "

    # Test 5: Selective Installation (dump only)
    run_test "Selective Installation" "
        automate_install '$TEST_DIR_SELECTIVE' '1' && \
        cd '$TEST_DIR_SELECTIVE' && \
        make -n dump && \
        ! make -n version
    "

    # Cleanup
    cleanup

    echo -e "\n${GREEN}All automated tests passed successfully!${NC}"
}

# Run main function
main "$@"