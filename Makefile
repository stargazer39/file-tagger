# Variables
VERSION=$(shell git rev-parse --short HEAD)
REMOTE=origin
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Check for uncommitted changes
check_changes:
	@if ! git diff-index --quiet HEAD --; then \
		echo "There are uncommitted changes. Please commit or stash them first."; \
		exit 1; \
	fi

# Check if the current commit is already tagged
check_tag:
	@if git describe --exact-match --tags $(VERSION) >/dev/null 2>&1; then \
		echo "The current commit is already tagged."; \
		exit 1; \
	fi

# Create a Git tag
tag: check_changes check_tag
	@echo "Creating Git tag..."
	git tag -a $(VERSION) -m "Release version $(VERSION)"
	git push $(REMOTE) $(VERSION)

# Push the changes
push: check_changes
	@echo "Pushing changes to the remote repository..."
	git push $(REMOTE) $(BRANCH)

# Composite target to tag and push
release: tag push
	@echo "Release process completed."
