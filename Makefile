all: build run 

build: ./server/api/main.go
	@go build -o bin/app  ./server/api/

run:
	@./bin/app -smtp-username="" -smtp-password="" -smtp-port=1025 -smtp-host=localhost -smtp-sender="Greenlight <hello@mailpit.local>"

up:
	@docker start db2 > /dev/null 
	@docker start mailpit > /dev/null 
	@echo "Postgres db & mailpit container started successfully"

down:
	@docker stop db2 mailpit > /dev/null 2>&1
	@echo "Postgres db & mailpit container stoped successfully"

clean:
	@rm -rf ./bin/*
