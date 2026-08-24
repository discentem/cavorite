#!/bin/bash
# End-to-end test: init a cavorite repo against an S3-compatible endpoint
# (rustfs by default), upload a file, and retrieve it back.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

CAVORITE_BIN=""
BUCKET_NAME="cavorite-test"
S3_ENDPOINT="http://localhost:9000"
S3_REGION="us-east-1"
ACCESS_KEY="${AWS_ACCESS_KEY_ID:-blah}"
SECRET_KEY="${AWS_SECRET_ACCESS_KEY:-blah}"
BUILD_MODE="bazel"
WORK_DIR=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --bin)
            CAVORITE_BIN="$2"
            shift 2
            ;;
        --bucket)
            BUCKET_NAME="$2"
            shift 2
            ;;
        --endpoint)
            S3_ENDPOINT="$2"
            shift 2
            ;;
        --region)
            S3_REGION="$2"
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
            echo "  --bin PATH        Path to the cavorite binary (default: builds it, see --go-build)"
            echo "  --go-build        Build cavorite with 'go build' instead of bazel"
            echo "  --bucket NAME     S3 bucket to use (default: cavorite-test)"
            echo "  --endpoint URL    S3 endpoint URL (default: http://localhost:9000)"
            echo "  --region REGION   AWS region (default: us-east-1)"
            echo "  --help            Show this help message"
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

if [ ! -f "$CAVORITE_BIN" ]; then
    echo "❌ Error: cavorite binary not found at $CAVORITE_BIN"
    exit 1
fi
CAVORITE_BIN="$(cd "$(dirname "$CAVORITE_BIN")" && pwd)/$(basename "$CAVORITE_BIN")"

# Bazel's output artifacts are read-only, and bazel-built (rules_go,
# fastbuild) binaries on macOS can be missing an LC_UUID load command, which
# newer dyld refuses to load ("missing LC_UUID load command"). Copy to a
# writable temp path and ad-hoc codesign it there to regenerate the load
# commands; this is a no-op on Linux.
if [ "$(uname -s)" = "Darwin" ] && command -v codesign &> /dev/null; then
    SIGNED_CAVORITE_BIN="$(mktemp -d)/cavorite"
    cp "$CAVORITE_BIN" "$SIGNED_CAVORITE_BIN"
    chmod +w "$SIGNED_CAVORITE_BIN"
    codesign -f -s - "$SIGNED_CAVORITE_BIN"
    CAVORITE_BIN="$SIGNED_CAVORITE_BIN"
fi

WORK_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR"' EXIT
cd "$WORK_DIR"

echo "Test workspace: $WORK_DIR"

echo ""
echo "Initializing cavorite repo against $S3_ENDPOINT/$BUCKET_NAME..."
AWS_ACCESS_KEY_ID="$ACCESS_KEY" \
AWS_SECRET_ACCESS_KEY="$SECRET_KEY" \
"$CAVORITE_BIN" init "$WORK_DIR" \
    --backend_address "${S3_ENDPOINT}/${BUCKET_NAME}" \
    --store_type=s3 \
    --region="$S3_REGION"

TEST_FILE="testpkg.bin"
echo "Creating test file: $TEST_FILE"
head -c 65536 /dev/urandom > "$TEST_FILE"
EXPECTED_CHECKSUM="$(shasum -a 256 "$TEST_FILE" | awk '{print $1}')"

echo ""
echo "Uploading $TEST_FILE..."
AWS_ACCESS_KEY_ID="$ACCESS_KEY" \
AWS_SECRET_ACCESS_KEY="$SECRET_KEY" \
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

echo ""
echo "Removing local copy and retrieving from $S3_ENDPOINT..."
rm -f "$TEST_FILE"

AWS_ACCESS_KEY_ID="$ACCESS_KEY" \
AWS_SECRET_ACCESS_KEY="$SECRET_KEY" \
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
echo "✓ S3 roundtrip test passed"
