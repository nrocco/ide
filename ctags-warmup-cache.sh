#!/bin/bash -e

CTAGS_CONFIG='/home/nrocco/ide/ctags-php'
GIT_DIR=`git rev-parse --git-dir`

if [[ '--vendors' == "$1" ]]
then

    time (
        find vendor -type f -name '*.php' | \
            /usr/bin/ctags -L - --tag-relative --options=${CTAGS_CONFIG} -o ${GIT_DIR}/tags_vendors
    )

else

    time (
        git ls-files | /bin/grep -E 'php$' | \
            /usr/bin/ctags -L - --tag-relative --options=${CTAGS_CONFIG} -o ${GIT_DIR}/tags
    )

fi
