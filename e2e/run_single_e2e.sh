#!/usr/bin/env bash
set -euo pipefail

RUNNER=$(readlink -f "$1")
NIYAMA=$(readlink -f "$2")
TEST_DIR_NAME="$3"
shift 3

WORK_DIR="${TEST_TMPDIR}/scratch"
rm -rf "$WORK_DIR"
mkdir -p "$WORK_DIR"

MODE=""
IRGEN_PATH=""
RUNTIME_LIB_PATH=""
TEST_FILES=()

for arg in "$@"; do
    if [[ "$arg" == "--deps" ]]; then MODE="deps"; continue; fi
    if [[ "$arg" == "--files" ]]; then MODE="files"; continue; fi

    if [[ "$MODE" == "deps" ]]; then
        if [[ "$arg" == *"irgen" ]] && [[ ! "$arg" == *".a" ]] && [[ ! "$arg" == *".so" ]]; then
            IRGEN_PATH=$(readlink -f "$arg")
        elif [[ "$arg" == *".a" ]]; then
            RUNTIME_LIB_PATH=$(readlink -f "$arg")
        fi
    elif [[ "$MODE" == "files" ]]; then
        TEST_FILES+=("$arg")
    fi
done

echo "Populating workspace for $TEST_DIR_NAME..."

for FILE in "${TEST_FILES[@]}"; do
    if [[ -d "$FILE" ]] || [[ "$FILE" == *"TEST_DIR" ]]; then continue; fi
    REL_PATH=${FILE#*"$TEST_DIR_NAME/"}
    if [[ "$REL_PATH" == "$FILE" ]]; then continue; fi

    DEST="$WORK_DIR/$REL_PATH"
    mkdir -p "$(dirname "$DEST")"
    cp -L "$FILE" "$DEST"
    chmod +w "$DEST"
done

if [[ -n "$IRGEN_PATH" ]]; then
    mkdir -p "$WORK_DIR/irgen/irgen_"
    ln -sf "$IRGEN_PATH" "$WORK_DIR/irgen/irgen_/irgen"
fi

if [[ -n "$RUNTIME_LIB_PATH" ]]; then
    ln -sf "$RUNTIME_LIB_PATH" "$WORK_DIR/libruntime_lib.a"
fi



cd "$WORK_DIR"
echo "--- Ready: $TEST_DIR_NAME ---"
exec "$RUNNER" "$NIYAMA" "."