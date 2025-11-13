# Include variables from .envrc file
include .envrc

# ====================================================================================== #
# HELPERS
# ====================================================================================== #
## help: print this help message
.PHONY: help
help:
	@echo "Usage"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/n] ' && read ans && [ "$${ans:-n}" = "y" ]


# ====================================================================================== #
# DEVELOPMENT
# ====================================================================================== #

.PHONY: all 
all: api/build api/run 

## api/build: build the app binary and save it to ./bin/app
.PHONY: api/build
api/build: ./server/api/main.go
	@go build -o bin/app  ./server/api/

## api/run: run the app with the default options
.PHONY: api/run
api/run:
	@./bin/app -db-dsn=${GREENLIGHT_DB_DSN} -smtp-username="" -smtp-password="" -limiter -smtp-port=1025 -smtp-host=localhost -smtp-sender="Greenlight <hello@mailpit.local>" -cors-trusted-origins="http://greenlight.local:8080 http://api.grenlight.local:8080 http://localhost:8080 http://192.168.0.118 http://localhost:9000 http://192.168.0.134:9000 http://192.168.0.134:8080"

## docker/up: starts postgresql and mailpit docker containers
.PHONY: docker/up
docker/up:
	@docker start db2 > /dev/null 
	@docker start mailpit > /dev/null 
	@echo "Postgres db & mailpit container started successfully"

## docker/down: stops postgresql and mailpit docker containers
.PHONY: docker/down
docker/down:
	@docker stop db2 mailpit > /dev/null 2>&1
	@echo "Postgres db & mailpit container stoped successfully"

## db/migrations/up: run db up migration to latest
.PHONY: db/migrations/up
db/migrations/up:
	@echo "Running SQL UP migrations"
	@migrate -path=./migrations -database=${GREENLIGHT_DB_DSN} up

## db/migrations/goto: run db migration to specific version {migration_version}
.PHONY: db/migrations/goto
db/migrations/goto: confirm
	@echo "Running Down Migration to version ${migration_version}"
	@migrate -path=./migrations -database=${GREENLIGHT_DB_DSN} goto ${migration_version}

## db/migrations/create: create new db migration file with name {name}
.PHONY: db/migrations/create
db/migrations/create: confirm
	@echo "Creating migrations files ${name}"
	@migrate create -dir ./migrations -seq -ext=.sql ${name}

# ====================================================================================== #
# QUALITY CONTROL
# ====================================================================================== #
## audit: performing QA tasks on the code [mode tidey, mod verify, fmt, vet, test] staticcheck tool
.PHONY: audit
audit:
	@echo 'Tidying and verifing module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formating code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo "Testing Code..."
	go test -race -vet=off ./...

## clean: removes the app binary from ./bin
.PHONY: clean
clean:
	@rm -rf ./bin/*

# Notes
# var ref ${var_name}
# pass var values to make using: make <target> var=value or it will search for var name in env variables available to make
# user lower case var name if it only used in makefile, use upper if not
# vars are case sensitive in make
# namespacing with / as separator e.g db/migrations/up db/migrations/down <side benifit of using /> it give us tab completion at the terminal make db/migrations and hit <tab> to test
# $ in makefile is used for variable ref to pass a literal $ to shell escape it with another $ that's why "$${ans:-n}" in cofirm is used
