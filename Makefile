.PHONY: release reset-tags
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: VERSION is required. Use: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "ERROR: Working directory not clean. Commit or stash changes before releasing."; \
		exit 1; \
	fi
	@git tag -a $(VERSION) -m "Release $(VERSION)" || { echo "ERROR: Failed to create Git tag."; exit 1; }
	@git push origin $(VERSION) || { echo "ERROR: Failed to push Git tag to origin."; exit 1; }
	@echo "SUCCESS: Release $(VERSION) tagged and pushed. GitHub Actions will handle the release."

reset-tags:
	@echo "Deleting all local tags..."
	@git tag | xargs -r git tag -d
	@echo "Fetching remote updates and pruning tags..."
	@git fetch --prune origin "+refs/tags/*:refs/tags/*"
	@echo "Deleting all remote tags..."
	@git tag -l | xargs -r -I {} git push origin :refs/tags/{}
	@echo "SUCCESS: All tags have been deleted locally and remotely."
