#!/usr/bin/env bash
set -euxo pipefail
rm -rf release
mkdir -p release
project_name=$(basename "${PWD}")
project_name="${project_name//\//}"
# build for all common arch and platforms
for supported_os in darwin linux windows
do
    for supported_arch in amd64 arm64
    do
        file_ext=""
        if [[ "$supported_os" == "windows" ]]; then
            file_ext=".exe"
        fi
        output_name="${project_name}_${supported_os}_${supported_arch}${file_ext}"
        GOOS="$supported_os" GOARCH="$supported_arch" go build -o "release/$output_name" main.go
    done
done
pushd release
sha256sum "${project_name}"* |tee "${project_name}".sha256sums
gpg --armor --output "${project_name}".sha256sums.asc --detach-sign "${project_name}".sha256sums
gpg --verify "${project_name}".sha256sums.asc "${project_name}".sha256sums

echo "done"