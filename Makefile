all: build run 

build: ./server/api/main.go
	@go build -o bin/app  ./server/api/

run:
	@./bin/app

clean:
	@rm -rf ./bin/*
