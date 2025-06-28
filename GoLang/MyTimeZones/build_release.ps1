param(
    [string]$Version = "v1.0.0"
)

# Create output directories
$base = "Releases\$Version"
$win = Join-Path $base "Win"
$linux = Join-Path $base "Linux"
$mac = Join-Path $base "Mac"
New-Item -ItemType Directory -Force -Path $win, $linux, $mac | Out-Null

# Build for Windows
$env:GOOS="windows"; $env:GOARCH="amd64"
go build -ldflags="-s -w -H=windowsgui" -x -v -o "$win\MyTimeZones.exe"
.\upx-5.0.1-win64\upx.exe --best -v -V -l "$win\MyTimeZones.exe"

#export CGO_ENABLED=1
# Build for Linux
$env:GOOS="linux"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -x -v -o "$linux\MyTimeZones"
#upx --best "$linux\MyWorkTimes"

# Build for Mac
$env:GOOS="darwin"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -x -v -o "$mac\MyTimeZones"
#upx --best "$mac\MyWorkTimes"

# Clean up env vars
Remove-Item Env:GOOS, Env:GOARCH
Write-Host "Builds complete! Check the Releases\$Version folder."

# Tag the release in git
#git tag $Version
#git push --tags

#Write-Host "Git tag $Version created and pushed."