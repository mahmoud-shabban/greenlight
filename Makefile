all: build run 

build: ./server/api/main.go
	@go build -o bin/app  ./server/api/

run:
	@./bin/app

up:
	@docker start db2
	@echo "Postgres db container started successfully"

down:
	@docker stop db2 > /dev/null 2>&1
	@echo "Postgres db container stoped successfully"

clean:
	@rm -rf ./bin/*
