#!/bin/sh
# Install crower and cr from a published release.
#
#   curl -fsSL https://raw.githubusercontent.com/fiwon123/crower/main/scripts/install.sh | sh
#   sh scripts/install.sh --version v0.8.0 --dir "$HOME/.local/bin"
#
# The archive is verified against the release's checksums.txt before anything
# is installed. Linux and macOS only; Windows uses the zip and add_to_path.bat.
set -eu

REPO="fiwon123/crower"
VERSION=""
INSTALL_DIR="${HOME}/.local/bin"
VERIFY=1

usage() {
	cat <<'TXT'
usage: install.sh [--version vX.Y.Z] [--dir DIR] [--no-verify]

  --version   install this tag instead of the newest release
  --dir       install into DIR (default: $HOME/.local/bin)
  --no-verify skip the checksums.txt check (not recommended)
TXT
}

while [ "$#" -gt 0 ]; do
	case "$1" in
	--version) VERSION="$2"; shift 2 ;;
	--dir) INSTALL_DIR="$2"; shift 2 ;;
	--no-verify) VERIFY=0; shift ;;
	-h|--help) usage; exit 0 ;;
	*) echo "install.sh: unknown argument $1" >&2; usage >&2; exit 2 ;;
	esac
done

os=$(uname -s)
case "$os" in
Linux) os=linux ;;
Darwin) os=darwin ;;
*) echo "install.sh: unsupported OS $os (use the Windows zip on Windows)" >&2; exit 1 ;;
esac

arch=$(uname -m)
case "$arch" in
x86_64|amd64) arch=amd64 ;;
aarch64|arm64) arch=arm64 ;;
*) echo "install.sh: unsupported architecture $arch" >&2; exit 1 ;;
esac

if [ -z "$VERSION" ]; then
	VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" |
		sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1)
fi

case "$VERSION" in
v[0-9]*.[0-9]*.[0-9]*) ;;
*) echo "install.sh: could not resolve a release version (got '$VERSION')" >&2; exit 1 ;;
esac

asset="crower_${VERSION}_${os}_${arch}.tar.gz"
base="https://github.com/${REPO}/releases/download/${VERSION}"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading ${asset}..."
curl -fsSL -o "$tmp/$asset" "$base/$asset"

if [ "$VERIFY" -eq 1 ]; then
	curl -fsSL -o "$tmp/checksums.txt" "$base/checksums.txt"

	line=$(grep " ${asset}\$" "$tmp/checksums.txt" || true)

	if [ -z "$line" ]; then
		echo "install.sh: checksums.txt has no entry for ${asset}" >&2
		exit 1
	fi

	if command -v sha256sum >/dev/null 2>&1; then
		printf '%s\n' "$line" | (cd "$tmp" && sha256sum -c -)
	else
		printf '%s\n' "$line" | (cd "$tmp" && shasum -a 256 -c -)
	fi
fi

tar -xzf "$tmp/$asset" -C "$tmp" crower cr

mkdir -p "$INSTALL_DIR"
for binary in crower cr; do
	cp "$tmp/$binary" "$INSTALL_DIR/$binary"
	chmod 0755 "$INSTALL_DIR/$binary"
done

echo "Installed crower and cr into ${INSTALL_DIR}"
"$INSTALL_DIR/crower" --version

case ":${PATH}:" in
*":${INSTALL_DIR}:"*) ;;
*)
	echo
	echo "Add ${INSTALL_DIR} to PATH:"
	echo "  crower init            # shows the line for your shell"
	echo "  crower init --apply    # adds it"
	;;
esac
