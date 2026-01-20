#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

# -----------------------------------------------------------------------------
# Path Resolution
# -----------------------------------------------------------------------------
# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Detect Project Root (works if script is in root or inside ./scripts/)
if [[ -d "$SCRIPT_DIR/watsap" ]]; then
    PROJECT_ROOT="$SCRIPT_DIR"
elif [[ -d "$SCRIPT_DIR/../watsap" ]]; then
    PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
else
    echo -e "\033[0;31m[ERROR] Could not find 'watsap' source directory. Run this from project root or scripts/ folder.\033[0m"
    exit 1
fi

# Define paths
readonly SOURCE_DIR="$PROJECT_ROOT/watsap"
readonly DIST_DIR="$PROJECT_ROOT/dist"
readonly ENV_FILE="$PROJECT_ROOT/.env"

# ANSI Colors
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

# -----------------------------------------------------------------------------
# Helper Functions
# -----------------------------------------------------------------------------

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1" >&2; }

check_dependencies() {
    local dependencies=("go" "upx")
    for cmd in "${dependencies[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            log_error "Command '$cmd' is missing. Please install it."
            exit 1
        fi
    done
}

load_env() {
    if [[ ! -f "$ENV_FILE" ]]; then
        log_warn ".env file not found at: $ENV_FILE"
        cat > "$ENV_FILE" <<EOF
export TG_BOT_TOKEN=""
export TG_CHAT_ID=""
EOF
        log_error "Created template .env. Please fill it and run again."
        exit 1
    fi

    # Securely source .env
    set -a
    source "$ENV_FILE"
    set +a

    if [[ -z "${TG_BOT_TOKEN:-}" || -z "${TG_CHAT_ID:-}" ]]; then
        log_error "TG_BOT_TOKEN or TG_CHAT_ID is missing in .env file."
        exit 1
    fi
    log_info "Environment variables loaded."
}

prepare_dirs() {
    if [[ -d "$DIST_DIR" ]]; then
        log_info "Cleaning old builds..."
        rm -f "$DIST_DIR"/watsap-*
    else
        mkdir -p "$DIST_DIR"
    fi
}

# -----------------------------------------------------------------------------
# Main Execution
# -----------------------------------------------------------------------------

clear
check_dependencies
load_env
prepare_dirs

# --- Step 1: Architecture Selection ---
echo "---------------------------------------"
echo "Select Target Architecture:"
echo "1) amd64 (64-bit)"
echo "2) 386   (32-bit)"
read -r -p "Choice [1-2]: " arch_choice

case "$arch_choice" in
    1) GOARCH="amd64" ;;
    2) GOARCH="386" ;;
    *) log_error "Invalid selection."; exit 1 ;;
esac

# --- Step 2: Build Type Selection ---
echo "---------------------------------------"
echo "Select Build Type:"
echo "1) Release (Optimized, Compressed, Hidden Console on Windows)"
echo "2) Debug   (Symbols included, Console visible)"
read -r -p "Choice [1-2]: " build_type

# --- Step 3: Compilation ---

# Critical: Switch to the source directory containing go.mod
log_info "Switching to source directory: $SOURCE_DIR"
cd "$SOURCE_DIR" || exit

# Inject environment variables into the binary
BASE_LDFLAGS="-X 'watsap/utils/config.TG_BOT_TOKEN=$TG_BOT_TOKEN' -X 'watsap/utils/config.TG_CHAT_ID=$TG_CHAT_ID'"

if [[ "$build_type" == "1" ]]; then
    # === RELEASE MODE ===
    log_info "Building RELEASE version ($GOARCH)..."

    # Linux Build (-s -w strips symbols)
    CGO_ENABLED=0 GOOS=linux GOARCH="$GOARCH" go build \
        -ldflags "$BASE_LDFLAGS -s -w" \
        -o "$DIST_DIR/watsap-linux-$GOARCH.bin" .

    # Windows Build (-H=windowsgui hides console)
    CGO_ENABLED=0 GOOS=windows GOARCH="$GOARCH" go build \
        -ldflags "$BASE_LDFLAGS -s -w -H=windowsgui" \
        -o "$DIST_DIR/watsap-windows-$GOARCH.exe" .

    # Compression
    log_info "Compressing binaries with UPX..."
    upx -9 -q -f --no-owner "$DIST_DIR/watsap-windows-$GOARCH.exe" "$DIST_DIR/watsap-linux-$GOARCH.bin" >/dev/null

else
    # === DEBUG MODE ===
    log_info "Building DEBUG version ($GOARCH)..."

    CGO_ENABLED=0 GOOS=linux GOARCH="$GOARCH" go build \
        -ldflags "$BASE_LDFLAGS" \
        -o "$DIST_DIR/watsap-linux-$GOARCH-debug.bin" .

    CGO_ENABLED=0 GOOS=windows GOARCH="$GOARCH" go build \
        -ldflags "$BASE_LDFLAGS" \
        -o "$DIST_DIR/watsap-windows-$GOARCH-debug.exe" .
fi

# --- Final Report ---
echo "---------------------------------------"
log_info "Build finished successfully!"
ls -lh "$DIST_DIR"
