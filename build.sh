#!/bin/bash

set -e

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

mkdir -p bin

for platform in "${PLATFORMS[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="bin/ndu-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building for $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name ./cmd/ndu
    if [ $? -ne 0 ]; then
        echo "Error building for $GOOS/$GOARCH"
        exit 1
    fi
done

echo "Build completed successfully!" 