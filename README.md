# Trans ms (proxy)


This microservice acts as a proxy between other microservices and a Trans server. The params are passed in a JSON body, and it can be configured to limit what commands can be executed.


## How to run trans-proxy

* Create the dir: ` ~/go/src/gitlab.com/yapo_team/legacy/commons`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd  ~/go/src/gitlab.com/yapo_team/legacy/commons
  $ git clone git@gitlab.com:yapo_team/legacy/commons/trans-proxy.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd trans-proxy
  $ make docker-start
  ```

* To get a list of available commands:

  ```
  $ make help
  Targets:
  test                 Run tests and generate quality reports
  cover                Run tests and output coverage reports
  coverhtml            Run tests and open report on default web browser
  checkstyle           Run gometalinter and output report as text
  setup                Install golang system level dependencies
  build                Compile the code
  run                  Execute the service
  start                Compile and start the service
  docker-start         Compile and start the service using docker
  docker-stop          Stop docker containers
  clone                Setup a new service repository based on trans-proxy
  fix-format           Run gofmt to reindent source
  info                 Display basic service info
  docs-start           Starts godoc webserver with live docs for the project
  docs-stop            Stops godoc webserver if running
  docs-compile         Compiles static documentation to docs folder
  docs-update          Generates a commit updating the docs
  docs                 Opens the live documentation on the default web browser
  docker-build         Create docker image based on docker/dockerfile
  docker-publish       Push docker image to containers.mpi-internal.com
  docker-attach        Attach to this service's currently running docker container output stream
  docker-compose-up    Start all required docker containers for this service
  docker-compose-down  Stop all running 
  ```

* If you change the code:

  ```
  $ make docker-start
  ```

* How to run the tests

  ```
  $ make [cover|coverhtml]
  ```

* How to check format

  ```
  $ make checkstyle
  ```

## Endpoints
### GET  /api/v1/healthcheck
Reports whether the service is up and ready to respond.

> When implementing a new service, you MUST keep this endpoint
and update it so it replies according to your service status!

#### Request
No request parameters

#### Response
* Status: Ok message, representing service health

```javascript
200 OK
{
	"Status": "OK"
}
```

### POST  /api/v1/execute/{command}
Sends the specified command to a trans-proxy server with the given params in the JSON body

#### Request
params: A JSON object, where the fields are the name of trans-proxy params (lowercase), and the values are the values required
by the trans-proxy command
```javascript
{
	"params":{
		"param1":"value1",
		...
	}
}
```

#### Response

```javascript
200 OK
{
	"status": "TRANS_OK"
	"response" - A JSON field containing all the values returned by the trans command
	
}
```

#### Error responses
```javascript
400 Bad Request
{
	"status": "TRANS_ERROR"
	"response": {
		"error" - An error message
	}
}
```

```javascript
500 Internal Server Error
{
	"ErrorMessage" - An error message
}
```


