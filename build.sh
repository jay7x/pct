#!/bin/bash
target=${1:-build}
arch=$(go env GOHOSTARCH)
platform=$(go env GOHOSTOS)
binPath="$(pwd)/dist/pct_${platform}_${arch}"

amd64=$(go env GOAMD64)
if [[ -n "$amd64" ]]; then
	binPath="${binPath}_${amd64}"
fi

if [ "$target" == "build" ]; then
  # Set goreleaser to build for current platform only
  goreleaser build --snapshot --clean --single-target
  git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/baker-round "$binPath/templates"
elif [ "$target" == "quick" ]; then
  go build -o "${binPath}"/pct
elif [ "$target" == "package" ]; then
  git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/baker-round "templates"
  goreleaser --skip-publish --snapshot --clean
fi
