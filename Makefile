ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: update
update:
	go mod tidy
	go mod verify

.PHONY: test
test:
	gotestsum --format=testname -- ./... -tags=units

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate-openapi
generate-openapi:
	swag init --pd 2

.DEFAULT_GOAL:= help