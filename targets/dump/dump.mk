# Include required variables
include $(dir $(lastword $(MAKEFILE_LIST)))/vars.mk
include $(mktools_path)/common/colors.mk

# Get the current folder name
CURRENT_FOLDER := $(notdir $(CURDIR))

# Help description
dump.help:
	@echo "  dump             - Create a context dump of the project structure"

# Dump project structure and file contents
.PHONY: dump
dump:
	@echo "$(CYAN)Generating file system listing...$(RESET)"
	@DUMP_FILE="$(CURRENT_FOLDER).context.txt"; \
	EXCLUDE_ARGS="$$(printf -- "--exclude=%s " $(EXCLUDE_PATHS))"; \
	EXCLUDE_ARGS="$$EXCLUDE_ARGS --exclude=$$DUMP_FILE --exclude=.git --exclude=.idea --exclude=Makefile"; \
	if [ -f .gitignore ]; then EXCLUDE_ARGS="$$EXCLUDE_ARGS --exclude-from=.gitignore"; fi; \
	> "$$DUMP_FILE" && \
	rsync -av --delete $$EXCLUDE_ARGS \
		--include=.git/ --include=.git/** --prune-empty-dirs . | \
		grep -v '^building file list' | grep -v '^done$$' | grep -v '/$$' | \
		grep -v 'sent [0-9]\+ bytes' | grep -v 'total size is' | grep -v 'speedup is' | \
		awk 'NF { if ($$1 ~ /d/) { print $$5 " (Directory)"; } else { print $$5 " (File)"; }}' >> "$$DUMP_FILE"
	@echo "\n=========================\nFile Content Dump\n=========================\n" >> "$(CURRENT_FOLDER).context.txt"
	@echo "$(CYAN)Appending file contents...$(RESET)"
	@DUMP_FILE="$(CURRENT_FOLDER).context.txt"; \
	EXCLUDE_ARGS="$$(printf -- "--exclude=%s " $(EXCLUDE_PATHS))"; \
	EXCLUDE_ARGS="$$EXCLUDE_ARGS --exclude=$$DUMP_FILE --exclude=.git --exclude=.idea --exclude=Makefile"; \
	if [ -f .gitignore ]; then EXCLUDE_ARGS="$$EXCLUDE_ARGS --exclude-from=.gitignore"; fi; \
	CONTENT_EXCLUDE_PATTERN="$$(printf ".*\\.%s$$|" $(CONTENT_EXCLUDE_EXT))"; \
	CONTENT_EXCLUDE_PATTERN="$${CONTENT_EXCLUDE_PATTERN%|}"; \
	rsync -av --delete $$EXCLUDE_ARGS \
		--include=.git/ --include=.git/** --prune-empty-dirs . | \
		grep -v '^building file list' | grep -v '^done$$' | grep -v '/$$' | \
		grep -v 'sent [0-9]\+ bytes' | grep -v 'total size is' | grep -v 'speedup is' | \
		awk -v exclude_pattern="$$CONTENT_EXCLUDE_PATTERN" ' \
		NF && $$1 !~ /d/ { \
			if ($$5 !~ exclude_pattern) { \
				print "\n" $$5 "\n```"; \
				system("cat \"" $$5 "\""); \
				print "\n```"; \
			} else { \
				print "\n" $$5 " (Content excluded)"; \
			} \
		}' >> "$$DUMP_FILE"
	@echo "$(GREEN)Context dump with file contents created at $(CURRENT_FOLDER).context.txt$(RESET)"