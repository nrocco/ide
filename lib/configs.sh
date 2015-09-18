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
