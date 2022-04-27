#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoTitle "Fixing format with gofmt:"
for file in ${GO_FILES}; do
    echo -n "checking ${file:2}"
    errors=$(gofmt -d -e -s ${file} | grep -c -E "^\+")
    if [ ${errors} -gt 0 ]; then
        echo " ... fixing $errors" issues
        gofmt -s -w $file
    else
        echo " ... ok"
    fi
done
echoTitle "Done"
