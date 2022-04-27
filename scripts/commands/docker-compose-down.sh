#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoTitle "Stopping Docker containers"
docker-compose -f docker/docker-compose.yml -p ${APPNAME} down
