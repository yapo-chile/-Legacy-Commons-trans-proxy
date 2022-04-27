#!/bin/bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

set -e

mkdir -p ${REPORT_ARTIFACTS}

CHECKSTYLE_FILE=${REPORT_ARTIFACTS}/checkstyle-report.xml

echoHeader "Running Checkstyle Tests"

COMMAND="golangci-lint -c golangci.yml"
if [[ $@ == **display** ]]; then
    COMMAND="${COMMAND} run ./... | tee /dev/tty > ${CHECKSTYLE_FILE} && echo"
else
    COMMAND="${COMMAND} --out-format \"colored-line-number\" run ./..."
fi
eval ${COMMAND}

