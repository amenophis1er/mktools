# Target definitions and dependencies
AVAILABLE_TARGETS := dump version

# Target groups
BASIC_TARGETS := version
DEVELOPMENT_TARGETS := dump

# Target dependencies (target:dependency1,dependency2)
TARGET_DEPS := \
    dump:colors \
    version:colors

# Target descriptions (for help command)
define TARGET_HELP
dump:Create a context dump of the project structure
version:Display version information
endef

# Target files mapping
define TARGET_FILES
dump:src/targets/dump/vars.mk src/targets/dump/dump.mk
version:src/common/version.mk
colors:src/common/colors.mk
endef

# Export as variables
export AVAILABLE_TARGETS
export TARGET_DEPS
export TARGET_HELP
export TARGET_FILES