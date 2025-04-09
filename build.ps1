$ErrorActionPreference = "Stop"

$platforms = @(
    @{GOOS="linux"; GOARCH="amd64"},
    @{GOOS="linux"; GOARCH="arm64"},
    @{GOOS="darwin"; GOARCH="amd64"},
    @{GOOS="darwin"; GOARCH="arm64"},
    @{GOOS="windows"; GOARCH="amd64"}
)

New-Item -ItemType Directory -Force -Path "bin" | Out-Null

foreach ($platform in $platforms) {
    $outputName = "bin/ndu-$($platform.GOOS)-$($platform.GOARCH)"
    if ($platform.GOOS -eq "windows") {
        $outputName += ".exe"
    }

    Write-Host "Building for $($platform.GOOS)/$($platform.GOARCH)..."
    $env:GOOS = $platform.GOOS
    $env:GOARCH = $platform.GOARCH
    go build -o $outputName ./cmd/ndu
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error building for $($platform.GOOS)/$($platform.GOARCH)"
        exit 1
    }
}

Write-Host "Build completed successfully!" 