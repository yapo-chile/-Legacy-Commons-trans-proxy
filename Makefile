include mk/vars.mk
include mk/help.mk
include mk/colors.mk

## Run tests and generate quality reports
test:
	@scripts/commands/test.sh

## Run tests and output coverage reports
cover:
	@scripts/commands/test_cover.sh cli

## Run tests and open report on default web browser
coverhtml:
	@scripts/commands/test_cover.sh html

## Run gometalinter and output report as text
checkstyle:
	@scripts/commands/test_style.sh display

## Install golang system level dependencies
setup:
	@scripts/commands/setup.sh

## Compile the code locally
build-local:
	@scripts/commands/build.sh

## Execute the service
run-local: build-local
	@./${APPNAME}

## Compile and start the service using docker
docker-start:  docker-compose-up 

## Stop docker containers
docker-stop: docker-compose-down

## Setup a new service repository based on trans
clone:
	@scripts/commands/clone.sh

## Run gofmt to reindent source
fix-format:
	@scripts/commands/fix-format.sh

## Display basic service info
info:
	@echo "YO           : ${YO}"
	@echo "ServerRoot   : ${SERVER_ROOT}"
	@echo "API Base URL : ${BASE_URL}"
	@echo "Healthcheck  : curl ${BASE_URL}/api/v1/healthcheck"
	@echo "Images from latest commit:"
	@echo -e "- ${DOCKER_IMAGE}:${DOCKER_TAG}"
	@echo -e "- ${DOCKER_IMAGE}:${COMMIT_DATE_UTC}"

deploy-k8s:
	@scripts/commands/deploy-k8s.sh

include mk/docs.mk
include mk/docker.mk
include mk/dev.mk
include mk/test.mk
include mk/deploy.mk
