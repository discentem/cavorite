#!/bin/bash
# End-to-end test for the localstore file-store plugin: no S3-compatible
# server needed, everything is backed by a local directory.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

CAVORITE_BIN=""
LOCALSTORE_BIN=""
BUILD_MODE="bazel"

while [[ $# -gt 0 ]]; do
    case $1 in
        --cavorite-bin)
            CAVORITE_BIN="$2"
            shift 2
            ;;
        --localstore-bin)
            LOCALSTORE_BIN="$2"
            shift 2
            ;;
        --go-build)
            BUILD_MODE="go"
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --cavorite-bin PATH    Path to the cavorite binary (default: builds it, see --go-build)"
            echo "  --localstore-bin PATH  Path to the localstore plugin binary (default: builds it, see --go-build)"
            echo "  --go-build             Build both binaries with 'go build' instead of bazel"
            echo "  --help                 Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

if [ -z "$CAVORITE_BIN" ]; then
    if [ "$BUILD_MODE" = "go" ]; then
        echo "Building cavorite via go build..."
        CAVORITE_BIN="$(mktemp -d)/cavorite"
        (cd "$PROJECT_ROOT" && go build -o "$CAVORITE_BIN" .)
    else
        echo "Building cavorite via bazel..."
        bazel build //:cavorite --stamp
        CAVORITE_BIN="$(bazel cquery //:cavorite --output=files 2>/dev/null)"
    fi
fi
if [ -z "$LOCALSTORE_BIN" ]; then
    if [ "$BUILD_MODE" = "go" ]; then
        echo "Building localstore plugin via go build..."
        LOCALSTORE_BIN="$(mktemp -d)/localstore"
        (cd "$PROJECT_ROOT" && go build -o "$LOCALSTORE_BIN" ./plugins/localstore)
    else
        echo "Building localstore plugin via bazel..."
        bazel build //plugins/localstore:localstore --stamp
        LOCALSTORE_BIN="$(bazel cquery //plugins/localstore:localstore --output=files 2>/dev/null)"
    fi
fi

for bin_var in CAVORITE_BIN LOCALSTORE_BIN; do
    bin_path="${!bin_var}"
    if [ ! -f "$bin_path" ]; then
        echo "❌ Error: $bin_var not found at $bin_path"
        exit 1
    fi
done

CAVORITE_BIN="$(cd "$(dirname "$CAVORITE_BIN")" && pwd)/$(basename "$CAVORITE_BIN")"
LOCALSTORE_BIN="$(cd "$(dirname "$LOCALSTORE_BIN")" && pwd)/$(basename "$LOCALSTORE_BIN")"

# Bazel's output artifacts are read-only, and bazel-built (rules_go,
# fastbuild) binaries on macOS can be missing an LC_UUID load command, which
# newer dyld refuses to load ("missing LC_UUID load command"). Copy to a
# writable temp path and ad-hoc codesign it there to regenerate the load
# commands; this is a no-op on Linux.
if [ "$(uname -s)" = "Darwin" ] && command -v codesign &> /dev/null; then
    BIN_STAGE_DIR="$(mktemp -d)"
    cp "$CAVORITE_BIN" "$BIN_STAGE_DIR/cavorite"
    cp "$LOCALSTORE_BIN" "$BIN_STAGE_DIR/localstore"
    chmod +wx "$BIN_STAGE_DIR/cavorite" "$BIN_STAGE_DIR/localstore"
    codesign -f -s - "$BIN_STAGE_DIR/cavorite"
    codesign -f -s - "$BIN_STAGE_DIR/localstore"
    CAVORITE_BIN="$BIN_STAGE_DIR/cavorite"
    LOCALSTORE_BIN="$BIN_STAGE_DIR/localstore"
else
    chmod +x "$LOCALSTORE_BIN"
fi

WORK_DIR="$(mktemp -d)"
STORAGE_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR" "$STORAGE_DIR"' EXIT
cd "$WORK_DIR"

echo "Test workspace: $WORK_DIR"
echo "Plugin storage: $STORAGE_DIR"

echo ""
echo "Initializing cavorite repo against localstore plugin..."
"$CAVORITE_BIN" init "$WORK_DIR" \
    --backend_address "$STORAGE_DIR" \
    --store_type=plugin \
    --plugin_address="$LOCALSTORE_BIN"

TEST_FILE="testpkg.bin"
echo "Creating test file: $TEST_FILE"
head -c 65536 /dev/urandom > "$TEST_FILE"
EXPECTED_CHECKSUM="$(shasum -a 256 "$TEST_FILE" | awk '{print $1}')"

echo ""
echo "Uploading $TEST_FILE..."
"$CAVORITE_BIN" upload "$TEST_FILE" --vv

CFILE="${TEST_FILE}.cfile"
if [ ! -f "$CFILE" ]; then
    echo "❌ Error: expected metadata file $CFILE was not created"
    exit 1
fi

echo ""
echo "Verifying metadata file contents..."
if ! jq -e --arg name "$TEST_FILE" '.name == $name' "$CFILE" >/dev/null; then
    echo "❌ Error: $CFILE has unexpected 'name'"
    cat "$CFILE"
    exit 1
fi
if ! jq -e --arg checksum "$EXPECTED_CHECKSUM" '.checksum == $checksum' "$CFILE" >/dev/null; then
    echo "❌ Error: $CFILE has unexpected 'checksum'"
    cat "$CFILE"
    exit 1
fi
echo "✓ Metadata file verified"

if [ ! -f "$STORAGE_DIR/$TEST_FILE" ]; then
    echo "❌ Error: uploaded file not found in plugin storage at $STORAGE_DIR/$TEST_FILE"
    exit 1
fi
echo "✓ File present in plugin storage"

echo ""
echo "Removing local copy and retrieving via localstore plugin..."
rm -f "$TEST_FILE"

"$CAVORITE_BIN" retrieve "$CFILE" --vv

if [ ! -f "$TEST_FILE" ]; then
    echo "❌ Error: retrieve did not restore $TEST_FILE"
    exit 1
fi

ACTUAL_CHECKSUM="$(shasum -a 256 "$TEST_FILE" | awk '{print $1}')"
if [ "$ACTUAL_CHECKSUM" != "$EXPECTED_CHECKSUM" ]; then
    echo "❌ Error: retrieved file checksum mismatch"
    echo "   Expected: $EXPECTED_CHECKSUM"
    echo "   Got:      $ACTUAL_CHECKSUM"
    exit 1
fi

echo "✓ Retrieved file checksum matches"
echo ""
echo "✓ localstore plugin roundtrip test passed"
