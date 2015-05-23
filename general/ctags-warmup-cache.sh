#!/bin/bash -e

# Log only when env variable isset.
function log() {
    if [[ "yes" == "$IDE_DEBUG" ]]
    then
        echo "[general/ctags-warmup-cache] $@"
    fi
}


GIT_DIR=`git rev-parse --git-dir`
CTAGS_CONFIG=`git config --global --null ide.ctags.config`


log "Warming up the ctags cache..."
log "Using ctags options file $CTAGS_CONFIG"

if [[ '--vendors' == "$1" ]]
then

    log "Creating tags in ${GIT_DIR}/tags_vendors"

    time (
        find vendor -type f -name '*.php' | \
            /usr/bin/ctags -L - --tag-relative --options=${CTAGS_CONFIG} -o ${GIT_DIR}/tags_vendors
    )

else

    log "Creating tags in ${GIT_DIR}/tags"

    time (
        git ls-files | /bin/grep -E 'php$' | \
            /usr/bin/ctags -L - --tag-relative --options=${CTAGS_CONFIG} -o ${GIT_DIR}/tags
    )

fi

log "Done."
