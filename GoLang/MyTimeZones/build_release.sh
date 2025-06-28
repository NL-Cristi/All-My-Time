#!/bin/bash
set -e

# Set your version here or pass as an argument
VERSION=${1:-v1.0.0}

BASE="Releases/$VERSION"
WIN="$BASE/Win"
LINUX="$BASE/Linux"
MAC="$BASE/Mac"

mkdir -p "$WIN" "$LINUX" "$MAC"

# Build for Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -x -v -o "$WIN/MyTimeZones.exe"
upx --best "$WIN/MyTimeZones.exe"

# Build for Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -x -v -o "$LINUX/MyTimeZones"
upx --best "$LINUX/MyTimeZones"

# Build for Mac
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -x -v -o "$MAC/MyTimeZones"
upx --best "$MAC/MyTimeZones"

echo "Builds complete! Check the Releases/$VERSION folder."

# Tag the release in git
git tag "$VERSION"
git push --tags

echo "Git tag $VERSION created and pushed."