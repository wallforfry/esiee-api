FONT_ESC := $(shell printf '\033')
FONT_BOLD := ${FONT_ESC}[1m
FONT_NC := ${FONT_ESC}[0m # No colour

all:
	@echo "Use a specific goal. To list all goals, type 'make help'"

.PHONY: version # Prints project version
version:
	@cat VERSION

.PHONY: run# Run the project
run: build
	./esiee-api

.PHONY: build # Builds the project
build:
	@go build

.PHONY: test # Runs unit tests
test:
	@go test -v ./...

.PHONY: build-docker # Build Docker image
build-docker:
	@docker build -t wallforfry/esiee-api:$(shell $(MAKE) version) -f Dockerfile .
	@docker tag wallforfry/esiee-api:$(shell $(MAKE) version) wallforfry/esiee-api:latest

.PHONY: push-docker # Build Docker image
push-docker:
	@docker push wallforfry/esiee-api:$(shell $(MAKE) version)
	@docker push wallforfry/esiee-api:latest

.PHONY: doc # Generate doc
doc:
	@swag init

.PHONY: help # Generate list of goals with descriptions
help:
	@echo "Available goals:\n"
	@grep '^.PHONY: .* #' Makefile | sed "s/\.PHONY: \(.*\) # \(.*\)/${FONT_BOLD}\1:${FONT_NC}\2~~/" | sed $$'s/~~/\\\n/g' | sed $$'s/~/\\\n\\\t/g'
