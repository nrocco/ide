#!/bin/bash

if [[ -z $1 || -z $2 ]]
then
    echo "Usage: $0 from to"
    exit 1
fi

if ! git --no-pager diff --exit-code --quiet --name-only $1 $2 -- app/DoctrineMigrations
then
    echo "[symfony2/detect-database-changes] Changes in migrations detected! you might want to run php app/console doctrine:migrations:migrate"
fi
