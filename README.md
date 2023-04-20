# EPAM Systems Golang Exercise

## Instructions

### Docker 
Run `docker build -t epam-systems:latest .` from root directory to build the Docker image.

In `build/` folder, a `docker-compose.yml` file has been provided for easier startup.

Run `docker-compose -f build/docker-compose.yml` from the root directory to start the necessary services. If you build the docker image with a different name, please change the image name in the `docker-compose.yml` file.

### Environment variables
`AUTH_SECRET`, `KAFKA_ADDR`, `DB_USER` and `DB_PWD` are mandatory environment variables that need to be passed in order for the service to work.

### Local

If you want to run the codebase locally, from project root run `go run *.go`. All environment variables can still be passed like in docker-compose.

Example : `AUTH_SECRET="zuTNubVdTIv2fLoNDsgHuDjcMBiA9ofV" KAFKA_ADDR="localhost:9092" DB_USER="user" DB_PWD="pass" go run *.go`

To run all tests, from root directory run `go test -cover -race ./...`