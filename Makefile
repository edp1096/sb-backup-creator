# SB Backup Creator Makefile

APP_NAME = sb-backup-creator
BUILD_FLAGS = -ldflags "-H windowsgui -s -w"

.PHONY: build clean run tidy

# Default target
build:
	go build $(BUILD_FLAGS) -o bin/$(APP_NAME).exe

# Clean build artifacts
clean:
	@if exist $(APP_NAME).exe del bin/$(APP_NAME).exe

# Build and run
run: build
	./$(APP_NAME).exe

# Update dependencies
tidy:
	go mod tidy

# Full setup (for first time)
setup:
	go mod init $(APP_NAME)
	go mod tidy
	make build