#!/usr/bin/env bash
set -e

VERSION="$1"
ARCH="darwin_arm64"

if [ -z "$VERSION" ]; then
  echo "Usage: ./make_release.sh <version>"
  exit 1
fi

ROOT_DIR="$(pwd)"
DIST_DIR="$ROOT_DIR/dist"
PKG_DIR="$DIST_DIR/picasso_${VERSION}_${ARCH}"

rm -rf "$DIST_DIR"
mkdir -p "$PKG_DIR"

echo "Building picasso binary..."
bazel build //cli:picasso
cp bazel-bin/cli/picasso "$PKG_DIR/picasso"


echo "Copying runtime files..."
cp -R "$ROOT_DIR/libs" "$PKG_DIR/libs"
cp -R "$ROOT_DIR/runtime" "$PKG_DIR/runtime"
cp "bazel-bin/libruntime_lib.a" "$PKG_DIR"
cp "bazel-bin/irgen/irgen_/irgen" "$PKG_DIR"

echo "Creating tarball..."
cd "$DIST_DIR"
tar -czf "picasso_${VERSION}_${ARCH}.tar.gz" "picasso_${VERSION}_${ARCH}"

echo "Done:"
echo "$DIST_DIR/picasso_${VERSION}_${ARCH}.tar.gz"