.DEFAULT_GOAL := help

# Colors for terminal output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color
BOLD := \033[1m

.PHONY: help release reset-tags

help: ## Display this help message
	@echo "$(BOLD)Usage:$(NC)"
	@echo "  make $(CYAN)<target>$(NC) $(YELLOW)[OPTIONS]$(NC)"
	@echo ""
	@echo "$(BOLD)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)Examples:$(NC)"
	@echo "  make $(CYAN)release$(NC) $(YELLOW)VERSION=v1.0.0$(NC)		Create and push a new release tag"
	@echo "  make $(CYAN)reset-tags$(NC)                  	Delete all Git tags locally and remotely"

release: ## Create and push a new release tag (requires VERSION=v*.*.*)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)ERROR: VERSION is required. Use: make release VERSION=v1.0.0$(NC)"; \
		exit 1; \
	fi
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)ERROR: Working directory not clean. Commit or stash changes before releasing.$(NC)"; \
		exit 1; \
	fi
	@git tag -a $(VERSION) -m "Release $(VERSION)" || { echo "$(RED)ERROR: Failed to create Git tag.$(NC)"; exit 1; }
	@git push origin $(VERSION) || { echo "$(RED)ERROR: Failed to push Git tag to origin.$(NC)"; exit 1; }
	@echo "$(GREEN)SUCCESS: Release $(VERSION) tagged and pushed. GitHub Actions will handle the release.$(NC)"

reset-tags: ## Delete all Git tags locally and remotely (requires confirmation)
	@echo "$(RED)WARNING: This will delete ALL tags locally and remotely$(NC)"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" != "y" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	@echo "Starting tag cleanup process..."
	@if ! git remote get-url origin >/dev/null 2>&1; then \
		echo "$(RED)ERROR: Remote 'origin' not configured$(NC)" && exit 1; \
	fi
	@echo "Deleting all local tags..."
	@git tag | xargs -r git tag -d
	@echo "Fetching remote updates and pruning tags..."
	@git fetch --prune origin "+refs/tags/*:refs/tags/*"
	@echo "Deleting all remote tags..."
	@git tag -l | xargs -r -I {} git push origin :refs/tags/{} || \
		(echo "$(RED)ERROR: Failed to delete remote tags$(NC)" && exit 1)
	@echo "$(GREEN)SUCCESS: All tags have been deleted locally and remotely.$(NC)"