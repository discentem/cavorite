#!/bin/bash
# Install the AWS CLI, used by the test scripts to talk to rustfs.
#
# macOS: downloads the official signed pkg and verifies its Developer ID
#        Installer signature/Team ID before installing.
# Linux: uses the distro package manager if available, otherwise falls back
#        to the official AWS CLI v2 zip installer.

set -e

INSTALL=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --install)
            INSTALL=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --install    Install the AWS CLI if not already present"
            echo "  --help       Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

if command -v aws &> /dev/null; then
    echo "✓ AWS CLI already installed: $(aws --version)"
    exit 0
fi

if [ "$INSTALL" = false ]; then
    echo "AWS CLI is not installed."
    echo "To install it, run:"
    echo "  $0 --install"
    exit 0
fi

install_macos() {
    local expected_issuer="AMZN Mobile LLC"
    local expected_team_id="94KV3E626L"
    local installed_bin_path="/usr/local/aws-cli/aws"
    local download_url="https://awscli.amazonaws.com/AWSCLIV2.pkg"
    local pkg_path="/tmp/AWSCLIV2.pkg"

    verify_awscli_binary() {
        local bin_path="$1"

        if [ ! -f "$bin_path" ]; then
            echo "❌ Error: AWS CLI binary not found at $bin_path"
            return 1
        fi

        if ! codesign --verify --strict "$bin_path" &>/dev/null; then
            echo "❌ Error: AWS CLI binary signature verification failed"
            return 1
        fi

        local team_id
        team_id=$(codesign -dvv "$bin_path" 2>&1 | grep "^TeamIdentifier=" | cut -d= -f2)
        if [ "$team_id" != "$expected_team_id" ]; then
            echo "❌ Error: Installed binary Team ID mismatch!"
            echo "   Expected: $expected_team_id"
            echo "   Got: ${team_id:-not set}"
            return 1
        fi

        echo "✓ AWS CLI binary signature verified (Team ID: $team_id)"
        return 0
    }

    verify_awscli_package() {
        local pkg_path="$1"

        echo "Verifying AWS CLI package signature..."
        if ! command -v pkgutil &> /dev/null; then
            echo "❌ ERROR: pkgutil is not available. Cannot verify package signature."
            return 1
        fi

        if ! pkgutil --check-signature "$pkg_path" &>/dev/null; then
            echo "❌ ERROR: Package signature verification failed"
            return 1
        fi
        echo "✓ Package signature verified"

        local cert_info actual_issuer actual_team_id
        cert_info=$(pkgutil --check-signature "$pkg_path" 2>&1)
        actual_issuer=$(echo "$cert_info" | grep "Developer ID Installer:" | head -1 | sed 's/.*Developer ID Installer: //' | sed 's/ (.*//')
        actual_team_id=$(echo "$cert_info" | grep "Developer ID Installer:" | head -1 | grep -oE '\([A-Z0-9]+\)' | tr -d '()')

        if [ "$actual_issuer" != "$expected_issuer" ]; then
            echo "❌ ERROR: Issuer mismatch!"
            echo "   Expected: $expected_issuer"
            echo "   Got: $actual_issuer"
            return 1
        fi

        if [ "$actual_team_id" != "$expected_team_id" ]; then
            echo "❌ ERROR: Team ID mismatch!"
            echo "   Expected: $expected_team_id"
            echo "   Got: $actual_team_id"
            return 1
        fi

        echo "✓ Issuer verified: $actual_issuer ($actual_team_id)"
        return 0
    }

    echo "Downloading AWS CLI..."
    if ! curl -L -o "$pkg_path" "$download_url" 2>/dev/null; then
        echo "❌ Error: Failed to download AWS CLI"
        exit 1
    fi

    if ! verify_awscli_package "$pkg_path"; then
        rm -f "$pkg_path"
        exit 1
    fi

    echo "Installing AWS CLI package..."
    if ! sudo installer -pkg "$pkg_path" -target / &>/dev/null; then
        echo "❌ Error: Failed to install AWS CLI"
        rm -f "$pkg_path"
        exit 1
    fi

    if ! verify_awscli_binary "$installed_bin_path"; then
        rm -f "$pkg_path"
        exit 1
    fi

    rm -f "$pkg_path"
}

install_linux() {
    if command -v apt-get &> /dev/null; then
        echo "Installing AWS CLI via apt-get..."
        sudo apt-get update -y
        sudo apt-get install -y awscli
        return
    fi

    echo "Installing AWS CLI via official installer..."
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" RETURN

    local arch url
    arch="$(uname -m)"
    case "$arch" in
        x86_64) url="https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" ;;
        aarch64|arm64) url="https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" ;;
        *) echo "❌ Error: unsupported architecture $arch" 1>&2; exit 1 ;;
    esac

    curl -L -o "$tmp_dir/awscliv2.zip" "$url"
    unzip -q "$tmp_dir/awscliv2.zip" -d "$tmp_dir"
    sudo "$tmp_dir/aws/install"
}

case "$(uname -s)" in
    Darwin) install_macos ;;
    Linux)  install_linux ;;
    *) echo "❌ Error: unsupported OS $(uname -s)" 1>&2; exit 1 ;;
esac

if ! command -v aws &> /dev/null; then
    echo "❌ Error: AWS CLI not found after installation"
    exit 1
fi

echo "✓ AWS CLI installed successfully: $(aws --version)"
