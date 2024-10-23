# Terminal colors
RESET := \033[0m
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m

# Common variables
CURRENT_DIR := $(shell pwd)
CURRENT_FILE := $(lastword $(MAKEFILE_LIST))
CURRENT_TARGET = $(firstword $(MAKECMDGOALS))

# Function to extract metadata from a target file
# Usage: $(call get_metadata,file,key)
define get_metadata
$(shell sed -n 's/^# \(.*\):\s*\(.*\)/\1|\2/p' $(1) | grep "^$(2)|" | cut -d'|' -f2)
endef

# Function to check if a target is active
# Usage: $(call is_target_active,file)
define is_target_active
$(shell if [ "$$(sed -n 's/^# active:\s*\(.*\)/\1/p' $(1))" = "true" ]; then echo "true"; else echo "false"; fi)
endef

# Function to extract help text from a target file
# Usage: $(call get_help_text,file)
define get_help_text
$(shell sed -n 's/^# help:\s*\(.*\)/\1/p' $(1))
endef

# Function to display help for mktools targets
define print_mktools_help
	@for f in $(sort $(wildcard src/targets/*/__init__.mk)); do \
		if [ "$$($(call is_target_active,$$f))" = "true" ]; then \
			target=$$(basename $$(dirname $$f)); \
			help_text=$$($(call get_help_text,$$f)); \
			if [ -n "$$help_text" ]; then \
				printf "  %-15s - %s\n" "$$target" "$$help_text"; \
			fi; \
		fi; \
	done
endef