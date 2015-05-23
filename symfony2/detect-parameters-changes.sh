#!/bin/bash

if [[ -z $1 || -z $2 ]]
then
    echo "Usage: $0 from to"
    exit 1
fi

if ! git --no-pager diff --exit-code --quiet --name-only $1 $2 -- app/config/parameters.yml.dist
then
    echo "[symfony2/detect-parameters-changes] Changes in parameters.yml detected! you might want to check it out."
fi
