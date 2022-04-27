#!/usr/bin/env bash
export UNAMESTR = $(uname)
export GO_FILES = $(shell find . -iname '*.go' -type f | grep -v vendor | grep -v pact) # All the .go files, excluding vendor/ and pact/
GENPORTOFF?=0
genport = $(shell expr ${GENPORTOFF} + \( $(shell id -u) - \( $(shell id -u) / 100 \) \* 100 \) \* 200 + 30100 + $(1))

# GIT variables
export BRANCH=$(shell git branch | sed -n 's/^\* //p')
export GIT_COMMIT=$(shell git rev-parse HEAD)
export GIT_COMMIT_DATE=$(shell TZ="America/Santiago" git show --quiet --date='format-local:%d-%m-%Y_%H:%M:%S' --format="%cd")
export COMMIT_DATE_UTC ?= $(shell TZ=UTC git show --quiet --date='format-local:%Y%m%d_%H%M%S' --format="%cd")
export BUILD_CREATOR=$(shell git log --format=format:%ae | head -n 1)

# REPORT_ARTIFACTS should be in sync with `RegexpFilePathMatcher` in
# `reports-publisher/config.json`
export REPORT_ARTIFACTS=reports

# APP variables
# This variables are for the use of your microservice. This variables must be updated each time you are creating a new microservice
export APPNAME=trans-proxy
export YO=`whoami`
export SERVICE_HOST=localhost
export SERVICE_PORT=8086
export SERVER_ROOT=${PWD}
export BASE_URL="http://${SERVICE_HOST}:${SERVICE_PORT}"
export MAIN_FILE=cmd/${APPNAME}/main.go
export LOGGER_SYSLOG_ENABLED=false
export LOGGER_SYSLOG_IDENTITY=trans-proxy
export LOGGER_STDLOG_ENABLED=true
export LOGGER_LOG_LEVEL=0


# Trans variables 
# Change to target a trans-proxy server,i.e.:
# TRANS_HOST=jenna.schibsted.cl
# TRANS_PORT=27205

# DOCKER variables
export DOCKER_REGISTRY=gitlab.com/yapo_team/legacy/commons/${APPNAME}

export DOCKER_TAG=$(shell echo ${BRANCH} | tr '[:upper:]' '[:lower:]' | sed 's,/,_,g')
export DOCKER_IMAGE=${DOCKER_REGISTRY}
export DOCKER_PORT=$(call genport,1)
export DOCKER ?= docker

# Documentation variables
export DOCS_DIR=docs
export DOCS_HOST=localhost:$(call genport,3)
export DOCS_PATH=github.mpi-internal.com/Yapo/${APPNAME}
export DOCS_COMMIT_MESSAGE=Generate updated documentation

export PROMETHEUS_ENABLED=true
