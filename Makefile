APP_NAME := crower
BUILD_DIR := build

WINDOWS_BIN := $(BUILD_DIR)/$(APP_NAME).exe
LINUX_BIN := $(BUILD_DIR)/$(APP_NAME)

.PHONY: all windows linux clean

all: windows linux

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

windows: $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_BIN)

linux: $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_BIN)

clean:
	rm -rf $(BUILD_DIR)
