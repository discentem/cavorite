#!/usr/bin/env bash
# Install a cavorite (or localstore plugin) tarball built by
# build-cavorite-package.sh / build-localstore-plugin-package.sh.

set -euo pipefail

PACKAGE_PATH=""
DEST_DIR="/usr/local/bin"

usage() {
    echo "Usage: $0 --pkg PATH [--dest DIR]"
    echo ""
    echo "Options:"
    echo "  --pkg PATH   Path to the .tar.gz package to install (required)"
    echo "  --dest DIR   Directory to install the binary into (default: /usr/local/bin)"
    echo "  --help       Show this help message"
}

while [[ $# -gt 0 ]]; do
    case $1 in
        --pkg)
            PACKAGE_PATH="$2"
            shift 2
            ;;
        --dest)
            DEST_DIR="$2"
            shift 2
            ;;
        --help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

if [ -z "$PACKAGE_PATH" ]; then
    echo "Error: --pkg is required" 1>&2
    usage
    exit 1
fi

if [ ! -f "$PACKAGE_PATH" ]; then
    echo "Error: package not found at $PACKAGE_PATH" 1>&2
    exit 1
fi

STAGE_DIR="$(mktemp -d)"
trap 'rm -rf "$STAGE_DIR"' EXIT

tar -C "$STAGE_DIR" -xzf "$PACKAGE_PATH"

BIN_NAME=""
for candidate in cavorite localstore; do
    if [ -f "$STAGE_DIR/$candidate" ]; then
        BIN_NAME="$candidate"
        break
    fi
done

if [ -z "$BIN_NAME" ]; then
    echo "Error: no known binary (cavorite, localstore) found in package" 1>&2
    exit 1
fi

echo "Installing ${BIN_NAME} to ${DEST_DIR}..."
install -d "$DEST_DIR"
install -m 755 "$STAGE_DIR/$BIN_NAME" "$DEST_DIR/$BIN_NAME"

echo ""
echo "Installed: ${DEST_DIR}/${BIN_NAME}"
