.PHONY: all test clean build install frontend sourcemod build
VERSION=v0.7.29
GO_CMD=go
GO_BUILD=$(GO_CMD) build
DEBUG_FLAGS = -gcflags "all=-N -l"
ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

all: frontend sourcemod buildp

fmt:
	go tool gci write . --skip-generated -s standard -s default
	go tool gofumpt -l -w .
	make -C frontend fmt

bump_deps:
	go get -u ./...
	make -C frontend update

buildp: frontend
	goreleaser release --clean

builds: frontend
	goreleaser release --clean --snapshot

watch:
	make -C frontend watch

serve:
	make -C frontend serve

frontend:
	make -C frontend

dist: frontend buildp
	zip -j gbans-`git describe --abbrev=0`-win64.zip build/win64/gbans.exe LICENSE README.md gbans_example.yml
	zip -r gbans-`git describe --abbrev=0`-win64.zip docs/
	zip -j gbans-`git describe --abbrev=0`-lin64.zip build/linux64/gbans LICENSE README.md gbans_example.yml
	zip -r gbans-`git describe --abbrev=0`-lin64.zip docs/

dist-master: frontend buildp
	zip -j gbans-master-win64.zip build/win64/gbans.exe LICENSE README.md gbans_example.yml
	zip -r gbans-master-win64.zip docs/
	zip -j gbans-master-lin64.zip build/linux64/gbans LICENSE README.md gbans_example.yml
	zip -r gbans-master-lin64.zip docs/

run:
	@go run $(GO_FLAGS) -race main.go

sourcemod:
	make -C sourcemod

sourcemod_devel: sourcemod
	docker cp sourcemod/plugins/gbans.smx srcds-localhost-1:/home/tf2server/tf-dedicated/tf/addons/sourcemod/plugins/
	docker restart srcds-localhost-1

test: test-go test-ts

test-ts:
	make -C frontend test

test-go:
	@go test $(GO_FLAGS) -race ./...

test-go-cover:
	@go test $(GO_FLAGS) -race -coverprofile coverage.out ./...
	@go tool cover -html=coverage.out

check: lint_golangci static vulncheck lint_ts typecheck_ts

vulncheck:
	go tool govulncheck

lint_golangci:
	go tool golangci-lint run --timeout 3m ./...

fix: fmt
	go tool golangci-lint run --fix

lint_ts:
	make -C frontend lint

typecheck_ts:
	make -C frontend typecheck

static:
	go tool staticcheck -go 1.24 ./...

clean:
	@go clean $(GO_FLAGS) -i
	rm -rf ./build/
	make -C frontend clean
	rm -rf ./sourcemod/plugins/gbans.smx

docker_test:
	@docker compose -f docker/docker-compose-test.yml up --force-recreate -V --remove-orphans
	@docker compose -f docker/docker-compose-test.yml rm -f

up_postgres:
	docker-compose --project-name dev -f docker/docker-compose-dev.yml down -v
	docker-compose --project-name dev -f docker/docker-compose-dev.yml up postgres --remove-orphans --renew-anon-volumes

docker_dump:
	docker exec gbans-postgres pg_dump -U gbans > gbans.sql

docker_restore:
	cat gbans.sql | docker exec -i docker-postgres-1 psql -U gbans

run_docker_snapshot: builds
	docker build . --no-cache -t gbans:snapshot
	docker run -it -v ./gbans.yml:/app/gbans.yml -v ./.cache:/app/.cache -p 6006:6006  gbans:snapshot

docs_install:
	make -C docs install

docs_start:
	make -C docs start

docs_deploy:
	make -C docs deploy

docs_build:
	make -C docs build
