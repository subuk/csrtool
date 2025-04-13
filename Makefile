.PHONY: build clean install wasm web-build web

# Version information
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse HEAD)

# Build flags
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# WASM exec path
GOROOT ?= $(shell go env GOROOT)
WASM_EXEC = $(GOROOT)/lib/wasm/wasm_exec.js

# Build the CLI tool
build: bin/
	go build $(LDFLAGS) -o bin/csrtool ./cmd/csrtool

# Install the CLI tool
install:
	go install $(LDFLAGS) ./cmd/csrtool

# Build WebAssembly module
wasm: web/public/
	GOOS=js GOARCH=wasm go build -o web/public/csrtool.wasm ./pkg/wasm
	cp "$(WASM_EXEC)" web/public/

# Build web application
web-build: wasm
	cd web && npm install && npm run build

# Run web development server
web: wasm
	cd web && npm install && npm run dev

# Clean build artifacts and generated files
clean:
	rm -rf bin/ web/public/csrtool.wasm web/public/wasm_exec.js web/dist/
	rm -f private.key csr.pem

# Create build directory
bin/:
	mkdir -p bin/

# Create web/public directory
web/public/:
	mkdir -p web/public/
