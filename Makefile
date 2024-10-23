.PHONY: release help

# At the start of your main Makefile
HELP_TEXT := Available targets:
HELP_TEXT += "\n  release          - Create a new release/tag and push it (e.g., make release VERSION=1.0.0)"

.PHONY: help
help:
	@echo "$(HELP_TEXT)"

release:
	@if [ -z "$(VERSION)" ]; then echo "Please specify the VERSION, e.g., make release VERSION=1.0.0" && exit 1; fi
	@if git rev-parse "v$(VERSION)" >/dev/null 2>&1; then echo "Version v$(VERSION) already exists!" && exit 1; fi
	@echo "Creating new release for version $(VERSION)"
	@git tag -a "v$(VERSION)" -m "Release version $(VERSION)"
	@git push origin "v$(VERSION)"
	@echo "Release v$(VERSION) pushed to origin."
