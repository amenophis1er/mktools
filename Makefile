# Makefile

.PHONY: release help

default: help

release:
	@if [ -z "$(VERSION)" ]; then echo "Please specify the VERSION, e.g., make release VERSION=1.0.0" && exit 1; fi
	@if git rev-parse "v$(VERSION)" >/dev/null 2>&1; then echo "Version v$(VERSION) already exists!" && exit 1; fi
	@echo "Creating new release for version $(VERSION)"
	@git tag -a "v$(VERSION)" -m "Release version $(VERSION)"
	@git push origin "v$(VERSION)"
	@echo "Release v$(VERSION) pushed to origin."

help:
	@echo "Available targets:"
	@echo "  - release          Create a new release/tag and push it (e.g., make release VERSION=1.0.0)"