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


reset-tags: ## Delete all Git tags locally and remotely, and attempt to clear workflow runs if possible
	@echo "$(RED)WARNING: This will delete ALL tags locally and remotely, and attempt to clear ALL workflow runs$(NC)"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" != "y" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	@echo "Starting cleanup process..."
	@if ! git remote get-url origin >/dev/null 2>&1; then \
		echo "$(RED)ERROR: Remote 'origin' not configured$(NC)" && exit 1; \
	fi

	@# Check and install gh if needed
	@if ! command -v gh >/dev/null 2>&1; then \
		echo "$(YELLOW)GitHub CLI (gh) is not installed.$(NC)"; \
		if command -v brew >/dev/null 2>&1; then \
			read -p "Would you like to install it using Homebrew? [y/N] " install_confirm; \
			if [ "$$install_confirm" = "y" ]; then \
				echo "Installing GitHub CLI..."; \
				brew install gh || { echo "$(RED)Failed to install gh$(NC)"; exit 1; }; \
				echo "$(GREEN)GitHub CLI installed successfully$(NC)"; \
			else \
				echo "$(YELLOW)WARNING: Skipping workflow runs cleanup...$(NC)"; \
			fi; \
		else \
			echo "$(YELLOW)WARNING: Homebrew not found. Please install gh manually. Skipping workflow runs cleanup...$(NC)"; \
		fi; \
	fi

	@# Try to authenticate gh if installed
	@if command -v gh >/dev/null 2>&1; then \
		if ! gh auth status >/dev/null 2>&1; then \
			echo "GitHub CLI needs authentication."; \
			read -p "Would you like to authenticate now? [y/N] " auth_confirm; \
			if [ "$$auth_confirm" = "y" ]; then \
				gh auth login || { echo "$(RED)Authentication failed$(NC)"; exit 1; }; \
			else \
				echo "$(YELLOW)WARNING: Skipping workflow runs cleanup...$(NC)"; \
			fi; \
		fi; \
	fi

	@# Try to delete workflow runs if gh is available and authenticated
	@if command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then \
		echo "Deleting all workflow runs..."; \
		gh run list --limit 1000 --json databaseId -q '.[].databaseId' | while read -r run_id; do \
			gh run delete "$$run_id" || \
				echo "$(YELLOW)WARNING: Failed to delete run $$run_id$(NC)"; \
		done; \
		echo "$(GREEN)Workflow runs cleanup completed$(NC)"; \
	fi

	@echo "Deleting all local tags..."
	@git tag | xargs -r git tag -d
	@echo "Fetching remote updates and pruning tags..."
	@git fetch --prune origin "+refs/tags/*:refs/tags/*"
	@echo "Deleting all remote tags..."
	@git tag -l | xargs -r -I {} git push origin :refs/tags/{} || \
		(echo "$(RED)ERROR: Failed to delete remote tags$(NC)" && exit 1)
	@echo "$(GREEN)SUCCESS: All tags have been deleted locally and remotely.$(NC)"