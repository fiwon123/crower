APP_NAME := crower
APP_NAME_SHORT := cr
BUILD_DIR := build
SCRIPTS_DIR := scripts
SCRIPT_NAME := add_to_path
INJECT_VERSION:= github.com/fiwon123/crower/cmd.Version

VERSION ?= $(shell git describe --tags --always)
DIST_DIR := $(BUILD_DIR)/dist

WINDOWS_BIN := $(BUILD_DIR)/$(APP_NAME).exe
SHORT_WINDOWS_BIN := $(BUILD_DIR)/$(APP_NAME_SHORT).exe
WINDOWS_ZIP := $(BUILD_DIR)/$(APP_NAME)_$(VERSION)_windows.zip
WINDOWS_SCRIPT :=  $(SCRIPTS_DIR)/$(SCRIPT_NAME).bat

LINUX_BIN := $(BUILD_DIR)/$(APP_NAME)
SHORT_LINUX_BIN := $(BUILD_DIR)/$(APP_NAME_SHORT)
LINUX_TAR := $(BUILD_DIR)/$(APP_NAME)_$(VERSION)_linux.tar.gz
LINUX_SCRIPT :=  $(SCRIPTS_DIR)/$(SCRIPT_NAME).sh


.PHONY: all windows linux zip_linux zip_windows clean release check \
	help dev-up dev-down dev-shell dev-opencode dev-opencode-auto dev-log dev-build reset \
	preflight require-agent prepare-agent-mounts check-identity opencode-version opencode-upgrade \
	shell opencode opencode-auto release-platform package-cli package-gui

all: release

# Create build folder if missing
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)
	mkdir -p $(DIST_DIR)

# Build Windows binary
windows: $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X $(INJECT_VERSION)=$(VERSION)" -o $(WINDOWS_BIN)
	cp $(WINDOWS_BIN) $(SHORT_WINDOWS_BIN)

# Build Linux binary
linux: $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X $(INJECT_VERSION)=$(VERSION)" -o $(LINUX_BIN)
	cp $(LINUX_BIN) $(SHORT_LINUX_BIN)

# Compress Windows binary
zip_windows: windows
	zip -j $(WINDOWS_ZIP) $(WINDOWS_BIN) $(SHORT_WINDOWS_BIN) $(WINDOWS_SCRIPT) README.md LICENSE LICENSE-APACHE NOTICE

# Compress Linux binary
zip_linux: linux
	tar -czvf $(LINUX_TAR) \
	          -C $(BUILD_DIR) $(notdir $(LINUX_BIN)) $(notdir $(SHORT_LINUX_BIN)) \
	          -C ../$(SCRIPTS_DIR) $(notdir $(LINUX_SCRIPT)) ../README.md ../LICENSE ../LICENSE-APACHE ../NOTICE

# Clean build folder
clean:
	rm -rf $(BUILD_DIR)

# Build & compress everything
release: zip_windows zip_linux
	@echo "Release ready: $(VERSION)"

# Platform release for CI (cross-compile)
release-platform: $(BUILD_DIR)
	mkdir -p $(DIST_DIR)
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH) with version $(VERSION)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags "-X $(INJECT_VERSION)=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)_$(GOOS)_$(GOARCH)$(EXT)
	cp $(BUILD_DIR)/$(APP_NAME)_$(GOOS)_$(GOARCH)$(EXT) $(BUILD_DIR)/$(APP_NAME_SHORT)_$(GOOS)_$(GOARCH)$(EXT)
	@$(MAKE) package-cli GOOS=$(GOOS) GOARCH=$(GOARCH)
	@$(MAKE) package-gui GOOS=$(GOOS) GOARCH=$(GOARCH)

ifeq ($(GOOS),windows)
EXT := .exe
SCRIPT_EXT := .bat
CLI_ARCHIVE := $(DIST_DIR)/$(APP_NAME)_$(VERSION)_$(GOOS)_$(GOARCH).zip
GUI_ARCHIVE := $(DIST_DIR)/$(APP_NAME)_gui_$(VERSION)_$(GOOS)_$(GOARCH).zip
LEGACY_CLI_ARCHIVE := $(DIST_DIR)/$(APP_NAME)_$(VERSION)_$(GOOS).zip
else
EXT :=
SCRIPT_EXT := .sh
CLI_ARCHIVE := $(DIST_DIR)/$(APP_NAME)_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz
GUI_ARCHIVE := $(DIST_DIR)/$(APP_NAME)_gui_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz
LEGACY_CLI_ARCHIVE := $(DIST_DIR)/$(APP_NAME)_$(VERSION)_$(GOOS).tar.gz
endif

package-cli: $(BUILD_DIR)
	@mkdir -p $(DIST_DIR)
ifeq ($(GOOS),windows)
	zip -j $(CLI_ARCHIVE) $(BUILD_DIR)/$(APP_NAME)_$(GOOS)_$(GOARCH)$(EXT) $(BUILD_DIR)/$(APP_NAME_SHORT)_$(GOOS)_$(GOARCH)$(EXT) $(SCRIPTS_DIR)/add_to_path$(SCRIPT_EXT) README.md LICENSE LICENSE-APACHE NOTICE
ifeq ($(GOARCH),amd64)
	cp $(CLI_ARCHIVE) $(LEGACY_CLI_ARCHIVE)
endif
else
	tar -czvf $(CLI_ARCHIVE) -C $(BUILD_DIR) $(APP_NAME)_$(GOOS)_$(GOARCH) $(APP_NAME_SHORT)_$(GOOS)_$(GOARCH) -C ../$(SCRIPTS_DIR) add_to_path$(SCRIPT_EXT) ../README.md ../LICENSE ../LICENSE-APACHE ../NOTICE
ifeq ($(GOARCH),amd64)
	cp $(CLI_ARCHIVE) $(LEGACY_CLI_ARCHIVE)
endif
endif
	@echo "CLI package: $(CLI_ARCHIVE)"

package-gui: $(BUILD_DIR)
	@mkdir -p $(DIST_DIR)
ifeq ($(GOOS),windows)
	zip -j $(GUI_ARCHIVE) $(BUILD_DIR)/$(APP_NAME)_$(GOOS)_$(GOARCH)$(EXT) $(BUILD_DIR)/$(APP_NAME_SHORT)_$(GOOS)_$(GOARCH)$(EXT) README.md NOTICE LICENSE LICENSE-APACHE
else
	tar -czvf $(GUI_ARCHIVE) -C $(BUILD_DIR) $(APP_NAME)_$(GOOS)_$(GOARCH) $(APP_NAME_SHORT)_$(GOOS)_$(GOARCH) ../README.md ../NOTICE ../LICENSE ../LICENSE-APACHE
endif
	@echo "GUI package: $(GUI_ARCHIVE)"
