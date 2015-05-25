function __ide_get_subcommands() {
    ls "$IDE_PATH/exec" | xargs
}

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

function __ide_config() {
    local location=${3-}
    local value=

    if ! value=$(git config $location -z "$1")
    then
        if [ -z "$2" ]
        then
            return 1
        fi
        value="$2"
    fi

    echo -n $value
    return 0
}

function __ide_config_local() {
    __ide_config "$1" "$2" "--local"
}

function __ide_config_global() {
    __ide_config "$1" "$2" "--global"
}

# Make sure we are in a git repository
if ! GIT_DIR=$(git rev-parse --git-dir 2> /dev/null)
then
    logerr "Not in a git repo!"
    exit 1
fi

IDE_PREFIX=${IDE_PREFIX-ide/$(basename $0)}
