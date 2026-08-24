#!/usr/bin/env bash
# Build release packages for both cavorite and the localstore plugin.
#
# Usage: scripts/build-all-packages.sh [VERSION]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

"$SCRIPT_DIR/build-cavorite-package.sh" "$@"
"$SCRIPT_DIR/build-localstore-plugin-package.sh" "$@"

echo ""
echo "All packages built. See dist/"
