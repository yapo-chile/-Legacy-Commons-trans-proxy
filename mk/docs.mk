## Starts godoc webserver with live docs for the project
docs-start:
	godoc -http ${DOCS_HOST} &> /dev/null & echo $$! > .docs.pid
	while ! test "$$(curl -s --write-out "%{http_code}" ${DOCS_HOST} -o /dev/null)" = 200; do sleep 0.1 ; done

## Stops godoc webserver if running
docs-stop:
	cat .docs.pid | xargs kill || true

## Compiles static documentation to docs folder
docs-compile: docs-stop docs-start
	@scripts/commands/docs-compile.sh

## Generates a commit updating the docs
docs-update: docs-compile
	git add docs
	git commit -m "${DOCS_COMMIT_MESSAGE}"

## Opens the live documentation on the default web browser
docs: docs-stop docs-start
	open http://${DOCS_HOST}/pkg/gitlab.com/yapo_team/legacy/commons/trans-proxy/
