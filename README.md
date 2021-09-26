ESIEE API
====================

This is the unofficial API for ESIEE Paris services (ADE, Aurion, etc.)

## Versioning

Simply override the application version in the `VERSION` file.

## Building the project

The `Makefile` contains all neccessary goals for build steps.

Simply type in `make help` to list all available goals.

To build and deploy :
- Increase the version number in `VERSION`
- `make build-docker`
- `docker login`
- `make push-docker`

## Running

Create a config.yaml
Simply run the application `./esiee-api`

The application then prints the following output:

```
esiee-api git:(master) ✗ ./esiee-api     
2021-09-26 11:42:06.590 INFO (main.go:116) main-logger Running in debug : false
2021-09-26 11:42:06.591 INFO (main.go:120) main-logger Starting Gocron
```

## Building and running in Docker

Build the application using `make build-docker` and run the following command: `docker run -d -p 8080:8080 wallforfry/esiee-api:latest`
