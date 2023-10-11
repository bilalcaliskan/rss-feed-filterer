#!/usr/bin/env bash

# Fetch the module path from go.mod
MODULE_PATH=$(grep module go.mod | awk '{print $2}')

fail() {
  echo "$1"
  exit 1
}


# get the directory of the current script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"


# initialize ARGS with any provided arguments that don't end in .go
ARGS=""
declare -A packages
for arg in "$@"; do
    if [[ $arg != *.go ]]; then
        ARGS="$ARGS $arg"
    else
        pkg=$(dirname "$arg")
        packages["$pkg"]=1
    fi
done

# if there are any .go files, run tests for their respective packages
if [[ ${#packages[@]} -ne 0 ]]; then
    for pkg in "${!packages[@]}"; do
        echo "running vet for package: $pkg"
        full_pkg_path="$MODULE_PATH/$pkg"
        go vet $ARGS "$full_pkg_path" || fail "vet"
    done
else
    echo "no .go files found in arguments"
fi
