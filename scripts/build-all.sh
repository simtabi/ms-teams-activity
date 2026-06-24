#!/usr/bin/env bash
# Build ready-to-run binaries for all targets into ./dist, with archives and a
# checksums.txt — the same layout the release pipeline produces. macOS targets
# need cgo (a C toolchain); Windows/Linux are pure Go.
#
#   ./scripts/build-all.sh [version]
set -euo pipefail
cd "$(dirname "$0")/.."

VERSION="${1:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
LDFLAGS="-s -w -X github.com/simtabi/ms-teams-activity/internal/cli.version=${VERSION}"
DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
LDFLAGS="${LDFLAGS} -X github.com/simtabi/ms-teams-activity/internal/cli.date=${DATE}"

rm -rf dist && mkdir -p dist
host_os="$(go env GOOS)"; host_arch="$(go env GOARCH)"

build() { # goos goarch cgo cc
  local goos="$1" goarch="$2" cgo="$3" cc="${4:-}"
  local out="dist/mta_${goos}_${goarch}" bin="mta"
  [ "$goos" = "windows" ] && bin="mta.exe"
  mkdir -p "$out"
  echo ">> ${goos}/${goarch} (cgo=${cgo})"
  if ! env GOOS="$goos" GOARCH="$goarch" CGO_ENABLED="$cgo" ${cc:+CC="$cc"} \
        go build -trimpath -ldflags "$LDFLAGS" -o "${out}/${bin}" ./cmd/mta; then
    echo "   skipped ${goos}/${goarch} (toolchain unavailable)"; rm -rf "$out"; return 0
  fi
  if [ "$goos" = "windows" ]; then
    (cd "$out" && zip -q "../mta_${goos}_${goarch}.zip" "$bin")
  else
    tar -C "$out" -czf "dist/mta_${goos}_${goarch}.tar.gz" "$bin"
  fi
}

# Linux + Windows are pure Go and always cross-build.
build linux   amd64 0
build linux   arm64 0
build windows amd64 0

# macOS needs cgo; the native arch always works. The other macOS arch only
# cross-builds when clang can target it (true on a normal Xcode install).
if [ "$host_os" = "darwin" ]; then
  build darwin "$host_arch" 1
  if [ "$host_arch" = "arm64" ]; then
    build darwin amd64 1 "clang -arch x86_64"
  else
    build darwin arm64 1 "clang -arch arm64"
  fi
else
  echo ">> skipping macOS builds (not on a Mac; cgo + macOS SDK required)"
fi

( cd dist && shasum -a 256 ./*.tar.gz ./*.zip 2>/dev/null > checksums.txt || true )

echo
echo "Artifacts in ./dist:"
ls -1 dist/*.tar.gz dist/*.zip 2>/dev/null
echo
echo "Ready-to-run binaries:"
find dist -type f -name 'mta' -o -name 'mta.exe' | sed 's/^/  /'
