
## Create production docker image
build:
	@echoHeader "Building production docker image"
	@set -x
	${DOCKER} build \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		-f docker/dockerfile \
		--build-arg APPNAME=${APPNAME} \
		--build-arg GIT_COMMIT=${COMMIT} \
		--label appname=${APPNAME} \
		--label branch=${BRANCH} \
		--label build-date=${CREATION_DATE} \
		--label commit=${COMMIT} \
		--label commit-author=${CREATOR} \
		--label commit-date=${COMMIT_DATE} \
		.
	${DOCKER} tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${COMMIT_DATE_UTC}
	@set +x

.PHONY: build
