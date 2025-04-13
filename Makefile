.PHONY: build clean install

# Version information
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse HEAD)

# Build flags
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Build the CLI tool
build:
	go build $(LDFLAGS) -o bin/csrtool ./cmd/csrtool

# Install the CLI tool
install:
	go install $(LDFLAGS) ./cmd/csrtool

# Clean build artifacts
clean:
	rm -rf bin/

# Create build directory
bin/:
	mkdir -p bin/
