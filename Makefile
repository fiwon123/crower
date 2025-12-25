APP_NAME := crower
APP_NAME_SHORT := cr
BUILD_DIR := build

# Detect last tag and increment patch
VERSION := $(shell git describe --tags --dirty --always)

WINDOWS_BIN := $(BUILD_DIR)/$(APP_NAME).exe
SHORT_WINDOWS_BIN := $(BUILD_DIR)/$(APP_NAME_SHORT).exe

LINUX_BIN := $(BUILD_DIR)/$(APP_NAME)
SHORT_LINUX_BIN := $(BUILD_DIR)/$(APP_NAME_SHORT)

WINDOWS_ZIP := $(BUILD_DIR)/$(APP_NAME)_$(VERSION)_windows.zip
LINUX_TAR := $(BUILD_DIR)/$(APP_NAME)_$(VERSION)_linux.tar.gz

.PHONY: all windows linux zip_linux zip_windows clean release

all: release

# Create build folder if missing
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build Windows binary
windows: $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/fiwon123/crower/cmd.Version=$(VERSION)" -o $(WINDOWS_BIN)
	cp $(WINDOWS_BIN) $(SHORT_WINDOWS_BIN)

# Build Linux binary
linux: $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/fiwon123/crower/cmd.Version=$(VERSION)" -o  $(LINUX_BIN)
	cp $(LINUX_BIN) $(SHORT_LINUX_BIN)

# Compress Windows binary
zip_windows: windows
	zip -j $(WINDOWS_ZIP) $(WINDOWS_BIN) $(SHORT_WINDOWS_BIN) README.md LICENSE

# Compress Linux binary
zip_linux: linux
	tar -czvf $(LINUX_TAR) -C $(BUILD_DIR) $(notdir $(LINUX_BIN)) $(notdir $(SHORT_LINUX_BIN)) ../README.md ../LICENSE

# Clean build folder
clean:
	rm -rf $(BUILD_DIR)

# Build & compress everything
release: zip_windows zip_linux
	@echo "Release ready: $(VERSION)"
