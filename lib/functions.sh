function logerr() {
    >&2 echo "$@"
}

function log() {
    # Log only when env variable isset.
    if [[ "yes" == "$IDE_DEBUG" ]]
    then
        echo "[$IDE_PREFIX] $@"
    fi
}

function log_section() {
    echo -e "==> \033[1;32m$1\033[0m"
}

function __ide_detect_project_language() {
    if   [ -f 'setup.py' ];          then echo 'python/package'
    elif [ -f 'app/AppKernel.php' ]; then echo 'php/symfony2'
    elif [ -f 'composer.json' ];     then echo 'php/package'
    elif [ -f 'manage.py' ];         then echo 'python/django'
    else
        echo 'generic/plain'
    fi
}

function __get_default_ctags_options() {
    local OPTS="--recurse --exclude=.git --exclude=.hg --exclude=.svn"

    case "$1" in
        php/*)
            OPTS="$OPTS --php-kinds=cif --fields=+aimS --languages=php"
            ;;

        python/*)
            OPTS="$OPTS --python-kinds=-i --languages=python"
            ;;
    esac

    echo "$OPTS"
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
