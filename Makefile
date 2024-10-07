include .env

# ==================================================================================== #
# GO RELATED HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go test -race -vet=off ./...
	go mod verify


# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build: build the cmd/api application
.PHONY: build
build:
	go mod verify
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## run: run the cmd/api application
.PHONY: run
run: tidy build
	./bin/api


# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #
.PHONY: get_migrate
get-migrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xz -- migrate

## migrations-new name=$1: create a new database migration
.PHONY: migrations-new
migrations-new:
	./migrate create -seq -ext=.sql -dir=./assets/migrations ${name}
# ./migrate create -ext sql -dir assets/migrations -seq ${name}

## migrations-up: apply all up database migrations
## migrations-up number=N: apply next N number of database migrations
.PHONY: migrations-up
migrations-up:
	./migrate -path=./assets/migrations -database="${MASTER_DATABASE_URL}" up ${number}

## migrations-down: apply all down database migrations
## migrations-down number=N: apply previous N number of down database migrations
.PHONY: migrations-down
migrations-down:
	./migrate -path=./assets/migrations -database="${MASTER_DATABASE_URL}" down ${number}

## migrations-goto version=$1: migrate to a specific version number
.PHONY: migrations-goto
migrations-goto:
	./migrate -path=./assets/migrations -database="${MASTER_DATABASE_URL}" goto ${version}

## migrations-force version=$1: force database migration
.PHONY: migrations-force
migrations-force:
	./migrate -path=./assets/migrations -database="${MASTER_DATABASE_URL}" force ${version}


## migrations-version: print the current in-use migration version
.PHONY: migrations-version
migrations-version:
	./migrate -path=./assets/migrations -database="${MASTER_DATABASE_URL}" version

## migrations-drop: drop full database
.PHONY: migrations-drop
migrations-drop:
	./migrate -path=./assets/migrations -database="${MASTER_DATABASE_URL}" drop
# ./migrate -path=./assets/migrations drop


# ==================================================================================== #
# DOCKER RELATED HELPERS
# ==================================================================================== #

## docker-dev: run docker compose up --build
.PHONY: docker-build
docker-build:
	docker compose -f docker-compose.yml down
	docker compose -f docker-compose.yml up --build

.PHONY: docker-start
docker-start:
	docker compose -f docker-compose.yml up -d 

## docker-clean: delete container and volumes
.PHONY: docker-clean
docker-clean:
	docker compose -f docker-compose.yml down -v
