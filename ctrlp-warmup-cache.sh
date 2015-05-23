#!/bin/bash -e

IGNORE_DIRS_DEFAULTS=".git .svn"
IGNORE_DIRS=( ${IGNORE_DIRS[@]} ${IGNORE_DIRS_DEFAULTS[@]} )

CACHE_DIR="$HOME/.cache/ctrlp"
GIT_DIR=`git rev-parse --git-dir`

time (
    project="$(realpath .)"
    cache_file="${CACHE_DIR}/${project//\//%}.txt"
    prune=

    for dir in "${IGNORE_DIRS[@]}"
    do
        prune="$prune -path $dir -prune -or"
    done

    echo "==> Storying results in $cache_file"

    /usr/bin/find * $prune \( -type f -and -print \) > "$cache_file"
)
