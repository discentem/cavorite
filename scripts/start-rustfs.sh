#!/bin/bash
# Start rustfs S3-compatible server and expose the PID

set -e

STORAGE_DIR="${HOME}/.cavorite-test-buckets"
ACCESS_KEY="blah"
SECRET_KEY="blah"
PORT="9000"
RUSTFS_BIN="${RUSTFS_BIN:-}"

while [[ $# -gt 0 ]]; do
    case $1 in
        --bin)
            RUSTFS_BIN="$2"
            shift 2
            ;;
        --dir)
            STORAGE_DIR="$2"
            shift 2
            ;;
        --access-key)
            ACCESS_KEY="$2"
            shift 2
            ;;
        --secret-key)
            SECRET_KEY="$2"
            shift 2
            ;;
        --port)
            PORT="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --bin PATH              Path to rustfs binary (default: /usr/local/bin/rustfs if present, else PATH)"
            echo "  --dir DIR              Storage directory for rustfs (default: ~/.cavorite-test-buckets)"
            echo "  --access-key KEY       S3 access key (default: blah)"
            echo "  --secret-key KEY       S3 secret key (default: blah)"
            echo "  --port PORT            Server port (default: 9000)"
            echo "  --help                 Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Resolve rustfs binary.
# Prefer the pinned install location so tests don't accidentally use an older
# rustfs earlier in PATH.
if [ -z "$RUSTFS_BIN" ] && [ -x "/usr/local/bin/rustfs" ]; then
    RUSTFS_BIN="/usr/local/bin/rustfs"
fi

if [ -z "$RUSTFS_BIN" ]; then
    if command -v rustfs &> /dev/null; then
        RUSTFS_BIN=$(command -v rustfs)
    else
        echo "❌ Error: rustfs not found in PATH"
        exit 1
    fi
fi

if [ ! -x "$RUSTFS_BIN" ]; then
    echo "❌ Error: rustfs binary is not executable: $RUSTFS_BIN"
    exit 1
fi

mkdir -p "$STORAGE_DIR"

echo "Starting rustfs server..." >&2
echo "  Binary: $RUSTFS_BIN" >&2
echo "  Storage: $STORAGE_DIR" >&2
echo "  Port: $PORT" >&2
echo "  Endpoint: http://localhost:$PORT" >&2

"$RUSTFS_BIN" server "$STORAGE_DIR" \
    --console-enable \
    --access-key "$ACCESS_KEY" \
    --secret-key "$SECRET_KEY" \
    --address ":$PORT" &> /tmp/rustfs.log &

RUSTFS_PID=$!

sleep 2

if ! kill -0 "$RUSTFS_PID" 2>/dev/null; then
    echo "❌ Error: rustfs server failed to start" >&2
    cat /tmp/rustfs.log >&2
    exit 1
fi

echo "✓ rustfs server started (PID: $RUSTFS_PID)" >&2
# Only output the PID to stdout so it can be captured
echo "$RUSTFS_PID"
