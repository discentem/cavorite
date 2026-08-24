#!/bin/bash
# Master test script that runs all setup and test steps in order.
# This orchestrates the complete local (and CI) testing workflow for
# cavorite: an S3 roundtrip test against rustfs, plus the localstore
# plugin roundtrip test (which needs no S3-compatible server at all).

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RUSTFS_BIN="/usr/local/bin/rustfs"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_section() {
    echo ""
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

error_exit() {
    echo -e "${RED}❌ Error: $1${NC}"
    exit 1
}

SKIP_INSTALLS=false
KEEP_SERVER_RUNNING=false
BUILD_ARGS=()

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-installs)
            SKIP_INSTALLS=true
            shift
            ;;
        --keep-server)
            KEEP_SERVER_RUNNING=true
            shift
            ;;
        --go-build)
            BUILD_ARGS+=(--go-build)
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Runs the complete local/CI testing workflow for cavorite"
            echo ""
            echo "Options:"
            echo "  --skip-installs    Skip dependency installation steps"
            echo "  --keep-server      Keep rustfs server running after tests"
            echo "  --go-build         Build cavorite/localstore with 'go build' instead of bazel"
            echo "  --help             Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

RUSTFS_PID=""
cleanup() {
    if [ -n "$RUSTFS_PID" ] && ! [ "$KEEP_SERVER_RUNNING" = true ]; then
        echo ""
        echo "Stopping rustfs server (PID: $RUSTFS_PID)..."
        kill "$RUSTFS_PID" 2>/dev/null || true
        wait "$RUSTFS_PID" 2>/dev/null || true
    fi
}
trap cleanup EXIT

# ========== Install Dependencies ==========
if [ "$SKIP_INSTALLS" = false ]; then
    log_section "STEP 1: Installing Dependencies"

    echo "Installing AWS CLI..."
    "$SCRIPT_DIR/install-awscli.sh" --install || error_exit "Failed to install AWS CLI"

    echo "Installing rustfs..."
    "$SCRIPT_DIR/install-rustfs.sh" --install || error_exit "Failed to install rustfs"
else
    log_section "Skipping Dependency Installation (--skip-installs)"
fi

# ========== Start rustfs Server ==========
log_section "STEP 2: Starting rustfs S3 Server"
# start-rustfs.sh backgrounds the actual server itself and prints its real PID
# on stdout, so it must run in the foreground here (not backgrounded with &)
# or RUSTFS_PID below would capture the wrapper script's PID instead, leaving
# the real server un-killable and leaked after this script exits.
RUSTFS_PID=$("$SCRIPT_DIR/start-rustfs.sh" --bin "$RUSTFS_BIN") || error_exit "Failed to start rustfs server"
echo "rustfs started with PID: $RUSTFS_PID"

echo "Waiting for rustfs to be ready..."
for i in {1..30}; do
    if nc -z localhost 9000 2>/dev/null; then
        echo "✓ rustfs is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        error_exit "rustfs server did not start"
    fi
    sleep 1
done

# ========== Create S3 Bucket ==========
log_section "STEP 3: Creating S3 Bucket"
"$SCRIPT_DIR/create-s3-bucket.sh" || error_exit "Failed to create S3 bucket"

# ========== S3 Roundtrip Test ==========
log_section "STEP 4: Running S3 (rustfs) Roundtrip Test"
"$SCRIPT_DIR/test-s3-roundtrip.sh" "${BUILD_ARGS[@]}" || error_exit "S3 roundtrip test failed"

# ========== localstore Plugin Roundtrip Test ==========
# No rustfs/S3 server needed for this one.
log_section "STEP 5: Running localstore Plugin Roundtrip Test"
"$SCRIPT_DIR/test-localstore-plugin.sh" "${BUILD_ARGS[@]}" || error_exit "localstore plugin test failed"

# ========== Complete ==========
log_section "✓ All Tests Completed Successfully!"

if [ "$KEEP_SERVER_RUNNING" = true ]; then
    echo "rustfs server is still running with PID: $RUSTFS_PID"
    echo "S3 endpoint available at: http://localhost:9000"
    echo ""
    echo "To stop the server later, run:"
    echo "  kill $RUSTFS_PID"
else
    echo "rustfs server will be stopped on exit"
fi
