APP_NAME := crower
BUILD_DIR := build

# Detect last tag and increment patch
LAST_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
VERSION := $(shell echo $(LAST_TAG) | sed 's/^v//')
MAJOR := $(shell echo $(VERSION) | cut -d. -f1)
MINOR := $(shell echo $(VERSION) | cut -d. -f2)
PATCH := $(shell echo $(VERSION) | cut -d. -f3)
NEW_PATCH := $(shell echo $$(( $(PATCH) + 1 )))
NEW_TAG := v$(MAJOR).$(MINOR).$(NEW_PATCH)

WINDOWS_BIN := $(BUILD_DIR)/$(APP_NAME).exe
LINUX_BIN := $(BUILD_DIR)/$(APP_NAME)

WINDOWS_ZIP := $(BUILD_DIR)/$(APP_NAME)_$(NEW_TAG)_windows.zip
LINUX_TAR := $(BUILD_DIR)/$(APP_NAME)_$(NEW_TAG)_linux.tar.gz

.PHONY: all windows linux zip_linux zip_windows clean release

all: release

# Create build folder if missing
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build Windows binary
windows: $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_BIN)

# Build Linux binary
linux: $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_BIN)

# Compress Windows binary
zip_windows: windows
	zip -j $(WINDOWS_ZIP) $(WINDOWS_BIN) README.md LICENSE

# Compress Linux binary
zip_linux: linux
	tar -czvf $(LINUX_TAR) -C $(BUILD_DIR) $(notdir $(LINUX_BIN)) ../README.md ../LICENSE

# Clean build folder
clean:
	rm -rf $(BUILD_DIR)

# Build & compress everything
release: zip_windows zip_linux
	@echo "Release ready: $(NEW_TAG)"
