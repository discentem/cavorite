#!/bin/bash
# Download and install rustfs (S3-compatible server)
# GitHub: https://github.com/rustfs/rustfs
# Pinned to: v1.0.0-beta.12, gnu-libc builds
#
# Supports macOS (aarch64) and Linux (x86_64, aarch64).

set -e

RUSTFS_VERSION="1.0.0-beta.12"
INSTALL_DIR="/usr/local/bin"
INSTALL=false

# Pinned zip + extracted-binary SHA256 hashes per platform, taken from the
# release's SHA256SUMS file and verified locally against the extracted binary.
platform_key() {
    local os arch
    os="$(uname -s)"
    arch="$(uname -m)"
    case "$os" in
        Darwin)
            case "$arch" in
                arm64) echo "macos-aarch64" ;;
                *) echo "" ;;
            esac
            ;;
        Linux)
            case "$arch" in
                x86_64) echo "linux-x86_64-gnu" ;;
                aarch64|arm64) echo "linux-aarch64-gnu" ;;
                *) echo "" ;;
            esac
            ;;
        *) echo "" ;;
    esac
}

zip_sha256_for() {
    case "$1" in
        macos-aarch64)     echo "f5266eda245fa4dab5acf28bef7bbab6c1da7f3e9575ddc7db803894107e09f5" ;;
        linux-x86_64-gnu)  echo "9b9a17fb006acd7ae2bcb8227ba4ba10b81d7cafe08081701643951fdca57fb8" ;;
        linux-aarch64-gnu) echo "a6615b98973d7dbaf6308b2b39faf1a7387b1eefcc4b96932d3bfe51b291a121" ;;
        *) echo "" ;;
    esac
}

binary_sha256_for() {
    case "$1" in
        macos-aarch64)     echo "0f9dedc7c606fe133ed33cc27464cc64705544e4e1360673f4c5e8a78e931bce" ;;
        linux-x86_64-gnu)  echo "8d1d132d3a509efa9dda3b285e714be3283c1091855119ab6a0d0848d72aabf0" ;;
        linux-aarch64-gnu) echo "9e5d8155287abc90124abe6a3a505468f0c1fc9d0f96466b949c63445f1c3262" ;;
        *) echo "" ;;
    esac
}

hash_file_sha256() {
    local path="$1"

    if command -v sha256sum &> /dev/null; then
        sha256sum "$path" | awk '{print $1}'
        return 0
    fi

    if command -v shasum &> /dev/null; then
        shasum -a 256 "$path" | awk '{print $1}'
        return 0
    fi

    echo "❌ Error: sha256sum or shasum is required but not found" >&2
    return 1
}

verify_rustfs_binary() {
    local binary_path="$1"
    local expected="$2"

    if [ ! -f "$binary_path" ]; then
        echo "❌ Error: rustfs binary not found at $binary_path"
        return 1
    fi

    local binary_sha
    binary_sha=$(hash_file_sha256 "$binary_path")
    if [ "$binary_sha" != "$expected" ]; then
        echo "❌ ERROR: Binary SHA256 does not match pinned hash"
        echo "   Got:      $binary_sha"
        echo "   Expected: $expected"
        return 1
    fi
    return 0
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --install)
            INSTALL=true
            shift
            ;;
        --dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Installs rustfs (pinned to v${RUSTFS_VERSION}) for macOS (aarch64) or Linux (x86_64/aarch64)"
            echo ""
            echo "Options:"
            echo "  --dir DIR         Installation directory (default: /usr/local/bin)"
            echo "  --install         Install after download"
            echo "  --help            Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

PLATFORM="$(platform_key)"
if [ -z "$PLATFORM" ]; then
    echo "❌ Error: unsupported platform $(uname -s)/$(uname -m)" 1>&2
    echo "   rustfs is only pinned for macOS/arm64 and Linux/x86_64,arm64 here." 1>&2
    exit 1
fi

ZIP_SHA256="$(zip_sha256_for "$PLATFORM")"
BINARY_SHA256="$(binary_sha256_for "$PLATFORM")"

# Check if rustfs is already installed. Prefer the target install location
# (so --dir installs are checked directly, not just whatever is on PATH)
# and fall back to PATH resolution.
RUSTFS_BIN=""
if [ -f "$INSTALL_DIR/rustfs" ]; then
    RUSTFS_BIN="$INSTALL_DIR/rustfs"
elif command -v rustfs &> /dev/null; then
    RUSTFS_BIN=$(command -v rustfs)
fi

if [ -n "$RUSTFS_BIN" ]; then
    echo "Found rustfs at: $RUSTFS_BIN"

    echo "Verifying rustfs binary hash..."
    if verify_rustfs_binary "$RUSTFS_BIN" "$BINARY_SHA256"; then
        INSTALLED_VERSION=$("$RUSTFS_BIN" --version 2>/dev/null || echo "unknown")
        echo "✓ rustfs is already installed and verified: $INSTALLED_VERSION"
        exit 0
    else
        if [ "$INSTALL" = false ]; then
            echo "❌ ERROR: Installed rustfs binary failed verification"
            exit 1
        fi

        echo "Installed rustfs at $RUSTFS_BIN does not match the pinned binary hash."
        echo "Downloading the exact pinned release artifact to replace it..."
    fi
fi

echo "Checking for rustfs availability..."

DOWNLOAD_URL="https://github.com/rustfs/rustfs/releases/download/${RUSTFS_VERSION}/rustfs-${PLATFORM}-v${RUSTFS_VERSION}.zip"

echo "Found rustfs release: $RUSTFS_VERSION ($PLATFORM)"
echo "Download URL: $DOWNLOAD_URL"

if [ "$INSTALL" = false ]; then
    echo ""
    echo "To install rustfs, run:"
    echo "  $0 --install"
    exit 0
fi

if ! command -v unzip &> /dev/null; then
    echo "❌ Error: unzip is required but not found"
    exit 1
fi

TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

echo "Downloading rustfs..."
if ! curl -L -o "$TEMP_DIR/rustfs.zip" "$DOWNLOAD_URL"; then
    echo "❌ Error: Failed to download rustfs"
    exit 1
fi

echo "Verifying download integrity..."
DOWNLOADED_SHA=$(hash_file_sha256 "$TEMP_DIR/rustfs.zip")
if [ "$DOWNLOADED_SHA" != "$ZIP_SHA256" ]; then
    echo "❌ ERROR: Downloaded file SHA256 does not match pinned hash"
    echo "   Downloaded: $DOWNLOADED_SHA"
    echo "   Expected:   $ZIP_SHA256"
    exit 1
fi
echo "✓ Download integrity verified"

echo "Extracting rustfs..."
unzip -q "$TEMP_DIR/rustfs.zip" -d "$TEMP_DIR"

RUSTFS_BIN=$(find "$TEMP_DIR" -name "rustfs" -type f | head -1)

if [ -z "$RUSTFS_BIN" ] || [ ! -f "$RUSTFS_BIN" ]; then
    echo "❌ Error: Could not find rustfs binary in downloaded archive"
    exit 1
fi

# rustfs is distributed unsigned, so hash verification is the only integrity check.
echo "Verifying rustfs binary hash..."
if ! verify_rustfs_binary "$RUSTFS_BIN" "$BINARY_SHA256"; then
    exit 1
fi
echo "✓ Binary hash verified"

mkdir -p "$INSTALL_DIR"
echo "Installing rustfs to $INSTALL_DIR..."
cp "$RUSTFS_BIN" "$INSTALL_DIR/rustfs"
chmod +x "$INSTALL_DIR/rustfs"

INSTALLED_RUSTFS_BIN="$INSTALL_DIR/rustfs"
echo "Verifying installed rustfs binary hash..."
if verify_rustfs_binary "$INSTALLED_RUSTFS_BIN" "$BINARY_SHA256"; then
    INSTALLED_VERSION=$("$INSTALLED_RUSTFS_BIN" --version 2>/dev/null || echo "unknown")
    echo "✓ rustfs installed successfully and verified: $INSTALLED_VERSION"
    echo "Location: $INSTALLED_RUSTFS_BIN"
else
    echo "❌ Error: Installed rustfs binary failed verification"
    exit 1
fi

if ! command -v rustfs &> /dev/null; then
    echo "⚠ rustfs installed but not in PATH"
    echo "Add $INSTALL_DIR to your PATH"
elif [ "$(command -v rustfs)" != "$INSTALLED_RUSTFS_BIN" ]; then
    echo "⚠ Another rustfs binary appears earlier in PATH: $(command -v rustfs)"
    echo "Use $INSTALLED_RUSTFS_BIN directly or update your PATH"
fi
