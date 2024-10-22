MKTOOLS_VERSION := 0.1.0

# Version check function
version_check:
	@echo "mktools version $(MKTOOLS_VERSION)"
	@curl -s https://api.github.com/repos/amenophis1er/mktools/releases/latest | \
		grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/' | \
		xargs -I {} sh -c 'if [ "$(MKTOOLS_VERSION)" != "{}" ]; then \
			echo "New version {} available! Visit https://github.com/amenophis1er/mktools/releases"; \
		fi'