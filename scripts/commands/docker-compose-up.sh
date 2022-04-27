#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"


echoTitle "Creating needed networks"
for network in ${DOCKER_COMPOSE_NETWORKS}; do
    networkId=`docker network ls -q -f name=${network}`
    if [ -z "$networkId" ];
    then
        echo "Creating network ${network}"
        docker network create ${network}
    fi
done

echoTitle "Starting containers"
docker-compose -f docker/docker-compose.yml -p ${APPNAME} up

echoTitle "Done"
