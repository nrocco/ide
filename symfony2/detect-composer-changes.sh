#!/bin/bash

if [[ -z $1 || -z $2 ]]
then
    echo "Usage: $0 from to"
    exit 1
fi

if ! git --no-pager diff --exit-code --quiet --name-only $1 $2 -- composer.json composer.lock
then
    echo "==> changes in composer detected! you might want to run composer install"
fi
