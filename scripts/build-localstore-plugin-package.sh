#!/usr/bin/env bash
# Build a release tarball for the localstore file-store plugin via bazel.
#
# Usage: scripts/build-localstore-plugin-package.sh [VERSION]
#
# Produces dist/cavorite-plugin-localstore_<VERSION>_<OS>_<ARCH>.tar.gz

set -euo pipefail

check_exit_code() {
    if [ "$1" != "0" ]; then
        echo "$2: $1" 1>&2
        exit 1
    fi
}

if ! command -v bazel &> /dev/null; then
    echo "Error: bazel is not installed"
    echo ""
    echo "Install bazelisk (recommended):"
    echo "  brew install bazelisk"
    exit 1
fi

TOOL="localstore"
PACKAGE_PREFIX="cavorite-plugin-localstore"
VERSION="${1:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

case "$(uname -s)" in
    Linux*)  OS="linux" ;;
    Darwin*) OS="darwin" ;;
    *)       echo "Error: unsupported OS $(uname -s)" 1>&2; exit 1 ;;
esac

case "$(uname -m)" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Error: unsupported arch $(uname -m)" 1>&2; exit 1 ;;
esac

echo "Building ${PACKAGE_PREFIX} ${VERSION} for ${OS}/${ARCH}..."

bazel build //plugins/localstore:localstore --stamp --embed_label="${VERSION}"
check_exit_code "$?" "Error building ${TOOL} plugin"

BIN_PATH="$(bazel cquery //plugins/localstore:localstore --output=files 2>/dev/null)"
if [ ! -f "$BIN_PATH" ]; then
    check_exit_code 1 "Failed to find built plugin binary at $BIN_PATH"
fi

DIST_DIR="$PROJECT_ROOT/dist"
STAGE_DIR="$(mktemp -d)"
trap 'rm -rf "$STAGE_DIR"' EXIT

mkdir -p "$DIST_DIR"
cp "$BIN_PATH" "$STAGE_DIR/${TOOL}"
chmod 755 "$STAGE_DIR/${TOOL}"

# Bazel-built (rules_go, fastbuild) binaries on macOS can be missing an
# LC_UUID load command, which newer dyld refuses to load. Ad-hoc codesigning
# regenerates the load commands and fixes this.
if [ "$OS" = "darwin" ] && command -v codesign &> /dev/null; then
    codesign -f -s - "$STAGE_DIR/${TOOL}" 2>/dev/null || true
fi

PACKAGE_NAME="${PACKAGE_PREFIX}_${VERSION}_${OS}_${ARCH}.tar.gz"
tar -C "$STAGE_DIR" -czf "$DIST_DIR/$PACKAGE_NAME" "${TOOL}"
check_exit_code "$?" "Error creating tarball"

echo ""
echo "Package created: dist/${PACKAGE_NAME}"
