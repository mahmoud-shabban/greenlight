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

greenlight_db_dsn ?= postgres://greenlight:pass@127.0.0.1/greenlight?sslmode=disable
limiter ?= true
smtp_username ?=  
smtp_password ?= 
smtp_port ?= 1025
smtp_host ?= localhost
smtp_sender ?= "Greenlight <hello@mailpit.local>"
cors_origins ?= "http://greenlight.local:8080 http://api.grenlight.local:8080 http://localhost:8080 http://192.168.0.118 http://localhost:9000 http://192.168.0.134:9000 http://192.168.0.134:8080"


output_dir ?= ./bin
git_description = $(shell git describe --always --dirty --tags --long)
current_time = $(shell date --iso-8601=seconds)

# -s for the linker to strip DWARF debbuging info and synmbol table from the binary
# -X to link the current_time value to the buildTime var in main module to burn in the binary build time available with -version
linker_flags ?= '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## api/build: build the app binary and save it to ./bin/app
.PHONY: api/build
api/build: ./server/api/main.go
	@go build -o ${output_dir}/api -ldflags=${linker_flags} ./server/api
	@GOARCH=amd64 GOOS=linux go build -o ${output_dir}/linux_amd64/api -ldflags=${linker_flags} ./server/api/

## api/run: run the app with the default options
.PHONY: api/run
api/run:
	@go run ./server/api \
	  -db-dsn=${greenlight_db_dsn} \
	  -limiter=${limiter} \
	  -smtp-username=${smtp_username} \
	  -smtp-password=${smtp_password} \
	  -smtp-port=${smtp_port} \
	  -smtp-host=${smtp_host} \
	  -smtp-sender=${smtp_sender} \
	  -cors-trusted-origins=${cors_origins}

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
# QUALITY CONTROL, Hooked to git as a pre-commit hook
# ====================================================================================== #
## qa/audit: performing QA tasks on the code [mode tidey, mod verify, fmt, vet, test] staticcheck tool
.PHONY: qa/audit
qa/audit: qa/vendor
	@echo 'Formating code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo "Testing Code..."
	go test -race -vet=off ./...

## qa/vendor: perform code and module vendoring 
.PHONY: qa/vendor 
qa/vendor:
	@echo 'Tidying and verifing module dependencies...'
	go mod tidy
	go mod verify
# 	@echo 'Vendoring the code...'
# 	go mod vendor
## clean: removes the app binary from ./bin
.PHONY: clean
clean:
	@rm -rf ./bin/*

# Notes for future
# var ref ${var_name}
# pass var values to make using: make <target> var=value or it will search for var name in env variables available to make
# user lower case var name if it only used in makefile, use upper if not
# vars are case sensitive in make
# namespacing with / as separator e.g db/migrations/up db/migrations/down <side benifit of using /> it give us tab completion at the terminal make db/migrations and hit <tab> to test
# $ in makefile is used for variable ref to pass a literal $ to shell escape it with another $ that's why "$${ans:-n}" in cofirm is used
