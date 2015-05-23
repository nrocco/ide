#!/bin/bash

# Log only when env variable isset.
function log() {
    if [[ "yes" == "$IDE_DEBUG" ]]
    then
        echo "[general/ctrlp-warmup-cache] $@"
    fi
}


GIT_DIR=`git rev-parse --git-dir`
CACHE_DIR="$HOME/.cache/ctrlp"
LOCAL_IGNORE=`git config --local --null ide.ctrlp.ignore`
GLOBAL_IGNORE=`git config --global --null ide.ctrlp.ignore`
PRUNE_DIRS=( ${GLOBAL_IGNORE[@]} ${LOCAL_IGNORE[@]} )


log "Warming up the ctrlp cache..."
log "Global ide.ctrlp.ignore == ${GLOBAL_IGNORE[@]}"
log "Local  ide.ctrlp.ignore == ${LOCAL_IGNORE[@]}"
log "Cache dir is $CACHE_DIR"


time (

    project="$(realpath .)"
    cache_file="${CACHE_DIR}/${project//\//%}.txt"
    prune=

    for dir in "${PRUNE_DIRS[@]}"
    do
        prune="$prune -path $dir -prune -or"
    done

    log "Storing results in $cache_file"
    /usr/bin/find * $prune \( -type f -and -print \) > "$cache_file"

)

log "Done."
