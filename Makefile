# ==============================================================================
# go-stats · Native Linux System Monitor & Systemd Daemon
# ==============================================================================
# Author: davidalvarezp (dap) <mail@davidalvarezp.com>
# Website: https://davidalvarezp.com/go-stats
# Repository: https://github.com/davidalvarezp/go-stats
# License: MIT
# ==============================================================================

APP_NAME  := go-stats
VERSION   := 1.0.0
PORT      := 10590
BUILD_DIR := bin
CMD_SRC   := cmd/$(APP_NAME)/main.go

# Compiler flags for lightweight static Linux binary
LDFLAGS   := -s -w -X 'main.Version=$(VERSION)' -X 'main.Author=David Alvarez Pampillon (dap)'

.PHONY: all build build-prod test lint install uninstall run clean help

# Default target
all: build-prod

help:
	@echo "go-stats v$(VERSION) build system"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Compile standard binary in $(BUILD_DIR)/$(APP_NAME)"
	@echo "  build-prod   Compile stripped production binary with optimizations"
	@echo "  test         Run all unit tests across internal packages"
	@echo "  lint         Run go vet static analysis"
	@echo "  install      Install binary to /usr/local/bin and enable systemd service"
	@echo "  uninstall    Stop and remove systemd service and binary"
	@echo "  run          Execute binary in foreground on port $(PORT)"
	@echo "  clean        Remove build artifacts and temporary files"

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_SRC)
	@echo "Built: $(BUILD_DIR)/$(APP_NAME)"

build-prod:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_SRC)
	@echo "Production binary compiled: $(BUILD_DIR)/$(APP_NAME) ($$(ls -lh $(BUILD_DIR)/$(APP_NAME) | awk '{print $$5}'))"

test:
	CGO_ENABLED=0 go test -v ./...

lint:
	CGO_ENABLED=0 go vet ./...

install: build-prod
	@chmod +x install.sh
	@sudo ./install.sh

uninstall:
	@chmod +x uninstall.sh
	@sudo ./uninstall.sh

run: build-prod
	./$(BUILD_DIR)/$(APP_NAME) -port $(PORT)

clean:
	rm -rf $(BUILD_DIR) *.log
	@echo "Cleaned build artifacts."
