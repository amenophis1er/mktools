MKTOOLS_VERSION := 0.1.0

# Only run version check if explicitly called
ifeq ($(MAKECMDGOALS),version_check)
  # This will only show when version_check is explicitly called
  $(info mktools version $(MKTOOLS_VERSION))
  # Perform version check only when explicitly called
  $(shell if [ "$$(cat $(mktools_path)/VERSION)" != "latest" ]; then \
    curl -s https://api.github.com/repos/amenophis1er/mktools/releases/latest | \
    grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/' | \
    xargs -I {} sh -c 'if [ "$(MKTOOLS_VERSION)" != "{}" ]; then \
      echo "New version {} available! Visit https://github.com/amenophis1er/mktools/releases"; \
    fi'; \
  fi)
endif

# Silent target that does nothing unless explicitly called
.PHONY: version_check
version_check:
	@true