# Version information will be embedded during installer generation
MKTOOLS_VERSION := 0.1.0

.PHONY: version
version:
	@echo "mktools version $(MKTOOLS_VERSION)"
	@if command -v curl >/dev/null 2>&1; then \
		latest_version=$$(curl -s https://api.github.com/repos/amenophis1er/mktools/releases/latest | \
		grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/'); \
		if [ "$(MKTOOLS_VERSION)" != "$$latest_version" ]; then \
			echo "New version $$latest_version available! Visit https://github.com/amenophis1er/mktools/releases"; \
		fi \
	fi

# Update help text
HELP_TEXT := Available targets:
HELP_TEXT += "\n  version         - Display version information"