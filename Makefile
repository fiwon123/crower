APP_NAME := crower
APP_NAME_SHORT := cr
GUI_NAME := crower-gui
GUI_PACKAGE = $(BUILD_DIR)/gui-package
BUILD_DIR := build
SCRIPTS_DIR := scripts
SCRIPT_NAME := add_to_path
INJECT_VERSION:= github.com/fiwon123/crower/cmd.Version

# The Fyne window is a GUI-subsystem program on Windows. Without -H windowsgui
# the loader gives it a console window, so double-clicking the executable opens
# a terminal instead of the window. The icon, manifest and version are a .syso
# in the package directory — the only hook Go offers for a Windows resource.
GUI_LDFLAGS = -X $(INJECT_VERSION)=$(VERSION)

ifeq ($(GOOS),windows)
GUI_LDFLAGS += -H windowsgui
endif

# Build tools run on the host even when GOOS points at the release target: a
# GOOS=windows `go run` would produce a Windows binary the Linux runner cannot
# execute. CGO is off for them too, so the cross compiler in CC is not asked to
# build the host's runtime/cgo.
HOST_GO := GOOS=$$(go env GOHOSTOS) GOARCH=$$(go env GOHOSTARCH) CGO_ENABLED=0

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
	shell opencode opencode-auto release-platform release-gui-platform package-cli package-gui verify-release

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

# Platform release for a CLI-only archive (local use).
#
# The release workflow builds one combined archive per platform with
# release-gui-platform; this target stays for a local CLI-only build.
#
# The binary inside the archive is always named $(APP_NAME)$(EXT) /
# $(APP_NAME_SHORT)$(EXT) — never arch-suffixed. The platform and architecture
# are already in the archive filename, and the documentation tells users to run
# `crower` and `cr`, so a suffixed binary is unreachable by the name it ships.
# verify-release enforces that.
release-platform: $(BUILD_DIR)
	mkdir -p $(DIST_DIR)
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH) with version $(VERSION)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags "-X $(INJECT_VERSION)=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)$(EXT)
	cp $(BUILD_DIR)/$(APP_NAME)$(EXT) $(BUILD_DIR)/$(APP_NAME_SHORT)$(EXT)
	@$(MAKE) package-cli GOOS=$(GOOS) GOARCH=$(GOARCH)
	@$(MAKE) verify-release

# Platform release for a platform that carries the Fyne GUI.
#
# Fyne needs CGO, so this target runs on the target OS runner (Linux, Windows
# or macOS), not as a cross-compile. It builds the cgo-free `crower`/`cr` pair
# and `crower-gui`, then packages all three into the one platform archive.
release-gui-platform: $(BUILD_DIR)
	mkdir -p $(DIST_DIR) $(GUI_PACKAGE)
	@echo "Building $(APP_NAME), $(APP_NAME_SHORT) and $(GUI_NAME) for $(GOOS)/$(GOARCH) with version $(VERSION)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags "-X $(INJECT_VERSION)=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)$(EXT)
	cp $(BUILD_DIR)/$(APP_NAME)$(EXT) $(BUILD_DIR)/$(APP_NAME_SHORT)$(EXT)
ifeq ($(GOOS),windows)
	@$(HOST_GO) go run ./internal/tools/winres -icon internal/fyneui/logo.png \
		-out cmd/crowergui/rsrc_windows_$(GOARCH).syso -arch $(GOARCH) \
		-name $(GUI_NAME) -version "$(VERSION)"
endif
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -tags fyne -ldflags "$(GUI_LDFLAGS)" -o $(BUILD_DIR)/$(GUI_NAME)$(EXT) ./cmd/crowergui
ifeq ($(GOOS),windows)
	@rm -f cmd/crowergui/rsrc_windows_$(GOARCH).syso
	@$(HOST_GO) go run ./internal/tools/winres -verify $(BUILD_DIR)/$(GUI_NAME).exe
	@$(HOST_GO) go run ./internal/tools/winres -verify-console $(BUILD_DIR)/$(APP_NAME).exe
	@$(HOST_GO) go run ./internal/tools/winres -verify-console $(BUILD_DIR)/$(APP_NAME_SHORT).exe
	@# The mingw cross-link uses msvcrt, the C runtime every Windows has; a
	@# native build would import the UCRT API sets instead, which is the heap of
	@# small DLLs a machine can be missing. Keep the exe on msvcrt (#468).
	@if command -v x86_64-w64-mingw32-objdump >/dev/null 2>&1; then \
		if x86_64-w64-mingw32-objdump -p $(BUILD_DIR)/$(GUI_NAME).exe | grep -q 'api-ms-win-crt'; then \
			echo "ERROR: $(GUI_NAME).exe imports the UCRT API sets; the cross-link must stay on msvcrt (#468)" >&2; \
			exit 1; \
		fi; \
		if ! x86_64-w64-mingw32-objdump -p $(BUILD_DIR)/$(GUI_NAME).exe | grep -q 'msvcrt.dll'; then \
			echo "ERROR: $(GUI_NAME).exe does not import msvcrt.dll; the cross-link changed runtime (#468)" >&2; \
			exit 1; \
		fi; \
	fi
endif
	@$(MAKE) package-gui GOOS=$(GOOS) GOARCH=$(GOARCH)
	@$(MAKE) verify-release GUI=1

# Fail the build if an archive does not carry the names users are told to run.
# GOOS/GOARCH are pinned to the host: this target inspects archives, so it must
# run on the machine doing the build rather than be cross-compiled to the
# archive's target. GUI=1 additionally requires crower-gui in every archive.
ARCHIVEVERIFY_FLAGS := $(if $(filter 1,$(GUI)),-gui,)

verify-release:
	@GOOS=$$(go env GOHOSTOS) GOARCH=$$(go env GOHOSTARCH) CGO_ENABLED=0 go run ./internal/tools/archiveverify $(ARCHIVEVERIFY_FLAGS) $(DIST_DIR)

ifeq ($(GOOS),windows)
EXT := .exe
SCRIPT_EXT := .bat
CLI_ARCHIVE := $(DIST_DIR)/$(APP_NAME)_$(VERSION)_$(GOOS)_$(GOARCH).zip
else
EXT :=
SCRIPT_EXT := .sh
CLI_ARCHIVE := $(DIST_DIR)/$(APP_NAME)_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz
endif

package-cli: $(BUILD_DIR)
	@mkdir -p $(DIST_DIR)
ifeq ($(GOOS),windows)
	zip -j $(CLI_ARCHIVE) $(BUILD_DIR)/$(APP_NAME)$(EXT) $(BUILD_DIR)/$(APP_NAME_SHORT)$(EXT) $(SCRIPTS_DIR)/add_to_path$(SCRIPT_EXT) README.md LICENSE LICENSE-APACHE NOTICE
else
	tar -czvf $(CLI_ARCHIVE) -C $(BUILD_DIR) $(APP_NAME)$(EXT) $(APP_NAME_SHORT)$(EXT) -C ../$(SCRIPTS_DIR) add_to_path$(SCRIPT_EXT) ../README.md ../LICENSE ../LICENSE-APACHE ../NOTICE
endif
	@echo "CLI package: $(CLI_ARCHIVE)"

package-gui: $(BUILD_DIR)
	@mkdir -p $(DIST_DIR) $(GUI_PACKAGE)
	cp $(BUILD_DIR)/$(APP_NAME)$(EXT) $(BUILD_DIR)/$(APP_NAME_SHORT)$(EXT) $(BUILD_DIR)/$(GUI_NAME)$(EXT) \
		$(SCRIPTS_DIR)/$(SCRIPT_NAME)$(SCRIPT_EXT) README.md NOTICE LICENSE LICENSE-APACHE LICENSE-FYNE $(GUI_PACKAGE)/
ifeq ($(GOOS),windows)
	cp $(SCRIPTS_DIR)/WINDOWS.txt $(GUI_PACKAGE)/
	cd $(GUI_PACKAGE) && zip -j $(abspath $(CLI_ARCHIVE)) *
else
	tar -czvf $(CLI_ARCHIVE) -C $(GUI_PACKAGE) $(APP_NAME)$(EXT) $(APP_NAME_SHORT)$(EXT) $(GUI_NAME)$(EXT) $(SCRIPT_NAME)$(SCRIPT_EXT) README.md NOTICE LICENSE LICENSE-APACHE LICENSE-FYNE
endif
	@echo "Platform package: $(CLI_ARCHIVE)"


# ---------------------------------------------------------------------------
# Local verification gate
#
# gofmt, vet, test, build — in that order, because a formatting failure hides
# everything after it.
#
# Two properties this target deliberately has:
#
#   - A missing toolchain FAILS the target. It does not skip. A gate that
#     cannot fail is worse than no gate: it reads green and proves nothing.
#   - After the run it names the packages that have no test files. `go test`
#     prints "[no test files]" for those and still exits 0, so a green run
#     would otherwise look like coverage that does not exist.
#
# The two cross builds exist because several packages choose their behaviour
# with build tags (`handler_windows.go`, `handler_darwin.go`, ...). A host build
# on Linux never compiles those files, so a Windows-only break would reach a
# release: CI builds Linux on every pull request, and the release matrix is the
# first job that compiles the other five targets. These builds are cheap and
# they can fail.
#
# Tests use isolated temp directories (internal/helper/tests.go), so parallel
# runs are not racing over a shared fixture (#47) — a red run is a regression
# until proven otherwise, not a baseline to discount.
# ---------------------------------------------------------------------------
# set -e matters: without it the recipe's status is the last command's (the
# note), so a failing vet, test or build still exited 0 — a gate that could not
# fail. The `|| true` keeps an empty untested-package list from tripping it.
check: ## gofmt, vet, test, build — the local verification gate
	@set -e; \
	if ! command -v go >/dev/null 2>&1; then \
		echo "ERROR: go is not installed — the check cannot run, so nothing is verified" >&2; \
		exit 1; \
	fi; \
	out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then echo "unformatted files:"; echo "$$out"; exit 1; fi; \
	go vet ./...; \
	go test ./...; \
	go build ./...; \
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ./...; \
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build ./...; \
	notested="$$(go list -f '{{if and (not .TestGoFiles) (not .XTestGoFiles)}}{{.ImportPath}}{{end}}' ./... | grep -v '^$$' || true)"; \
	if [ -n "$$notested" ]; then \
		echo ""; \
		echo "NOTE: these packages have no test files — this run does not cover them:"; \
		echo "$$notested" | sed 's/^/  /'; \
	fi

# ---------------------------------------------------------------------------
# Dev sandbox
#
# A pinned toolchain and a shell that behave the same on every machine. The
# sandbox is optional: the project still builds and tests with a plain host
# Go toolchain (see the `check` target above).
#
# Cross-platform note: this Makefile targets POSIX hosts — Linux, macOS and
# WSL2 — where `make` runs recipes under a POSIX shell. Native Windows users
# should use `python scripts/dev-up.py dev-up` (or `scripts/dev-up.ps1`), which
# provides the same targets without a POSIX shell. Everyone can run the project
# with a plain `docker compose up`, which needs no Make at all.
# ---------------------------------------------------------------------------

# Extra flags for the opencode invocation. `make opencode-auto` /
# `make dev-opencode-auto` set this to --auto, which auto-approves permissions
# that are not explicitly denied; the secret-file deny rules still apply.
OPENCODE_ARGS ?=

# --- Host identity -----------------------------------------------------------
# Run bind-mounted dev services as the host developer instead of root. Derived
# from the workspace OWNER, not the current process, so invoking make inside an
# already-running bind-mounted container still sends the host identity.
# GNU stat (-c) first, then BSD/macOS stat (-f), then plain `id`. The `id`
# fallback is what makes this work everywhere `make` does.
HOST_UID ?= $(shell stat -c '%u' . 2>/dev/null || stat -f '%u' . 2>/dev/null || id -u)
HOST_GID ?= $(shell stat -c '%g' . 2>/dev/null || stat -f '%g' . 2>/dev/null || id -g)
HOST_HOME ?= $(HOME)
export HOST_UID HOST_GID HOST_HOME

# --- Optional in-container agent ---------------------------------------------
# Starting the sandbox never requires opencode. The agent layer is optional so
# `make dev-up` works on macOS and Windows too, where the host binary cannot
# execute inside the Linux container.
#
# NOTE: not named OPENCODE — the opencode runtime exports an OPENCODE env var
# (=1), which would override a `?=` default through make's env import.
OPENCODE_BIN ?= $(HOME)/.opencode/bin/opencode

# The dev container is Linux, so only an ELF binary can execute inside it. A
# macOS Mach-O or Windows PE install is detected here and skipped, instead of
# being mounted and failing later with "exec format error".
OPENCODE_MAGIC := $(shell head -c 4 "$(OPENCODE_BIN)" 2>/dev/null | od -An -tx1 2>/dev/null | tr -d ' \n')
AGENT_AVAILABLE := $(if $(filter 7f454c46,$(OPENCODE_MAGIC)),yes,no)

ifeq ($(AGENT_AVAILABLE),yes)
COMPOSE_UP := docker compose -f docker-compose.yml -f docker-compose.agent.yml
else
COMPOSE_UP := docker compose
endif

help: ## List the available targets
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@printf '\n  \033[36m%-30s\033[0m %s\n' "docker compose up" "run the sandbox with no host prerequisites (any OS)"
	@printf '  \033[36m%-30s\033[0m %s\n' "scripts/dev-up.py" "native Windows: python scripts/dev-up.py dev-up"

# --- Preflight ---------------------------------------------------------------
# Refuses to start as root, and pre-creates the read-only mount sources the
# agent overlay needs. Missing mount sources are fatal on Windows and create
# root-owned directories on Linux, so they are created here rather than assumed.
# This target NEVER fails on a missing or non-Linux agent binary: that is what
# `require-agent` is for, and it is depended on only by the agent targets.
preflight: check-identity prepare-agent-mounts
	@if [ "$(AGENT_AVAILABLE)" = "yes" ]; then \
	  echo "[sandbox] in-container agent available — merging docker-compose.agent.yml."; \
	else \
	  echo "[sandbox] in-container agent unavailable (no Linux binary at $(OPENCODE_BIN))." >&2; \
	  echo "[sandbox]   Starting the sandbox without it." >&2; \
	fi

check-identity:
	@test "$(HOST_UID)" != "0" || { echo "ERROR: refusing to run the dev sandbox as UID 0; repair workspace ownership or set HOST_UID to the host user." >&2; exit 1; }

prepare-agent-mounts:
	@if [ "$(AGENT_AVAILABLE)" = "yes" ]; then \
	  mkdir -p "$(HOST_HOME)/.config/opencode" "$(HOST_HOME)/.config/gh" && touch "$(HOST_HOME)/.gitconfig"; fi

# Only the agent-specific targets need the in-container binary.
require-agent:
	@test "$(AGENT_AVAILABLE)" = "yes" || { \
	  echo "ERROR: the in-container AI agent needs a Linux opencode binary at $(OPENCODE_BIN)." >&2; \
	  echo "  A macOS/Windows binary cannot execute inside the Linux dev container." >&2; \
	  echo "  Recommended: use WSL2, where the host binary is Linux." >&2; \
	  echo "  Alternative: run 'opencode' on the host against the running stack." >&2; \
	  exit 1; }

# The host owns the opencode binary; the dev sandbox mounts it read-only
# (`~/.opencode/bin/opencode -> /usr/local/bin/opencode:ro`), so it can only be
# upgraded here, on the host. A *file* bind mount pins the inode, so a running
# sandbox keeps the old binary until it restarts — `opencode-upgrade` restarts
# it for you. Override the install URL if your installer differs.
OPENCODE_INSTALL_URL ?= https://opencode.ai/install

opencode-version: ## Print the host opencode version
	@command -v opencode >/dev/null 2>&1 || { echo "ERROR: opencode not on PATH — install it from $(OPENCODE_INSTALL_URL)" >&2; exit 1; }
	@opencode --version

opencode-upgrade: ## Upgrade host opencode and refresh the dev sandbox
	@command -v opencode >/dev/null 2>&1 || { echo "ERROR: opencode not on PATH — install it from $(OPENCODE_INSTALL_URL)" >&2; exit 1; }
	@echo "opencode before: $$(opencode --version 2>/dev/null || echo unknown)"
	@opencode upgrade </dev/null || { echo "opencode upgrade unavailable — falling back to the install script"; curl -fsSL $(OPENCODE_INSTALL_URL) | bash; }
	@echo "opencode after:  $$(opencode --version 2>/dev/null || echo unknown)"
	@if [ -n "$$(docker compose ps -q dev 2>/dev/null)" ]; then \
	  echo "restarting the dev sandbox so the bind-mounted binary refreshes..."; \
	  docker compose restart dev; \
	else \
	  echo "dev sandbox not running — nothing to refresh"; \
	fi

# --- Sandbox lifecycle -------------------------------------------------------
# HOST_UID/HOST_GID need no explicit export here: docker-compose.yml defaults
# them, and the exported values above supply the host's real identity when make
# is running on the host. Only the GitHub token is resolved at up-time.
dev-up: preflight ## Start the dev sandbox (foreground; Ctrl+C stops it)
	@GH_TOKEN=$$(scripts/resolve-gh-token.sh 2>/dev/null) GITHUB_TOKEN=$$(scripts/resolve-gh-token.sh 2>/dev/null) $(COMPOSE_UP) up --build dev

dev-down: ## Stop the dev sandbox (keeps volumes)
	@docker compose down

dev-log: ## Tail dev sandbox logs
	@docker compose logs -f dev

dev-build: check-identity ## Rebuild the dev image after Dockerfile/dependency changes
	@docker compose build dev

reset: ## Stop everything and wipe volumes (clean slate)
	@docker compose down -v

shell: dev-shell
opencode: dev-opencode
opencode-auto: dev-opencode-auto

dev-shell: preflight ## Open a shell inside the dev sandbox
	@$(COMPOSE_UP) exec dev bash || $(COMPOSE_UP) run --rm dev bash

dev-opencode: preflight require-agent ## Run the AI coding agent inside the dev sandbox
	@if [ -z "$$(docker compose ps -q dev 2>/dev/null)" ]; then echo "[opencode] starting dev stack..."; GH_TOKEN=$$(scripts/resolve-gh-token.sh 2>/dev/null) GITHUB_TOKEN=$$(scripts/resolve-gh-token.sh 2>/dev/null) $(COMPOSE_UP) up -d dev; fi
	@if [ -t 0 ]; then $(COMPOSE_UP) exec -it dev zsh -lc "cd /workspace && opencode $(OPENCODE_ARGS)" || $(COMPOSE_UP) exec -it dev bash -lc "opencode $(OPENCODE_ARGS)"; else $(COMPOSE_UP) exec -T dev zsh -lc "cd /workspace && opencode $(OPENCODE_ARGS)" || $(COMPOSE_UP) exec -T dev bash -lc "opencode $(OPENCODE_ARGS)"; fi

# Same recipe, with --auto: auto-approve permissions that are not explicitly
# denied (every secret-file deny rule still applies).
dev-opencode-auto: OPENCODE_ARGS = --auto
dev-opencode-auto: dev-opencode
