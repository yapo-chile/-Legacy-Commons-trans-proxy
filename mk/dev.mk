## Build and start the service in development mode (detached)
run: build-dev "docker-compose-up -d"

## Build and start the service in development mode (attached)
start: docker-compose-up

## Stop running services
stop: docker-compose-down

.PHONY: run start stop

## Run docker compose commands with the project configuration
dc-%:
	docker-compose -f docker/docker-compose.yml \
		--project-name ${APPNAME} \
		--project-directory . \
		$*
