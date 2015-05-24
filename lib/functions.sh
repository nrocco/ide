GIT_DIR=`git rev-parse --git-dir`

function logerr() {
    >&2 echo "$@"
}

function log() {
    # Log only when env variable isset.
    if [[ "yes" == "$IDE_DEBUG" ]]
    then
        echo "[${IDE_PREFIX:-ide}] $@"
    fi
}
