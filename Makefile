GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
# project path
ROOT=$(shell pwd)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

## init-test-db: 初始化测试环境db
init-test-db:
	@echo "init test db..."
	@cd scripts && ./init_test_db.sh $(ROOT)

## init-local-db: 初始化本地环境db
init-local-db:
	@echo "init local db..."
	@cd scripts && ./init_local_db.sh $(ROOT)

migrate-local-db:
	@echo "init local db..."
	@cd scripts && ./migrate_local_db.sh $(ROOT)

.PHONY: test
# test
test:
	go test ./... -race -cover

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
